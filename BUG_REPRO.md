# Bug Reproduction

Bug: redaction, header handling, and retention results share mutable slice storage with their inputs or with later results.

Trigger: redact a payload, mutate nested lists or headers, and then inspect the original request and an earlier retention result.

Observed error: the captured input is changed by redaction/header processing, and an earlier result can be rewritten by a later result.
