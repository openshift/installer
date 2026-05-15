# GitHub Actions Workflows

This directory contains GitHub Actions workflows for the openshift-splat-team installer fork.

## Workflows

### vendor-check.yml
**Purpose:** Ensures the vendor directory is in sync with go.mod and go.sum

**Runs on:** PRs and pushes to master/main

**What it checks:**
- Runs `go mod tidy` and `go mod vendor`
- Verifies no changes to go.mod, go.sum, or vendor/

**How to fix failures:**
```bash
go mod tidy
go mod vendor
git add go.mod go.sum vendor/
git commit -m "Update vendor directory"
```

---

### lint.yml
**Purpose:** Runs golangci-lint to check code quality

**Runs on:** PRs and pushes to master/main

**What it checks:**
- Uses `.golangci.yaml` configuration
- Runs all configured linters

**How to fix failures:**
```bash
golangci-lint run --config .golangci.yaml
# Fix reported issues
```

---

### format-check.yml
**Purpose:** Verifies code is properly formatted

**Runs on:** PRs and pushes to master/main

**What it checks:**
- gofmt formatting
- goimports (warning only)

**How to fix failures:**
```bash
gofmt -w .
# Or use goimports for import organization too
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

---

### build-test.yml
**Purpose:** Builds the installer and runs unit tests

**Runs on:** PRs and pushes to master/main

**What it does:**
- Builds openshift-install binary via `hack/build.sh`
- Runs all unit tests with race detection
- Uploads coverage to Codecov (if configured)

**How to run locally:**
```bash
# Build
./hack/build.sh

# Test
go test -v -race ./...
```

## Adding New Workflows

When adding new workflows:
1. Follow the existing naming convention
2. Use the same trigger events (pull_request, push)
3. Use `go-version-file: 'go.mod'` for consistent Go versions
4. Document the workflow in this README
5. Provide clear error messages with remediation steps

## Notes

These workflows supplement (not replace) the OpenShift CI presubmit tests defined in the openshift/release repository. They provide quick feedback for common issues before the full CI suite runs.
