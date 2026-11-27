# go-redact

Go integration with a separately managed Presidio Analyzer HTTP service.

The current deliverable is an HTTP-only smoke harness, not a production
redaction library. Detection belongs to Presidio; this repository does not
implement local passport, bank-account, or other detector patches.

## HTTP harness

See [the harness README](poc/presidio-http-harness/README.md) for running the
health/analysis checks and tests. Go 1.26.9 or newer is required.

The existing Windows launcher uses a workspace-specific Ubuntu WSL toolchain.
On Linux with a suitable Go installation, run from the harness module:

```sh
cd poc/presidio-http-harness
GOTOOLCHAIN=local go test -race ./...
GOTOOLCHAIN=local go run . -endpoint http://127.0.0.1:5002 -timeout 5s
```

The smoke command sends one fixed synthetic example after a successful health
check. It reports statuses and detection count without printing prompt or
response values. It neither installs nor starts Presidio.

Successful HTTP calls do not establish detector accuracy or complete PII removal.
