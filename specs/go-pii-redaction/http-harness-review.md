---
project: "go-pii-redaction"
date: "2026-11-27"
status: "review-complete-launcher-fixed"
scope: "accepted HTTP-only harness; no detector/model work"
---

# HTTP Harness: Atom Multi-Model Review

## Outcome

Two independent models reviewed the staged 11-file harness change and the approved brief through Atom's fixed-evidence mode: claude-sonnet-5.5 (1/1 completed) and gpt-6.1-terra (1/1 completed). Neither is the authoring model, gpt-6.1-sol. One pass per model; no consensus or statistical confidence is claimed from repetition. All required build/vet/changed-file-format/test/race gates passed before the models ran. No additional real Presidio request occurred.

One launcher error-reporting/exit-contract issue was corrected after developer approval. No detector/HTTP-protocol correctness blocker was established. Remaining findings are explicit scope limitations, bounded tradeoffs or low-priority robustness/style suggestions; their individual dispositions follow. The original model output is retained, including findings rejected here rather than silently removed.

Raw successful review: http-harness-review-retry.md. Command: atom review --staged --model claude-sonnet-5.5 --model gpt-6.1-terra --runs 1 --concurrency 2 --fixed-evidence --context-file <approved brief>, launched from poc\presidio-http-harness so the gate checks the actual nested Go module instead of the known empty root module.

## Failed Default Attempt Is Not A Clean Review

The installed Atom defaults were tried first. gpt-5.4 was unavailable; gemini-3.8-flash exited with an empty parent answer that Atom correctly rejected. The run returned exit 1 with no usable model review, retained in http-harness-review-raw.md. A small non-code probe confirmed Gemini can answer a trivial prompt; inspection of the exact review session confirmed the review answer was empty, not a harness finding. No installed tooling or permission guard was changed to force acceptance. The successful retry used two different current models, not the failed review mislabeled as clean.

## Toolchain And Advisory Limits

The successful Atom report records required Go build, vet, changed-file gofmt, tests and race as PASSED. Advisory staticcheck 2025.1.1 cannot decode Go 1.26 export data ("export data version 5 is greater than maximum supported version 2"); its WARNING is a checker compatibility failure, not a source diagnostic or a passing staticcheck result. errcheck and govulncheck are absent and did not run. No checker was installed/upgraded. Fixed-evidence review excludes tool use and has the usual caller-context limitation; Atom labels context provenance unvalidated. Prior live smoke evidence remains exactly the previously approved two requests.

## Deduplicated Findings And Triage

Raw model severities are medium/low; the consequence-based dispositions below do not manufacture agreement merely because Atom grouped nearby lines. In particular, different findings around lines 40-43 and 58-60 remain separate issues.

