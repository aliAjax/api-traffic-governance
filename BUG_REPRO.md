# Bug Reproduction

Bug: concurrent semaphore accounting and mirror queue operations are not consistently protected.

Trigger: run the concurrent limiter and mirror checks while acquisitions, releases, pushes, and length reads overlap.

Observed error: the race detector reports concurrent reads and writes in `internal/limiter/concurrency.go` and `internal/mirror/queue.go`; the running count can remain non-zero and queue capacity/length assertions fail.
