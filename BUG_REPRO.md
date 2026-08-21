# Bug Reproduction

Bug: request cancellation and deadline state do not consistently reach health and fault operations.

Trigger: run a slow upstream health check and cancel a fault operation while it is waiting or holding a read lock.

Observed error: the health check exceeds its deadline, canceled fault work continues, and the read-lock delay path can block follow-up work.
