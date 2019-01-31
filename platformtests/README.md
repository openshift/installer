# Platform Tests

This directory contains test suites checking per-platform assumptions.
These tests require platform access that is not appropriate for platform-agnostic unit tests.
The `t` directory name (unlike `tests`) allows us to share [the project `vendor` directory managed by `dep`](../docs/dev/dependencies.md#go).

Platforms:

* [aws](aws)
