# Bug Reproduction

Bug: audit verification can return while its mutex is still held, and audit sealing does not bind all required state.

Trigger: corrupt an audit entry, call verification, and attempt an append afterward.

Observed error: `panic: test timed out after 300ms`; a goroutine is blocked in `sync.(*Mutex).Lock` from `AuditLog.Append` after verification fails.
