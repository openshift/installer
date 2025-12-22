---
description: Suggest appropriate contributors for a platform
argument-hint: "[platform-name]"
---

## Name
platform-help

## synopsis
```
/platform-help [platform-name]
```

## Description
The `platform-help` command analyzes files associated with the provided platform to determine recent contributors. These contributors are considered knowledgeable on matters related to the platform.

The command performs the following analysis:
- Find all files related to the "[platform-name]" platform in the pkg directory and subdirectories
- Use git blame to find the contributors
- Select the 2 users with the most recent contributions

This command is particularly useful, to the [OpenShift Installer](https://github.com/openshift/installer) code base, when the user/caller is unfamiliar with the team. This has potential
uses with code reviews and general team slack help.

## Implementation

### Step 1: Determine the base branch
- Detect the main branch (usually `main` or `master`)
- Verify the base branch exists

```bash
git symbolic-ref refs/remotes/origin/HEAD | sed 's@^refs/remotes/origin/@@'
```

### Step 2: Find relevant files
- Use the base branch
- Find all files related to the "[platform-name]" platform in the pkg directory and subdirectories

### Step 3: Find Contributors to those files
- Use git blame to find the contributors to the files in the previous step

### Step 4: Select recent contributors
- Select the 2 users with the most recent contributions ignore "openshift-merge-bot"

## Examples

1. **Basic Usage**

```
/platform-help gcp
```

Output:
```

⏺ GCP Platform Contributors

  Based on analysis of 105 GCP-related files in the pkg directory, here are the 2 most recent contributors to the GCP platform:

  1. user-1 (User 1)

  - Email: user-1@something.com
  - Most recent contribution: 2025-12-11 19:20:35
  - By far the most active contributor with 2,955 total lines of code

  2. user-2 (User 2)

  - Email: user-2@something.com
  - Most recent contribution: 2025-11-20 13:37:56
  - Significant contributor with 598 total lines of code

  These two contributors have the most recent commits to GCP platform files and would be excellent resources for questions related to GCP in the OpenShift Installer codebase.
```

## Arguments
- `platform-name`: The name of the platform to find relevant contributors (ex: aws)
