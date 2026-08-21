# Bug Reproduction

Bug: loading an empty state leaves storage maps unusable on the zero-value path.

Trigger: start with an empty `state.json` and create the first service.

Observed error: `panic: assignment to entry in nil map`, with the stack entering `internal/repository.(*Store).CreateService`.
