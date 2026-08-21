# Bug Reproduction

Bug: circuit-breaker half-open probes, failure transitions, state reads, and retry-budget boundaries are inconsistent under concurrency.

Trigger: release several requests at the half-open boundary, fail a probe, read state concurrently, and exercise the configured retry limit.

Observed error: a second half-open probe is accepted, a failed probe remains half-open, the race detector reports a `State` read/write race, and the configured retry boundary is rejected.
