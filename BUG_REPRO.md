# Bug Reproduction

Bug: xDS ACK, quota diff, and snapshot storage do not preserve the current-version and ownership invariants.

Trigger: ACK an old version, change quota content without changing list length, and mutate route input after compiling a snapshot.

Observed error: an unknown old ACK is accepted, equal-length quota changes are missed, and a stored snapshot changes when the caller mutates its input slice.
