# Bug Reproduction

Bug: reservation lifecycle and token waiting do not enforce committed-state, audit-time, and zero-deadline invariants.

Trigger: commit the same reservation twice, roll back the committed reservation, inspect its commit time, and wait with a zero deadline while tokens are unavailable.

Observed error: replay commit and rollback return nil, the commit timestamp is missing, and zero-deadline waiting lasts too long.
