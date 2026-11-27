---
project: "go-pii-redaction"
date: "2026-11-26"
poc_slug: "presidio-http-harness"
verdict: "Confirmed"
status: "accepted-http-only"
---

# POC Report: Presidio HTTP Harness

## Question

Can a small Go command successfully call GET /health and POST /analyze on the existing local Presidio service and decode a valid analysis response within bounded time, without logging prompt text or response bodies?

## Verdict

**Confirmed for HTTP invocation only.** The new isolated Go command successfully made GET /health and POST /analyze against http://127.0.0.1:5002 inside Ubuntu WSL. Both returned HTTP 200; analysis decoded as one structurally valid detection. One live invocation, two serial requests, no retries or measured reruns. This is not a detection-accuracy or redaction/privacy claim.

## Approval And Scope

Developer redirected work away from detector/model repairs, selected the HTTP-only harness, then explicitly selected approve-harness-and-smoke for http-harness-brief.md. The approved budget was at most one health and one synthetic analysis call; that budget is now consumed. The original Refuted coverage result remains unchanged. The NER model comparison remains parked with no provisioning permission.

This separate disposable standard-library module is under poc\presidio-http-harness. It contains main.go and focused tests, not a production redaction package. No source was copied from an old experiment into production. No model/dependency download, new analyzer process, service edit/restart, corpus replay, local passport/bank-account detector, Go-side replacement, external provider call or subagent occurred.

## Real-Service Evidence

The Windows Atom controller built the command and invoked its temporary Linux child with the already cached Go 1.26.9 toolchain and GOTOOLCHAIN=local. HTTP requests ran wholly inside Ubuntu WSL. The child removes its exact temporary executable on exit. The request input is a fixed synthetic example; diagnostics below contain no prompt or response values.

```text
COMMAND: atom poc run presidio-http-harness --timeout 120 --no-colour
STARTED_UTC=2026-11-26T11:10:10.7930218Z
atom poc run presidio-http-harness  (command spike)

  PASS  vet     ok
  PASS  build   ok
  PASS  execute ok

  stdout:
    GET /health status=200
    POST /analyze status=200 detections=1
    result=passed scope=http_only

  Verdict: CONFIRMED  (exit 0) -- the question was answered yes
ENDED_UTC=2026-11-26T11:10:20.7818244Z
OUTER_EXIT=0
```

The child and outer exits are both 0 here; the explicit child Confirmed verdict and actual API statuses, not merely the outer recording status, establish the bounded result. The orchestration timestamps include build/launch and are not an HTTP latency measurement.

## Verification

Focused fake-HTTP unit tests cover successful and empty detections, exact request routing/parameters, status/schema/oversize failures, no retries/redirects, timeout behavior, invalid configuration and privacy-safe diagnostics. These local test servers are not the real analyzer and their fixture results are not detection evidence.

```text
COMMAND: go test -v . (fake HTTP servers only)
=== RUN   TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality
=== RUN   TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality/[]
=== RUN   TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality/[{"entity_type":"PERSON","start":11,"end":19,"score":0.85}]
--- PASS: TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality (0.01s)
    --- PASS: TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality/[] (0.01s)
    --- PASS: TestSmoke_InvokesBothAPIsWithoutGradingDetectionQuality/[{"entity_type":"PERSON","start":11,"end":19,"score":0.85}] (0.00s)
=== RUN   TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks
=== RUN   TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/unhealthy
=== RUN   TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/analysis_failure
=== RUN   TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/malformed_analysis
=== RUN   TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/oversized_analysis
--- PASS: TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks (0.03s)
    --- PASS: TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/unhealthy (0.00s)
    --- PASS: TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/analysis_failure (0.00s)
    --- PASS: TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/malformed_analysis (0.00s)
    --- PASS: TestSmoke_FailuresStopWithoutRetriesOrBodyLeaks/oversized_analysis (0.02s)
=== RUN   TestSmoke_RejectsInvalidConfigurationBeforeRequests
--- PASS: TestSmoke_RejectsInvalidConfigurationBeforeRequests (0.00s)
=== RUN   TestSmoke_DoesNotFollowRedirects
--- PASS: TestSmoke_DoesNotFollowRedirects (0.00s)
=== RUN   TestSmoke_EnforcesCallTimeout
--- PASS: TestSmoke_EnforcesCallTimeout (0.02s)
=== RUN   TestDetectionCount_RejectsMalformedSchemaAndInvalidRanges
--- PASS: TestDetectionCount_RejectsMalformedSchemaAndInvalidRanges (0.00s)
PASS
ok  	spike/presidio-http-harness	3.548s
EXIT=0
```

Full Windows isolated-module gates:

```text
COMMAND: atom check .\poc\presidio-http-harness --no-colour
atom check C:\Users\gimotwan\work\go-redact\poc\presidio-http-harness

  PASS  build  go build ./... passed (17.7s)
  PASS  vet    go vet ./... passed (2.3s)
  PASS  fmt    gofmt -l . reported no unformatted files (0.1s)
  PASS  test   go test ./... passed (6.5s)
  PASS  race   go test -race ./... passed (17.1s)

  PASSED in 43.6s
EXIT=0
```

Linux execution-platform tests:

```text
COMMAND: cached Linux Go 1.26.9 test -race -cover ./... (fake HTTP servers only)
ok  	spike/presidio-http-harness	1.227s	coverage: 64.9% of statements
EXIT=0
```

Linux unit coverage is exactly 64.9% of statements, not exhaustive coverage. The live smoke run is not included in that instrumented unit coverage. The actual Windows-to-WSL entry path is exercised by the real run above. Some I/O/error/launcher branches remain untested. The root module still has no production packages; its known no-package vet/test failures are separate from this passing isolated module.

## Behavior And Limits

The command accepts an HTTP(S) origin (no credentials, query, fragment or path prefix), defaults to loopback port 5002, and uses a positive configurable per-call timeout (default five seconds) within a fifteen-second request window. Response reads are capped at 1 MiB. Requests never follow redirects or use proxy forwarding. Health must succeed before the one analysis request is attempted. Analysis requires a JSON array, required typed fields, code-point ranges within the supplied input and scores in [0,1]. An empty array also passes: detector recall/precision is deliberately not the harness gate. Errors emit safe reasons, not raw URLs, request/response bodies or values.

The Windows launcher is specific to this workspace's Ubuntu distribution and cached toolchain path. This is not a portable runtime installer or lifecycle helper. Native use needs Go 1.26.9 or newer; the installed WSL default Go 1.22.2 is not sufficient with GOTOOLCHAIN=local. See the spike README for exact commands.

## Explicitly Not Covered

No Go redaction/replacement API, overlap handling, production error taxonomy, concurrency/load, LLM fail-closed integration, provider transmission, arbitrary user-prompt CLI, service recovery/idle availability, authentication, deployment packaging or universal endpoint compatibility. This smoke run did not re-audit the model or service package versions: that is not necessary to establish the two HTTP calls and is not claimed. A returned detection does not prove correct boundaries, precision, recall or removal of PII. Existing upstream detection gaps remain known limitations, not bugs patched in Go.

## Result Gate

Accept this HTTP-only harness result as the completed requested scope, leaving detector/model work parked. No further experiment, production client, redaction implementation or additional live calls are automatically authorized.

## Result Acceptance (2026-11-26)

Developer selected accept-http-harness: accept the HTTP-only harness result and stop at this scope. The completed harness is retained for use/review. Detector/model work remains parked; no further live invocation, production client, redaction implementation or experiment is authorized by this acceptance.
