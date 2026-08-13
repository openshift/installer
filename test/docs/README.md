# Installer Test Documentation

This directory contains structured test documentation for the OpenShift
Installer. Test plans and test cases are written as markdown files, organized
by platform or topic, and cross-linked to the Go test files that implement
them.

The goal is traceability: every testable behavior should have a documented
scenario that links to a JIRA feature or bug, the source code under test, and
the automated test that verifies it.

For a cross-platform overview of which features are tested and how, see the
[Platform Feature Matrix](platform_feature_matrix.md).

## Traceability

The test documentation supports end-to-end traceability through this chain:

```text
JIRA issue -> Feature -> Test Plan -> Test Cases -> Scenarios -> CI Job (Prow)
```

Each link in the chain is navigable:

| From | To | How |
|------|----|-----|
| JIRA issue | Test plan | JIRA issue description should link to the test plan file (see [JIRA linking](#jira-linking)) |
| Test plan | JIRA issue | `References` section contains JIRA links |
| Test plan | Test cases | Section 3.2 links to case files in `../cases/` |
| Test case | Test plan | `Related Links` section links back to the parent plan |
| Test case | Scenarios | Scenarios table at the top with anchor links |
| Scenario | Source code | `Test file` field with repo-relative path |
| Scenario | CI job logs | `CI job history` field with Prow URL (for Prow CI scenarios) |

### Navigating to Prow logs

For scenarios with status `Automated (Prow CI)`, include a `CI job history`
link using this URL pattern:

```text
https://prow.ci.openshift.org/job-history/gs/test-platform-results/logs/<full-job-name>
```

The full job name follows the Prow naming convention:

```text
periodic-ci-openshift-release-<branch>-ci-<version>-<job-suffix>
```

For example, the GCD E2E job `e2e-gcd-ovn-private-techpreview` on the `main`
branch for version `5.0` becomes:

```text
periodic-ci-openshift-release-main-ci-5.0-e2e-gcd-ovn-private-techpreview
```

To find the full job name for a CI job, search the
[openshift/release](https://github.com/openshift/release) repo for the job
suffix in `ci-operator/config/openshift/installer/openshift-installer-main.yaml`.

### JIRA linking

To complete the reverse link from JIRA to test docs, add a comment or update
the JIRA issue description with a link to the test plan in the repo:

```text
Test plan: https://github.com/openshift/installer/blob/main/test/docs/<platform>/plans/<plan>.md
Test cases: https://github.com/openshift/installer/blob/main/test/docs/<platform>/cases/<case>.md
```

This makes the chain navigable in both directions: from JIRA into the test
docs, and from the test docs back to JIRA.

## Directory Structure

```text
test/docs/
  README.md                          # This file
  <platform-or-topic>/
    plans/
      <feature-or-component>.md      # Test plans (higher-level)
    cases/
      <scenario-group>.md            # Test cases (scenario-level)
```

Each platform or topic gets its own directory. Within it, `plans/` holds
higher-level test plans and `cases/` holds detailed test case documents.

### Current contents

```text
test/docs/
  install_config_field_index.md       # Install-config field to feature mapping
  platform_feature_matrix.md          # Cross-platform feature test coverage matrix
  gcd/
    plans/
      gcd-feature.md                 # Feature test plan for GCD (OCPSTRAT-3006)
    cases/
      gcd_sovereign_install.md       # Detailed test cases for GCD installation
```

### Planned directories

As test documentation expands, add directories following the same pattern:

```text
test/docs/
  gcd/           # Google Cloud Dedicated (sovereign cloud)
  gcp/           # Google Cloud Platform (public cloud)
  aws/           # Amazon Web Services
  azure/         # Microsoft Azure
  vsphere/       # VMware vSphere
  baremetal/     # Bare metal / Agent-based
  nutanix/       # Nutanix
  ibmcloud/      # IBM Cloud
  powervs/       # IBM Power VS
  openstack/     # Red Hat OpenStack
  bugs/          # Bug-specific test plans (cross-platform)
```

The `bugs/` directory is for complex bugs that warrant their own test plan,
typically named after the JIRA key (e.g., `bugs/cases/OCPBUGS-12345.md`).

## Document Types

### Test Plan

A test plan is a higher-level document covering a feature, release, or
component. It answers: what are we testing, why, how, and what does "done"
look like.

**Location:** `<platform>/plans/<name>.md`

**Required sections:**

| Section | Purpose |
|---------|---------|
| 1. Introduction | Overview, scope (in/out), key features, JIRA references |
| 2. Testing Strategy | Schedule, test types (unit/integration/E2E), environments |
| 3. Test Areas and Test Cases | Activity table, links to case documents, coverage map |
| 4. Test Details | Technical details specific to the feature under test |
| 5. Risks | Risk table with impact and mitigation |
| 6. Exit Criteria | Conditions that must be met to consider testing complete |

### Test Case

A test case document describes specific test scenarios with enough detail that
someone could execute them manually or verify the automated test covers the
right behavior.

**Location:** `<platform>/cases/<name>.md`

**Required sections:**

| Section | Purpose |
|---------|---------|
| Metadata | Table with feature JIRA, component, test type, assignee |
| Description | What the test cases cover and how they are organized |
| Scenarios | Table with scenario number, link, and status (Automated/Manual) |
| Per-scenario | Test file, function, description, execution, pass/fail criteria |
| Manual Test Checklist | Checkbox checklist of all manual scenarios (only if file has manual tests) |
| Related Links | JIRA issues, PRs, parent test plan |

**Per-scenario fields:**

Each scenario should include:

- **Test file** - path to the `_test.go` file (or CI job name for E2E)
- **Function under test** - the Go function or CI workflow being validated
- **Automation status** - Automated (unit test), Automated (Prow CI), or Manual
- **Description** - what the scenario validates and why it matters
- **Execution** - code snippet or CLI command showing how to run it
- **Expected Result** - what a passing test produces
- **Pass/Fail Criteria** - table of inputs and expected outputs

## How to Add a New Test Plan

### 1. Create the directory

If the platform or topic directory does not exist yet, create it:

```sh
mkdir -p test/docs/<platform>/plans test/docs/<platform>/cases
```

### 2. Write the test plan

Create `test/docs/<platform>/plans/<name>.md` using this template:

```markdown
# Feature Test Plan: <Feature Name>

## 1. Introduction

### 1.1. Overview

<One paragraph describing the feature and its significance.>

### 1.2. Scope

**In Scope**

- <What is being tested>

**Out of Scope**

- <What is explicitly not being tested>

### 1.3. Key Features

- **<JIRA-KEY>** - <Feature description>

### 1.4. References

- Strategy: [<JIRA-KEY>](https://redhat.atlassian.net/browse/<JIRA-KEY>)
- <Other relevant PRs, docs, or links>

## 2. Testing Strategy

### 2.1. Schedule

| Milestone | Date / Sprint | Notes |
|-----------|---------------|-------|
| <Milestone> | TBD | <Notes> |

### 2.2. Test Types

- **Unit Tests:** <What unit tests cover>
- **Integration Tests:** <What integration tests cover>
- **E2E Tests:** <What E2E tests cover>

### 2.3. Test Environments

- **Unit/Integration:** <Environment description>
- **E2E CI:** <CI job and environment details>

## 3. Test Areas and Test Cases

### 3.1. Key Test Activities

| Activity | Type | Est. Duration |
|----------|------|---------------|
| <Activity> | <Unit/Integration/E2E> | <Duration> |

### 3.2. Test Cases

| Test Case | Doc | Test File(s) |
|-----------|-----|--------------|
| <Case name> | [<filename>.md](../cases/<filename>.md) | <Go test file paths> |

### 3.3. Unit Test Coverage Map

| Behavior | Source | Test File |
|----------|--------|-----------|
| <Behavior> | `<source-file>` `<Function>()` | `<test-file>` |

## 4. Test Details

<Technical details, invariants, configuration tables, E2E workflow steps.>

## 5. Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| <Risk> | <Impact> | <Mitigation> |

## 6. Exit Criteria

- <Condition 1>
- <Condition 2>
```

### 3. Link to test cases

The test plan should reference its test case documents in section 3.2
using relative links:

```markdown
| GCD Sovereign Install | [gcd_sovereign_install.md](../cases/gcd_sovereign_install.md) | Multiple |
```

## How to Add New Test Cases

### 1. Create the test case file

Create `test/docs/<platform>/cases/<name>.md` using this template:

````markdown
# Test Case: <Case Name>

## Metadata

| Field | Value |
|-------|-------|
| **Feature** | [<JIRA-KEY>](https://redhat.atlassian.net/browse/<JIRA-KEY>) - <Description> |
| **Component** | Installer / <platform> |
| **Test type** | <Unit, Integration, E2E> |
| **Assignee** | <Name or TBD> |

## Description

<What these test cases cover and how they are organized.>

## Scenarios

| # | Scenario | Status | Priority |
|---|----------|--------|----------|
| 1 | [<Scenario title>](#1-scenario-title) | Automated (unit test) | |
| 2 | [<Scenario title>](#2-scenario-title) | Manual | P1 |

---

## 1. <Scenario title>

- **Test file:** `<path/to/file_test.go>`
- **Function under test:** `<FunctionName>()`
- **Automation status:** Automated (unit test)

### Description

<What this scenario validates and why it matters.>

### Execution

```go
result := pkg.FunctionUnderTest(input)
// Expected: <expected value>
```

### Expected Result

- <What the function returns or what the test asserts>

### Pass/Fail Criteria

| Input | Expected Output |
|-------|-----------------|
| `<input>` | `<output>` |

---

## 2. <Scenario title>

- **Automation status:** Manual
- **Requires:** <Running cluster, credentials, etc.>

### Description

<What this scenario validates and why manual verification is needed.>

### Prerequisites

- <What must be in place before running this test>

### Execution

<Step-by-step instructions with commands to run.>

### Pass/Fail Criteria

| Check | Pass Condition |
|-------|---------------|
| <Check> | <Condition> |

---

## Manual Test Checklist

<!-- MANUAL_TESTS_START -->

- [ ] **2. <Scenario title>**
  - [ ] <Sub-check from pass/fail criteria>

<!-- MANUAL_TESTS_END -->

---

## Related Links

- Feature test plan: [<plan-name>.md](../plans/<plan-name>.md)
- <JIRA links, PR links>
````

### 2. Add the case to the test plan

Update the parent test plan's section 3.2 (Test Cases table) to include a row
linking to the new case file.

### 3. Add to the coverage map

If the test case covers unit-tested functions, add rows to the test plan's
section 3.3 (Unit Test Coverage Map) linking each behavior to its source and
test file.

## Conventions

### File naming

- **Plans:** Use kebab-case. Name after the feature or component:
  `gcd-feature.md`, `ipi-install.md`, `shared-vpc.md`.
- **Cases:** Use snake_case. Name after the scenario group:
  `gcd_sovereign_install.md`, `aws_govcloud_validation.md`.
- **Bug cases:** Name after the JIRA key: `OCPBUGS-12345.md`.

### Cross-linking

- Plans link down to cases using relative paths: `[name](../cases/file.md)`
- Cases link up to plans using relative paths: `[name](../plans/file.md)`
- Both link to JIRA using full URLs:
  `[JIRA-KEY](https://redhat.atlassian.net/browse/JIRA-KEY)`
- Both link to source code using repo-relative paths:
  `pkg/types/gcp/platform.go`

### JIRA references

Always include JIRA references in the metadata. For features, link the
strategy or epic. For bugs, link the bug issue directly. Use the format:

```markdown
| **Feature** | [OCPSTRAT-3006](https://redhat.atlassian.net/browse/OCPSTRAT-3006) - Description |
```

### Automation status

Use one of these values in test case scenarios:

| Value | Meaning |
|-------|---------|
| Automated (unit test) | Covered by a Go `_test.go` file |
| Automated (Prow CI) | Covered by a CI job in the release repo |
| Manual | Requires manual execution by a human |
| TBD | Automation status not yet determined |

Every scenario must have an automation status in two places:

1. **Scenarios table** - the `Status` column at the top of the case file
2. **Per-scenario metadata** - the `Automation status` field in the scenario

### Automation priority

The Scenarios table has an optional `Priority` column for tracking which
Manual or TBD scenarios should be automated and when:

| Value | Meaning |
|-------|---------|
| P1 | Automate next sprint |
| P2 | Automate soon (next 1-2 sprints) |
| P3 | Keep manual (automation not cost-effective) |
| (empty) | Already automated, or priority not yet assessed |

Leave the Priority column empty for already-automated scenarios. Use this
column to help Technical Leads and Software Engineers prioritize automation
work during sprint planning.

To list all scenarios awaiting automation, sorted by priority:

```bash
grep -n "| P[123]" test/docs/*/cases/*.md | sort -t'|' -k5
```

### Manual tests

Manual test scenarios are tests that cannot be (or are not yet) automated and
require a human to execute and verify. Common reasons for manual testing:

- Verifying cloud console state that is not exposed through CLI or API
- Validating behavior that depends on network topology (e.g., private DNS
  reachability from outside the VPC)
- Confirming resource cleanup after cluster destroy
- Validating error messages during interactive install-config creation
- Exploratory testing of new features before automation is written

#### Marking a scenario as manual

In the scenario's metadata, set:

```markdown
- **Automation status:** Manual
- **Requires:** <what is needed - running cluster, credentials, etc.>
```

In the scenarios table at the top of the file, set the Status column to
`Manual`:

```markdown
| 14 | [Verify cluster nodes use correct machine type](#14-...) | Manual |
```

#### Manual Test Checklist section

Every case file that contains manual scenarios must include a
**Manual Test Checklist** section near the end (before Related Links). This
section provides a runnable checklist with `- [ ]` checkboxes that can be
copied and used to track progress during a manual verification run.

Wrap the checklist with HTML comments so it can be found programmatically:

```markdown
<!-- MANUAL_TESTS_START -->
...checklist content...
<!-- MANUAL_TESTS_END -->
```

See [gcd_sovereign_install.md](gcd/cases/gcd_sovereign_install.md) for a
complete example.

#### Finding manual tests

To list all case files that have manual tests:

```bash
grep -rl "MANUAL_TESTS_START" test/docs/*/cases/
```

To list all case files with manual tests for a specific platform:

```bash
grep -rl "MANUAL_TESTS_START" test/docs/gcd/cases/
```

To see a quick summary of which scenarios are manual across all case files:

```bash
grep -n "| Manual" test/docs/*/cases/*.md
```

To extract just the manual test checklist from a specific file:

```bash
sed -n '/MANUAL_TESTS_START/,/MANUAL_TESTS_END/p' test/docs/gcd/cases/gcd_sovereign_install.md
```

To count manual vs automated scenarios across all platforms:

```bash
echo "Manual:    $(grep -rc '| Manual' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"
echo "Automated: $(grep -rc '| Automated' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"
```

### Keeping docs in sync

Test documentation should be updated when:

- A new feature adds testable installer behavior
- A test file is added, renamed, or deleted
- A CI job is added or modified
- A complex bug fix warrants its own test case
- Test coverage gaps are identified during review

The documentation does not need to mirror every individual Go test function.
Focus on documenting behaviors and scenarios, not individual assertions.

## Execution Status and Reporting

This section helps Engineering Managers and Technical Leads assess test
execution status and coverage readiness.

### Checking CI execution status

View recent Prow job runs for the installer:

- **All installer jobs:** [Prow job list](https://prow.ci.openshift.org/?repo=openshift%2Finstaller)
- **Specific job history:** Use the URL pattern from the [Traceability](#navigating-to-prow-logs) section

### Generating a coverage report

Use `check.py` and `grep` to produce a quick coverage summary:

```bash
# Validate all test docs and show errors
python3 test/docs/check.py

# Count scenarios by automation status across all platforms
echo "=== Scenario counts ==="
echo "Automated (unit): $(grep -rc '| Automated (unit test)' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"
echo "Automated (Prow): $(grep -rc '| Automated (Prow CI)' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"
echo "Manual:           $(grep -rc '| Manual' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"
echo "TBD:              $(grep -rc '| TBD' test/docs/*/cases/*.md | awk -F: '{s+=$NF}END{print s}')"

# List scenarios awaiting automation, sorted by priority
echo "=== Automation backlog ==="
grep -n "| P[123]" test/docs/*/cases/*.md | sort -t'|' -k5

# Count platforms with test documentation
echo "=== Platform coverage ==="
grep -c "| Yes" test/docs/platform_feature_matrix.md
```

### Release readiness assessment

Use the [Platform Feature Matrix](platform_feature_matrix.md) to assess
readiness. Key indicators:

- **NT (Not Tested)** cells for features being shipped indicate coverage gaps
- **Test Documentation** rows (bottom of matrix) show which platforms have
  documented test plans and cases
- Footnotes provide context on CI job counts and known limitations

## Test Lifecycle

### When to update test documentation

Update scenario documentation when:

- A test file is renamed, moved, or deleted
- The function under test changes signature or behavior
- A CI job is renamed or reconfigured in the release repo
- A Manual scenario is automated (change status and add test file reference)

### When to retire a scenario

Remove a scenario when:

- The feature it tests is removed from the installer
- The scenario is fully replaced by a different scenario (not just refactored)
- The CI job it references is permanently deleted

Before removing, check that no other documents link to the scenario's anchor.

### Deprecation process

When a scenario will be retired but is not yet removed:

1. Add `(deprecated)` to the scenario title in the Scenarios table
2. Add a note at the top of the scenario detail section explaining why and
   what replaces it
3. Remove the scenario in the next PR that touches the same case file

Do not leave deprecated scenarios for more than one release cycle.

### Keeping source references current

When a test file or function is renamed, update all references to it across
the test docs. Use grep to find stale references:

```bash
# Find references to a specific test file
grep -rn "old_file_test.go" test/docs/

# Find references to a specific function
grep -rn "OldFunctionName" test/docs/
```

## References

These external examples informed the structure used here:

- [Kueue operator test docs](https://github.com/openshift/kueue-operator/pull/1994) -
  test plans and test cases with CI job traceability
- [Node runc deprecation cases](https://github.com/asahay19/origin/blob/f1280dc76adb20a53ea1291b2496545a4eec5da5/test/extended/node/runcdeprecationcases.md) -
  scenario-level test case format with metadata tables and pass/fail criteria
- [Node runc upgrade cases](https://github.com/asahay19/origin/blob/4258dcd738b8cb62aa946dd4886beb962133b7e8/test/extended/node/runc_upgrade_cases.md) -
  step-by-step test case format with prerequisites and expected results
