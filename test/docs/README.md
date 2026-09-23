# Installer Test Documentation

This directory contains structured test documentation for the OpenShift
Installer, organized by platform or topic:

- **Plans** (`<platform>/plans/`) cover a feature end to end: scope, strategy,
  which behavior is covered by which Go test or CI job, risks, exit criteria.
- **Cases** (`<platform>/cases/`) are one file per test a human runs, written
  as setup, steps with expected results, and cleanup.

The goal is traceability: every testable behavior links to a JIRA feature or
bug, to the source under test, and to whatever verifies it - a Go test, a CI
job, or a documented manual procedure.

For a cross-platform overview of which features are tested and how, see the
[Platform Feature Matrix](platform_feature_matrix.md).

## Traceability

The test documentation supports end-to-end traceability through this chain:

```text
JIRA issue -> Feature -> Test Plan -> Go tests / CI job (Prow)
                                   -> Test Cases (manual)
```

Each link in the chain is navigable:

| From | To | How |
|------|----|-----|
| JIRA issue | Test plan | JIRA issue description should link to the test plan file (see [JIRA linking](#jira-linking)) |
| Test plan | JIRA issue | `References` section contains JIRA links |
| Test plan | Test cases | Section 3.2 links to case files in `../cases/` |
| Test case | Test plan | `Test plan` metadata field links back to the parent plan |
| Test plan | Source code | Section 3.3 maps each behavior to its source and test file |
| Test plan | CI job logs | Section 4 links the Prow job history for the feature's E2E job |

### Navigating to Prow logs

When a test plan names an E2E CI job, link its job history using this URL
pattern:

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
      <one-case>.md                  # Test cases (one file per case)
```

Each platform or topic gets its own directory. Within it, `plans/` holds
higher-level test plans and `cases/` holds one file per manual test case.

### Current contents

```text
test/docs/
  install_config_field_index.md       # Install-config field to feature mapping
  platform_feature_matrix.md          # Cross-platform feature test coverage matrix
  gcd/
    plans/
      gcd-feature.md                 # Feature test plan for GCD (OCPSTRAT-3006)
    cases/
      gcd_node_machine_and_disk.md   # Manual: nodes use C3 + Hyperdisk
      gcd_private_dns_only.md        # Manual: no public DNS zone or record
      gcd_destroy_cleanup.md         # Manual: destroy leaves nothing behind
      gcd_reject_public_publish.md   # Manual: publish External is rejected
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

A test case document describes **one** test that a human runs, in enough detail
that someone unfamiliar with the feature can execute it. The structure follows
IEEE-829: setup, then alternating steps and expectations, then cleanup.

**Location:** `<platform>/cases/<name>.md` - one file per case.

**Required sections:**

| Section | Purpose |
|---------|---------|
| Title + metadata table | Feature JIRA, component, type, priority, parent plan |
| Intro paragraph | Two or three lines on what the case proves and why it matters |
| `## Setup` | Environment, credentials, and state required before starting |
| `## Test` | Alternating `### Step` / `### Expect` pairs |
| `## Cleanup` | What to tear down afterwards, or "None" |

**What does not belong here:** automated tests. A Go test or a Prow job is its
own documentation, and a prose copy of it drifts out of date. Record those in
the parent plan's coverage map (section 3.3) and its E2E workflow description
(section 4) instead, so each behavior is described in exactly one place.

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

| Test Case | Doc | Priority |
|-----------|-----|----------|
| <Case name> | [<filename>.md](../cases/<filename>.md) | <P1/P2/P3> |

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
| Cluster uses private DNS only | [gcd_private_dns_only.md](../cases/gcd_private_dns_only.md) | P1 |
```

## How to Add New Test Cases

### 1. Create the test case file

Create one file per case, `test/docs/<platform>/cases/<name>.md`, using this
template:

````markdown
# Test Case: <What the test proves, as a statement>

| Field | Value |
|-------|-------|
| **Feature** | [<JIRA-KEY>](https://redhat.atlassian.net/browse/<JIRA-KEY>) - <Description> |
| **Component** | Installer / <platform> |
| **Type** | Manual |
| **Priority** | <P1/P2/P3> |
| **Test plan** | [<plan-name>.md](../plans/<plan-name>.md) |

<Two or three lines: what this case proves, and why it cannot be automated.>

## Setup

- <State, credentials, or tooling needed before starting.>

```bash
export SOME_VAR="<value>"
```

## Test

### Step

<One action, with the exact command to run.>

```bash
<command>
```

### Expect

<The observable result. Be specific enough to judge pass or fail.>

### Step

<Next action.>

### Expect

<Next result.>

## Cleanup

<What to tear down, or "None." when every step is read-only.>
````

Keep a case to a single page. If it grows past roughly a dozen steps, it is
probably two cases.

### 2. Add the case to the test plan

Update the parent test plan's section 3.2 (Test Cases table) to include a row
linking to the new case file.

### 3. Validate

```sh
python3 test/docs/check.py cases
```

## Conventions

### File naming

- **Plans:** Use kebab-case. Name after the feature or component:
  `gcd-feature.md`, `ipi-install.md`, `shared-vpc.md`.
- **Cases:** Use snake_case. Name after what the single case checks:
  `gcd_private_dns_only.md`, `aws_govcloud_validation.md`.
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

### Type

The `Type` metadata field records why the case is a document rather than code:

| Value | Meaning |
|-------|---------|
| Manual | Requires a human to execute and judge |
| TBD | Written ahead of the automation that will replace it |

Automated coverage is not documented in `cases/`. Record it in the parent
plan's coverage map (section 3.3) or its E2E workflow description (section 4).

### Priority

`Priority` tracks whether the case should be automated, and when:

| Value | Meaning |
|-------|---------|
| P1 | Automate next sprint |
| P2 | Automate soon (next 1-2 sprints) |
| P3 | Keep manual (automation is not cost-effective) |

To list the automation backlog, highest priority first:

```bash
grep -l '^| \*\*Type\*\* | Manual' test/docs/*/cases/*.md |
  xargs grep -H '^| \*\*Priority\*\*' | sort -t'|' -k3
```

### Manual tests

A case is manual when a human has to be in the loop. Common reasons:

- Verifying cloud-side state that the installer's tests cannot observe
- Validating behavior that depends on network topology (e.g., private DNS
  reachability from outside the VPC)
- Confirming resource cleanup after cluster destroy
- Validating error messages during interactive install-config creation
- Exploratory testing of a new feature before automation is written

Because one file is one case, discovery is just a listing:

```bash
# All manual cases
grep -rl '^| \*\*Type\*\* | Manual' test/docs/*/cases/

# Manual cases for one platform
grep -rl '^| \*\*Type\*\* | Manual' test/docs/gcd/cases/

# Count them
grep -rl '^| \*\*Type\*\* | Manual' test/docs/*/cases/ | wc -l
```

When a case is automated, delete its file and add the behavior to the parent
plan's coverage map in the same PR as the new test.

### Keeping docs in sync

Test documentation should be updated when:

- A new feature adds testable installer behavior
- A test file is added, renamed, or deleted
- A CI job is added or modified
- A complex bug fix warrants its own test case
- Test coverage gaps are identified during review

The documentation does not need to mirror every individual Go test function.
Map behaviors to their test file in the plan's coverage map; leave individual
assertions to the tests.

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

# Count cases still needing a human, per platform
echo "=== Manual cases ==="
grep -rl '^| \*\*Type\*\* | Manual' test/docs/*/cases/ | cut -d/ -f3 | uniq -c

# List the automation backlog, highest priority first
echo "=== Automation backlog ==="
grep -l '^| \*\*Type\*\* | Manual' test/docs/*/cases/*.md |
  xargs grep -H '^| \*\*Priority\*\*' | sort -t'|' -k3

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

### When to update a case

Update a case when the commands it runs change: a renamed CLI flag, a changed
error message, a new prerequisite. A case whose steps no longer run is worse
than no case at all.

### When to retire a case

Delete the case file, and drop its row from the plan's section 3.2, when:

- The case is now automated - add the behavior to the plan's coverage map in
  the same PR as the new test
- The feature it covers is removed from the installer

There is no deprecation period. One case per file means deleting a file is a
complete, reviewable removal.

### Keeping plan references current

The plan's coverage map points at source and test files, so it goes stale when
those are renamed. Use grep to find the references:

```bash
grep -rn "old_file_test.go" test/docs/
grep -rn "OldFunctionName" test/docs/
```

## References

The case format is the IEEE-829 shape (setup, steps and expected results,
cleanup) recommended for OpenShift test documentation. These examples informed
the rest of the structure:

- [Kueue operator test docs](https://github.com/openshift/kueue-operator/pull/1994) -
  test plans and test cases with CI job traceability
- [Node runc upgrade cases](https://github.com/asahay19/origin/blob/4258dcd738b8cb62aa946dd4886beb962133b7e8/test/extended/node/runc_upgrade_cases.md) -
  step-by-step test case format with prerequisites and expected results
