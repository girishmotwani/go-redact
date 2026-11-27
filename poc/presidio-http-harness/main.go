// This disposable command checks Presidio's HTTP interface.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"
)

const sample = "My name is Jane Doe."
const responseLimit = 1 << 20

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		fmt.Println("Usage: presidio-http-harness [-endpoint http://127.0.0.1:5002] [-timeout 5s]")
		return
	}
	if runtime.GOOS == "windows" {
		if err := launchWSL(os.Args[1:], os.Stdout, (*exec.Cmd).Run); err != nil {
			fmt.Println("result=failed reason=" + err.Error())
			os.Exit(2)
		}
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	err := smoke(ctx, os.Args[1:], os.Stdout)
	cancel()
	if err != nil {
		fmt.Println("result=failed reason=" + err.Error())
		os.Exit(2)
	}
	fmt.Println("result=passed scope=http_only")
}

func launchWSL(arguments []string, output io.Writer, run func(*exec.Cmd) error) error {
	const launcher = `set -eu
binary=$(mktemp /tmp/go-redact-http-harness.XXXXXX)
trap 'rm -f -- "$binary"' EXIT
/home/girish/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.9.linux-amd64/bin/go build -o "$binary" .
set +e
"$binary" "$@"
status=$?
exit "$status"
`
	launcherArguments := []string{"--distribution", "Ubuntu", "--exec", "env", "GOTOOLCHAIN=local", "sh", "-c", launcher, "http-harness"}
	command := exec.Command("wsl.exe", append(launcherArguments, arguments...)...)
	var captured bytes.Buffer
	// Publish child output only after success; WSL/build errors can echo arguments.
	command.Stdout, command.Stderr = &captured, io.Discard
	if err := run(command); err != nil {
		return fmt.Errorf("wsl_launch_or_child_failed")
	}
	if _, err := io.Copy(output, &captured); err != nil {
		return fmt.Errorf("status_output_failed")
	}
	return nil
}

func smoke(ctx context.Context, arguments []string, output io.Writer) error {
	flags := flag.NewFlagSet("presidio-http-harness", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	endpoint := flags.String("endpoint", "http://127.0.0.1:5002", "analyzer base URL")
	timeout := flags.Duration("timeout", 5*time.Second, "per-call timeout")
	if flags.Parse(arguments) != nil || flags.NArg() != 0 || *timeout <= 0 {
		return fmt.Errorf("invalid_flags_or_timeout")
	}
	base, err := url.Parse(*endpoint)
	if err != nil || (base.Scheme != "http" && base.Scheme != "https") || base.Hostname() == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery || base.Fragment != "" || base.RawPath != "" || (base.Path != "" && base.Path != "/") {
		return fmt.Errorf("endpoint_must_be_http_origin_without_credentials_query_or_fragment")
	}
	client := &http.Client{Timeout: *timeout, Transport: &http.Transport{Proxy: nil}, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	defer client.CloseIdleConnections()
	origin := strings.TrimRight(base.String(), "/")
	if _, err := call(ctx, client, http.MethodGet, origin, "/health", nil); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, "GET /health status=200"); err != nil {
		return fmt.Errorf("status_output_failed")
	}
	payload, err := json.Marshal(map[string]any{"text": sample, "language": "en", "score_threshold": 0.5, "return_decision_process": false})
	if err != nil {
		return fmt.Errorf("payload_encoding_failed")
	}
	body, err := call(ctx, client, http.MethodPost, origin, "/analyze", payload)
	if err != nil {
		return err
	}
	count, err := detectionCount(body, utf8.RuneCountInString(sample))
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(output, "POST /analyze status=200 detections=%d\n", count); err != nil {
		return fmt.Errorf("status_output_failed")
	}
	return nil
}

func call(ctx context.Context, client *http.Client, method, origin, path string, payload []byte) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, origin+path, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("%s_request_creation_failed", path)
	}
	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/json")
	}
	response, err := client.Do(request)
	if err != nil {
		// Transport errors can echo URLs; the harness emits only a safe reason.
		return nil, fmt.Errorf("%s_transport_or_timeout_failed", path)
	}
	body, readErr := io.ReadAll(io.LimitReader(response.Body, responseLimit+1))
	closeErr := response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s_status_%d", path, response.StatusCode)
	}
	if readErr != nil || closeErr != nil {
		return nil, fmt.Errorf("%s_response_read_or_close_failed", path)
	}
	if len(body) > responseLimit {
		return nil, fmt.Errorf("%s_response_too_large", path)
	}
	return body, nil
}

func detectionCount(body []byte, characters int) (int, error) {
	var detections []struct {
		Type  *string  `json:"entity_type"`
		Start *int     `json:"start"`
		End   *int     `json:"end"`
		Score *float64 `json:"score"`
	}
	if !utf8.Valid(body) || !bytes.HasPrefix(bytes.TrimSpace(body), []byte("[")) || json.Unmarshal(body, &detections) != nil {
		return 0, fmt.Errorf("analysis_array_invalid")
	}
	for _, detection := range detections {
		if detection.Type == nil || strings.TrimSpace(*detection.Type) == "" || detection.Start == nil || detection.End == nil || detection.Score == nil || *detection.Start < 0 || *detection.Start >= *detection.End || *detection.End > characters || *detection.Score < 0 || *detection.Score > 1 {
			return 0, fmt.Errorf("analysis_detection_invalid")
		}
	}
	return len(detections), nil
}
