# Bug Reproduction

Bug: subscription cancellation and publisher revision handling lose lifecycle and state guards.

Trigger: cancel a subscriber while a worker waits, then exercise rejected, concurrent, and overflow publishes.

Observed error: the worker wait reaches `context deadline exceeded`; rejected publishes advance the revision and overflow is accepted.
