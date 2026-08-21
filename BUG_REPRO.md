# Bug Reproduction

Bug: wrapped domain errors are not preserved or recognized across the application and HTTP layers.

Trigger: pass wrapped not-found, invalid-input, conflict, and validation errors through the service and error writer.

Observed error: wrapped not-found is written as HTTP 500 (`wrapped not found code = 500`); validation errors lose their sentinel and remain plain rejection text.