| Finding | Location | Raised by | Disposition and reason |
| --- | --- | --- | --- |
| Launcher failures can emit infrastructure diagnostics and exit with arbitrary codes despite the documented safe failure/exit-2 behavior. | poc/presidio-http-harness/main.go:40-43 | gpt-6.1-terra, medium | Fixed after developer approval. launchWSL suppresses stderr and withholds stdout on any launch/build/child failure, returns wsl_launch_or_child_failed, and main exits 2. Regression tests inject a runner to cover start/exit failures, successful output/arguments, and output-write failure without WSL or Presidio. |
| User-specific Linux toolchain path and Ubuntu distribution are not portable. | main.go:28 | claude-sonnet-5.5, medium | Deferred as an explicit limitation, not a newly discovered contract violation. Brief and both READMEs deliberately identify the launcher as workspace-specific, while native Linux use with Go 1.26.9+ is documented. A generic installer/configuration layer is outside this HTTP-only spike. |
| errors.As suggested for a direct exec.Cmd.Run result; safe reasons do not wrap raw errors. | main.go:40,58 | claude-sonnet-5.5, low | Reject as an API/style requirement here: this is a CLI, not an exported typed-error API. exec.Cmd.Run directly yields the concrete exit error, and raw URL-bearing transport errors are intentionally not wrapped into diagnostics. The separate launcher failure issue is retained above. |
| Final success-line write failure is ignored. | main.go:58 | gpt-6.1-terra, low | Low-priority robustness suggestion; the two prior status writes already propagate errors, and this does not change the actual HTTP result. Deferred: the native final success-line write remains unchecked; see http-harness/findings.md. Launcher output-write failures are now checked. |
| Configurable per-call timeout can be cut short by the fixed total deadline. | main.go:60 | claude-sonnet-5.5, medium | Reject as a bug: the approved contract explicitly combines a per-call upper bound with a fifteen-second total request budget. A timeout is a maximum, not a guaranteed minimum wait. Retain the total cap. |
| Proxy restriction lacks a dedicated test; custom transport lacks separate default dial/TLS timeouts. | main.go:77 | claude-sonnet-5.5, medium | No unbounded behavior established: Proxy:nil explicitly disables proxy forwarding, and http.Client.Timeout plus request context bound connection/read phases. Defer extra transport test granularity; no separate phase timing or tuned transport contract is claimed. |
| Non-200 response body is read before status is checked, and read errors are then secondary to status. | main.go:112 | claude-sonnet-5.5, medium | Bounded tradeoff, not a success or privacy bug: failure is returned either way, reading is capped at 1 MiB and the call deadline, and no body is logged. Defer early status rejection unless a requirement for it is added. |
| Constants lack comments and conditions are long. | main.go:18,126-137 | claude-sonnet-5.5, low | Reject as blocking work. Constants are self-describing in a fixed-input disposable command, and gofmt already passed. No behavior failure was identified; avoid an unrelated style refactor. |
| Captured evidence contains local paths and becomes historical. | poc/presidio-http-harness/check-output.txt:1 | claude-sonnet-5.5, low | Retain honest dated evidence as required by Atom; do not imply it is a fresh check for a later source revision. Review before public push; final validation and any changes will be recorded separately. No secrets/tokens/customer text were identified in the scoped files. |
| Test-handler t.Error/Fprint usage and absence of a unit test for Windows entry point. | poc/presidio-http-harness/main_test.go:28 | claude-sonnet-5.5, low | t.Error is allowed from test goroutines and reports failure; fixture write failures result in failed response checks. The actual Windows-to-WSL path was exercised by the saved live run, not claimed as instrumented unit coverage. Focused launcher regressions were added with the fix; no real API rerun occurred. |

## Publication Boundary

The user requested a PR, confirmed creating public girishmotwani/go-redact, then requested Atom multi-model review first. Only a local Git bootstrap and feature branch exist; no GitHub repository, remote push or PR has been created. The candidate contains the harness, its short root README and its own brief/report/evidence, not old detector/model spikes, global historical logs, vendored Atom tooling or local state. Bootstrap metadata was directly inspected; commits use the authenticated owner's GitHub noreply address without changing global Git identity. Developer approved fixing the launcher, validating, then creating the public repository and PR. Never merge automatically.

## Approved Triage And Post-Fix Verification

Developer selected "Fix launcher, validate, and create PR" on 2026-11-27. The existing HTTP implementation and request budget are unchanged. Windows failures intentionally collapse infrastructure and child errors into one safe reason rather than forwarding arbitrary diagnostics or exit codes. Successful child status output is retained. The final correction was directly inspected and tested; the two model reports above describe the pre-fix snapshot, not a second model pass on the correction.

Focused regression command: go test ./... -run TestLaunchWSL -count=1 (exit 0). Full fresh Atom gates, captured in poc/presidio-http-harness/post-review-check-output.txt, passed build, vet, module-wide formatting, tests and race.

Linux execution-platform verification, cached Go 1.26.9, GOTOOLCHAIN=local, fake HTTP servers only:

```text
ok  	spike/presidio-http-harness	1.429s	coverage: 74.4% of statements
```

Linux command: go test -race -cover ./... (exit 0). Coverage is not exhaustive. No new live service call occurred; the original report/run-output remain dated integration evidence, not post-fix live evidence. Deferrals are collected in http-harness/findings.md; scope and publication decisions are in http-harness/decisions.md.
