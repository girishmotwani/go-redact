package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestLaunchWSL_SuppressesLaunchAndChildFailures(t *testing.T) {
	for _, failure := range []error{errors.New("private launch details"), &exec.ExitError{Stderr: []byte("private build details")}} {
		var output bytes.Buffer
		err := launchWSL([]string{"-endpoint", "http://private.invalid"}, &output, func(command *exec.Cmd) error {
			fmt.Fprintln(command.Stdout, "private stdout")
			fmt.Fprintln(command.Stderr, "private stderr")
			return failure
		})
		if err == nil || err.Error() != "wsl_launch_or_child_failed" || output.Len() != 0 {
			t.Fatalf("unsafe failure: error=%v output=%q", err, output.String())
		}
	}
}

func TestLaunchWSL_ForwardsArgumentsAndSuccessfulOutput(t *testing.T) {
	const status = "result=passed scope=http_only\n"
	var output bytes.Buffer
	err := launchWSL([]string{"-timeout", "1s"}, &output, func(command *exec.Cmd) error {
		arguments := command.Args
		if arguments[0] != "wsl.exe" || arguments[len(arguments)-2] != "-timeout" || arguments[len(arguments)-1] != "1s" {
			t.Fatalf("unexpected command arguments: %q", arguments)
		}
		_, err := fmt.Fprint(command.Stdout, status)
		return err
	})
	if err != nil || output.String() != status {
		t.Fatalf("error=%v output=%q", err, output.String())
	}
}

type failedOutput struct{}

func (failedOutput) Write([]byte) (int, error) {
	return 0, errors.New("private output failure")
}

func TestLaunchWSL_ReportsOutputFailureSafely(t *testing.T) {
	err := launchWSL(nil, failedOutput{}, func(command *exec.Cmd) error {
		_, err := fmt.Fprintln(command.Stdout, "result=passed scope=http_only")
		return err
	})
	if err == nil || err.Error() != "status_output_failed" {
		t.Fatalf("unexpected output error: %v", err)
	}
}

func TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality(t *testing.T) {
	for _, body := range []string{"[]", `[{"entity_type":"PERSON","start":11,"end":19,"score":0.85}]`} {
		t.Run(body, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				calls.Add(1)
				if request.URL.Path == "/health" && request.Method == http.MethodGet {
					fmt.Fprint(writer, "healthy")
					return
				}
				var payload struct {
					Text      string  `json:"text"`
					Language  string  `json:"language"`
					Threshold float64 `json:"score_threshold"`
					Decision  *bool   `json:"return_decision_process"`
				}
				if request.URL.Path != "/analyze" || request.Method != http.MethodPost || request.Header.Get("Content-Type") != "application/json" || json.NewDecoder(request.Body).Decode(&payload) != nil || payload.Text != sample || payload.Language != "en" || payload.Threshold != 0.5 || payload.Decision == nil || *payload.Decision {
					t.Error("unexpected analysis route, headers or payload")
					writer.WriteHeader(400)
					return
				}
				fmt.Fprint(writer, body)
			}))
			t.Cleanup(server.Close)
			var output bytes.Buffer
			if err := smoke(context.Background(), []string{"-endpoint", server.URL + "/"}, &output); err != nil {
				t.Fatal(err)
			}
			if calls.Load() != 2 || !strings.Contains(output.String(), "GET /health status=200") || !strings.Contains(output.String(), "POST /analyze status=200 detections=") || strings.Contains(output.String(), sample) || strings.Contains(output.String(), "Jane") {
				t.Fatal("wrong request count or unsafe/missing status output")
			}
		})
	}
}

func TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks(t *testing.T) {
	for _, test := range []struct {
		name   string
		health int
		body   string
		status int
		calls  int32
	}{
		{"unhealthy", 503, "private upstream body", 200, 1},
		{"analysis failure", 200, "private upstream body", 500, 2},
		{"malformed analysis", 200, "private upstream body", 200, 2},
		{"oversized analysis", 200, strings.Repeat("x", responseLimit+1), 200, 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				calls.Add(1)
				if request.URL.Path == "/health" {
					writer.WriteHeader(test.health)
					fmt.Fprint(writer, "private health body")
					return
				}
				writer.WriteHeader(test.status)
				fmt.Fprint(writer, test.body)
			}))
			t.Cleanup(server.Close)
			var output bytes.Buffer
			err := smoke(context.Background(), []string{"-endpoint", server.URL}, &output)
			if err == nil || calls.Load() != test.calls || strings.Contains(err.Error()+output.String(), "private") || strings.Contains(err.Error()+output.String(), sample) {
				t.Fatal("failure was accepted, retried or leaked values")
			}
		})
	}
}

func TestSmoke_RejectsInvalidConfigurationBeforeRequests(t *testing.T) {
	for _, arguments := range [][]string{
		{"-endpoint", "http://user:secret@127.0.0.1:5002"},
		{"-endpoint", "file:///private"},
		{"-endpoint", "http://localhost:5002/analyze"},
		{"-endpoint", "http://localhost:5002?secret=value"},
		{"-endpoint", "http://localhost:5002#secret"},
		{"-endpoint", "http://:5002"},
		{"-timeout", "0s"}, {"-timeout", "-1s"}, {"-timeout", "private"},
		{"-unknown", "secret"}, {"unexpected"},
	} {
		var output bytes.Buffer
		err := smoke(context.Background(), arguments, &output)
		if err == nil || output.Len() != 0 || strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "private") {
			t.Fatal("unsafe configuration accepted or echoed")
		}
	}
}

func TestSmoke_DoesNotFollowRedirects(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		calls.Add(1)
		http.Redirect(writer, request, "/unexpected", http.StatusFound)
	}))
	t.Cleanup(server.Close)
	if err := smoke(context.Background(), []string{"-endpoint", server.URL}, &bytes.Buffer{}); err == nil || calls.Load() != 1 {
		t.Fatal("redirect accepted or followed")
	}
}

func TestSmoke_EnforcesCallTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		<-request.Context().Done()
	}))
	t.Cleanup(server.Close)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	started := time.Now()
	err := smoke(ctx, []string{"-endpoint", server.URL, "-timeout", "20ms"}, &bytes.Buffer{})
	if err == nil || time.Since(started) >= time.Second || strings.Contains(err.Error(), server.URL) {
		t.Fatal("timeout missing or transport error echoed URL")
	}
}

func TestDetectionCount_RejectsMalformedSchemaAndInvalidRanges(t *testing.T) {
	for _, body := range []string{
		`null`, `{}`, `[null]`, `[] trailing private data`,
		`[{"entity_type":"PERSON","end":19,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":null,"end":19,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":11.5,"end":19,"score":0.8}]`,
		`[{"entity_type":"","start":11,"end":19,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":-1,"end":19,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":11,"end":11,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":11,"end":21,"score":0.8}]`,
		`[{"entity_type":"PERSON","start":11,"end":19}]`,
		`[{"entity_type":"PERSON","start":11,"end":19,"score":1.1}]`,
		`[{"entity_type":"PERSON","start":11,"end":19,"score":-0.1}]`,
	} {
		if _, err := detectionCount([]byte(body), 20); err == nil || strings.Contains(err.Error(), "private") {
			t.Fatal("invalid schema accepted or echoed")
		}
	}
}
