COMMAND: atom review --staged --model claude-sonnet-5.5 --model gpt-6.1-terra --runs 1 --concurrency 2 --fixed-evidence --context-file <approved http-harness-brief.md>; cwd=poc\presidio-http-harness
# Code review

Reviewed staged changes: 593 changed lines across 11 files.

- `README.md`
- `poc/presidio-http-harness/README.md`
- `poc/presidio-http-harness/check-output.txt`
- `poc/presidio-http-harness/go.mod`
- `poc/presidio-http-harness/linux-test-output.txt`
- `poc/presidio-http-harness/main.go`
- `poc/presidio-http-harness/main_test.go`
- `poc/presidio-http-harness/run-output.txt`
- `poc/presidio-http-harness/unit-output.txt`
- `specs/go-pii-redaction/http-harness-brief.md`
- `specs/go-pii-redaction/http-harness-report.md`

| Model | Runs completed |
| --- | --- |
| `claude-sonnet-5.5` | 1/1 |
| `gpt-6.1-terra` | 1/1 |

## Toolchain evidence

| Command | Verdict |
| --- | --- |
| `go build ./...` | PASSED |
| `go vet ./...` | PASSED |
| `gofmt (changed Go files)` | PASSED - Checked 2 changed Go file(s) from staged changes: 'poc/presidio-http-harness/main.go', 'poc/presidio-http-harness/main_test.go' |
| `go test ./...` | PASSED |
| `go test -race ./...` | PASSED |
| `staticcheck ./...` | WARNING - staticcheck ./... reported findings with exit code 1 |
| `errcheck ./...` | not installed - not checked |
| `govulncheck ./...` | not installed - not checked |

    Formatting checks selected snapshots from staged changes: 'poc/presidio-http-harness/main.go', 'poc/presidio-http-harness/main_test.go'

## Review evidence

Execution: fixed-evidence; tool policy: source-exclusions-requested.
Scope: diff-plus-caller-context. Context quality: **unvalidated**.
Other repository/dependency material is not automatically included.
Logical prompt (before platform stdin normalization): 44374 UTF-8 bytes; SHA-256 `26fb80f852620ee267a3413bfeb462cb4eb496250284adf7d174bfc2735da3ed`.

Document provenance: caller-supplied/unverified.

- context-1: 5100 UTF-8 bytes; SHA-256 `b6e1c5e108b9295562052c2cab105702a19bf22b50491765859d14cd70d8d80e`.

## Findings

### `poc/presidio-http-harness/main.go`

#### poc/presidio-http-harness/main.go:28 - medium

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 28 (medium): The Windows launcher hardcodes a user-specific absolute toolchain path (/home/girish/...) and the 'Ubuntu' distribution, so it silently breaks on any other machine or user.

#### poc/presidio-http-harness/main.go:40-43 - medium

Severity is disputed: medium, low.

Reported at lines 40, 43; nearby lines are grouped, so this entry may cover more than one problem.

Several models reported problems in this group, but on different lines, so this is not recorded as agreement.

Reported by: `claude-sonnet-5.5` 1/1, `gpt-6.1-terra` 1/1

- `claude-sonnet-5.5` run 1, line 40 (low): The type assertion `err.(*exec.ExitError)` should use errors.As, and the launch failure drops the underlying error with no wrapped context.
- `gpt-6.1-terra` run 1, line 43 (medium): WSL or launcher build failures propagate their original exit code and diagnostics instead of exiting 2 with a safe failure reason.

#### poc/presidio-http-harness/main.go:58-60 - medium

Severity is disputed: medium, low.

Reported at lines 58, 60; nearby lines are grouped, so this entry may cover more than one problem.

Reported by: `claude-sonnet-5.5` 1/1, `gpt-6.1-terra` 1/1

- `claude-sonnet-5.5` run 1, line 58 (low): Errors are plain fmt.Errorf strings used as sentinel reasons rather than wrapped or typed errors, so callers and tests can only match by string.
- `gpt-6.1-terra` run 1, line 58 (low): Ignoring the final output error lets the command exit successfully even when its success result could not be written.
- `claude-sonnet-5.5` run 1, line 60 (medium): Overall deadline is fixed at 15s in main while the per-call timeout is configurable, so a -timeout above 7.5s is silently cut short by the context rather than honored.

#### poc/presidio-http-harness/main.go:77 - medium

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 77 (medium): `&http.Transport{Proxy: nil}` is not a restriction that tests verify, and this custom Transport drops DefaultTransport settings (dial timeout, TLS handshake timeout) so only the overall client Timeout bounds connection phases.

#### poc/presidio-http-harness/main.go:112 - medium

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 112 (medium): A response body read error is masked by a non-200 status check ordering and, more importantly, on a non-200 response the body is read fully (up to 1MiB) before the status is checked, wasting bandwidth and time.

#### poc/presidio-http-harness/main.go:18 - low

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 18 (low): Package-level constants `sample` and `responseLimit` lack comments explaining why (fixed synthetic text, 1 MiB cap).

#### poc/presidio-http-harness/main.go:126 - low

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 126 (low): Validation order `!utf8.Valid || HasPrefix('[') || Unmarshal` is a long one-line condition, and the entity validation condition at line 137 is a very long unreadable boolean that is hard to maintain.

### `poc/presidio-http-harness/check-output.txt`

#### poc/presidio-http-harness/check-output.txt:1 - low

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 1 (low): Committed generated tool output files (check-output.txt, unit-output.txt, run-output.txt, linux-test-output.txt) embed a user's local paths and go stale against the source.

### `poc/presidio-http-harness/main_test.go`

#### poc/presidio-http-harness/main_test.go:28 - low

Reported by: `claude-sonnet-5.5` 1/1

- `claude-sonnet-5.5` run 1, line 28 (low): Handler writes with t.Error from a server goroutine and ignores fmt.Fprint errors; the test also never exercises the Windows launcher path.

EXIT=0
