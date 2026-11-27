# Presidio HTTP Harness

Can a small Go command successfully call GET /health and POST /analyze on the existing local Presidio service and decode a valid analysis response within bounded time, without logging prompt text or response bodies?

**Result:** Confirmed for HTTP invocation only. One real run returned health=200 and analyze=200 with one structurally valid detection. This does not grade detector accuracy or implement redaction.

## Run From The Workspace Root (Windows PowerShell)

```powershell
$env:GOTOOLCHAIN = 'local'
atom poc run presidio-http-harness --timeout 120 --no-colour
```

Each invocation makes at most two real requests: GET /health, then POST /analyze with a fixed synthetic example. Do not run a corpus or model comparison. The Windows launcher uses Ubuntu and this workspace's already-cached Linux Go 1.26.9 toolchain; it does not install anything, start/restart the analyzer, or bridge Windows networking. The existing analyzer must already be reachable on WSL loopback.

To customize the endpoint/timeout, run from this module using Go 1.26.9+:

```powershell
Push-Location .\poc\presidio-http-harness
go run . -endpoint http://127.0.0.1:5002 -timeout 5s
Pop-Location
```

On Linux, use Go 1.26.9+ directly from this module, or the already cached toolchain executable. The default WSL Go 1.22.2 is too old. `-h` prints usage without HTTP calls. On any failure, the executable exits 2 with a safe reason. Windows launch/build/child failures use the generic wsl_launch_or_child_failed reason; raw child diagnostics are withheld. `go run` can report that as outer exit 1; inspect its printed child error. Atom's recorded child verdict is authoritative.

## Contract

Only an HTTP(S) origin is supported: no credentials, query, fragment or path prefix. Default per-call timeout is five seconds; total request window is fifteen seconds. Bodies are limited to 1 MiB. No redirects, proxies, retries, arbitrary user-text inputs, local recognizers, model tuning or replacement logic. Analysis must decode as a valid detection array; empty detections pass because this tests HTTP invocation, not accuracy. Logs contain API labels, statuses and counts, never prompts or response values.

## Test

```powershell
atom check .\poc\presidio-http-harness --no-colour
```

From the module: `go test -race -cover ./...` uses local fake HTTP servers only. Post-review Windows full checks and Linux race tests passed; Linux unit coverage is 74.4%, not exhaustive. Fresh Windows evidence is post-review-check-output.txt; Linux evidence and review dispositions are in specs\go-pii-redaction\http-harness-review.md. Original pre-review logs are check-output.txt, unit-output.txt, linux-test-output.txt and run-output.txt; they remain historical evidence. The full report is specs\go-pii-redaction\http-harness-report.md from the workspace root.

This is an isolated disposable command, not the production go-redact package. The accepted detector-quality baseline remains Refuted; detector/model experiments remain parked. No new live request budget or further implementation is implied by the saved result.
