# HTTP Harness: Publication Decisions

- The developer accepted the HTTP-only harness and the single successful two-request smoke result. Detector/model experiments remain parked. No additional live service call is authorized by this publication.
- The developer requested a public girishmotwani/go-redact repository and a PR, then required Atom multi-model review before publication.
- Atom defaults failed to return usable reviews. That failure is retained, not relabeled clean. A successful retry used claude-sonnet-5.5 and gpt-6.1-terra, neither the authoring model; each ran once with required gates enabled.
- The developer selected "Fix launcher, validate, and create PR" after the merged findings were presented. Only the Windows launcher error path and focused tests changed after model review. A runner injection boundary allows offline regression tests; no generic process framework or detector patch was added.
- Windows infrastructure and child failures now share a safe reason and exit 2. Arbitrary child diagnostics/codes are deliberately not forwarded. Successful child output is retained.
- Publish only the harness, its project overview, metadata, brief, accepted report, evidence, and review/triage. Historical detector/model spikes, global local research logs, vendored Atom tooling, and runtime state remain outside the PR.
- Use gh directly for this accepted standalone spike rather than fabricating production-phase artifacts for the phase-oriented Atom shipping interface. The installed atom ship has no --pr flag from the skill example. Fresh full Atom checks remain mandatory before the feature commit.
- Create a PR for human review; never merge it automatically. No production-phase completion is claimed.
