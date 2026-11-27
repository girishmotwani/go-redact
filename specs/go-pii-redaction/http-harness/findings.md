# HTTP Harness: Open Findings

Scope: accepted disposable HTTP-only harness, reviewed 2026-11-27. No blocking correctness finding remains after the approved Windows launcher correction. Full individual triage and rejected findings: ../http-harness-review.md.

- Windows runtime portability is deferred: Ubuntu and the cached user-specific Go executable are intentional workspace requirements. Native Linux use with Go 1.26.9+ is documented; a portable installer is outside scope.
- The native final success-line write error remains unchecked (Minor robustness limitation). The HTTP status writes and Windows output-copy failure are checked; broader CLI output refactoring is deferred.
- Dedicated proxy/dial/TLS phase tests are deferred. Proxy forwarding is disabled in code, and client/context deadlines bound connection/read phases; exhaustive transport coverage is not claimed.
- Early rejection before reading non-200 response bodies is deferred. Reads are capped at 1 MiB and bounded by deadlines; failure bodies are never logged.
- staticcheck 2025.1.1 cannot read Go 1.26 export data. errcheck and govulncheck are absent. No passing result is claimed for these advisory tools; none was installed/upgraded.
- The root module has no production packages; checks must target the nested harness module. This is not a production redaction library or evidence of detection quality.
