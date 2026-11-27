---
project: "go-pii-redaction"
date: "2026-11-26"
status: "approved-harness-and-single-smoke-run"
poc_slug: "presidio-http-harness"
---

# Brief: Presidio HTTP Harness

## Scope Selected By The Developer

HTTP-only Go harness: health/analyze, configurable endpoint, timeouts and safe error reporting. No Go-side detection or replacement. Detector quality belongs in Presidio and is not a prerequisite for this deliverable. The model comparison remains parked.

## Question

Can a small Go command successfully call GET /health and POST /analyze on the existing local Presidio service and decode a valid analysis response within bounded time, without logging prompt text or response bodies?

## Existing Evidence And Why This Is Small

The accepted availability report already exercised health and analysis. The accepted detection-feasibility report records 910 completed Go analysis requests without transport/schema errors. Its Refuted verdict concerns coverage, not HTTP integration. Reuse those findings; do not start another detector-quality or model-comparison study. This harness makes the HTTP smoke check easy to rerun rather than repeating the corpus benchmark. Disposable experiment source is not promoted into a production client.

## Concrete Deliverable

One separate standard-library Go command under poc\presidio-http-harness, with focused tests and short usage documentation. Configure the analyzer base URL (default http://127.0.0.1:5002) and a positive per-call timeout (default five seconds). Use the existing approved arrangement: Windows Atom orchestration, HTTP client inside Ubuntu WSL, cached Go 1.26.9 and GOTOOLCHAIN=local. No tool/dependency download.

A run makes one GET /health, then, only if healthy, one POST /analyze using the fixed synthetic text "My name is Jane Doe.", English, score_threshold 0.5 and return_decision_process false. No arbitrary user-prompt inputs in this first command. Require HTTP 200; cap each response at 1 MiB; parse analyze as a JSON array with valid entity_type, integer start/end code-point ranges, and score values. An empty array is a valid invocation result, not a detector-quality failure. Output only safe endpoint labels, statuses, detection count and success/failure; never text, detected values, response bodies or credentials.

Validate configuration, status errors, bounded reading, malformed JSON/fields, timeout behavior and safe diagnostics with local fake-HTTP unit tests. Distinguish those tests from the real-service smoke run. No automatic retries, redirects/proxy forwarding of the analysis request, service restart, or fallback detector.

## Success, Failure And Non-Coverage

Success: the real health call succeeds and real analyze returns a structurally valid array, with literal safe output retained. Failure/blocked run: nonzero exit and explicit safe reason; no partial success claimed, no implicit restart or repair. An unavailable service is a blocked integration check, not permission to change shared infrastructure.

This does not prove detection recall/precision, privacy completeness, redaction/offset replacement, overlap merging, fail-closed LLM integration, concurrency/load, availability recovery, or a production Go API. It does not change the original coverage verdict or thresholds. No recognizers, model changes, local passport/bank-account patches, held-out corpora, new manifests, provider calls or subagents.

## Bounded Execution And Verification

Approval authorizes only this isolated harness and one live smoke invocation: at most two serial HTTP requests, five seconds each, fifteen seconds overall request budget. Build/test orchestration may take up to 120 seconds. Use synthetic data only and the unchanged retained service; no installations, analyzer-arm startup, systemd edits, WSL shutdown/restart, unrelated process termination or model downloads. If the live attempt fails, retain its output and stop rather than spending another request budget silently.

Run targeted tests during development, then full atom check for the separate module, including race. Capture the literal real run through atom poc run presidio-http-harness. Inspect the child verdict rather than interpreting an outer recording exit as success. The root module's known no-package vet/test failures remain a separate pre-existing limitation, not a passing production baseline. Preserve accepted old spikes and frozen inputs unchanged.

## Gate

Approve building this exact HTTP-only harness and the single bounded real-service smoke run. This does not authorize production client/redaction implementation, model work, or the parked NER provisioning. Present the actual outcome afterward; do not automatically advance to a larger client.

## Approval (2026-11-26)

Developer selected approve-harness-and-smoke after presentation of this exact HTTP-only scope and at-most-two-request real-service budget. Build/test this isolated command and run that one bounded smoke invocation. No production client, redaction implementation, detector/model work, service restart or additional live request budget is approved.
