# Agent-Based Installer Testing Guide

This document describes how to test the agent-based installer components in the OpenShift installer repository.

## Overview

The agent-based installer allows deploying OpenShift clusters without a centralized infrastructure management service. Testing these components involves both unit tests and integration tests with specific environment requirements.

## Unit Tests

Agent-related unit tests live alongside the code they exercise, following standard Go conventions:

- Agent asset tests: `pkg/asset/agent/`
- Agent manifests: `pkg/asset/agent/manifests/`
- Agent image tests: `pkg/asset/agent/image/`
- Agent workflow tests: `pkg/asset/agent/workflow/`

Run unit tests with:

```sh
go test ./pkg/asset/agent/...
```

Or via the containerized runner:

```sh
hack/go-test.sh
```

## Integration Tests

Agent integration tests are located in `cmd/openshift-install/` and require additional tooling.

### Prerequisites

- **oc** (OpenShift CLI): Required for image extraction and release payload inspection.
- **nmstatectl**: Required for validating network configuration in agent manifests.
- **Registry credentials**: Many tests query `registry.ci.openshift.org` for release image information. In CI, credentials are provided via the `AUTH_FILE` environment variable.

### Running Integration Tests

```sh
hack/go-integration-test.sh
```

This script runs tests matching `-run .Integration`, which includes agent integration tests such as `TestAgentIntegration`.

## Test Data and Fixtures

Agent integration test fixtures reside in:

- `cmd/openshift-install/testdata/agent/`

These fixtures contain sample `install-config.yaml` and `agent-config.yaml` files used to validate the agent installer workflow end-to-end.

## Common Patterns

### Mocking Assisted Service Responses

Agent tests that interact with the assisted-service API use mock HTTP servers. See existing tests for examples of setting up test servers that simulate cluster registration and host validation workflows.

### Validating Agent Manifests

Manifest generation tests verify that the agent installer produces the correct set of Kubernetes manifests. Test helpers in `pkg/asset/agent/manifests/util_test.go` provide shared utilities for building test install configs and asserting manifest content.
