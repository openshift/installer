---
description: Find platform experts by analyzing recent git contributions
argument-hint: "[platform-name] [count]"
---

## Name
platform-experts

## Synopsis
```
/platform-experts [platform-name] [count]
/platform-experts --list
```

## Description
The `platform-experts` command analyzes files associated with the provided platform to determine recent contributors. These contributors are considered knowledgeable on matters related to the platform.

The command performs the following analysis:
- Dynamically discover available platforms from the codebase
- Validate the platform name exists in the codebase
- Find all files related to the "[platform-name]" platform across multiple directories
- Use git log to efficiently find contributors from the last 12 months
- Exclude vendor directories from contribution analysis (vendored dependencies don't count)
- Filter out bot accounts (openshift-merge-bot, openshift-ci-robot, dependabot)
- Rank contributors by both number of commits and lines changed
- Select the top N contributors (default: 2, configurable)

This command is particularly useful for the [OpenShift Installer](https://github.com/openshift/installer) code base when the user/caller is unfamiliar with the team. This has potential uses with code reviews and general team slack help.

## Implementation

**IMPORTANT: Use ONLY bash built-ins and git commands.**

This implementation uses only:
- Bash built-ins (if/while/for, string manipulation, arrays)
- Git commands (log, shortlog, show)
- Standard Unix tools (grep, sed, awk, sort, uniq, head, cut, wc)

### Step 1: Determine the base branch
- Detect the main branch (usually `main` or `master`)
- Verify the base branch exists

```bash
base_branch=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@')
if [ -z "$base_branch" ]; then
  base_branch="main"
fi
```

### Step 1.5: Validate platform name and discover platforms
- Dynamically discover platforms by checking existing directories
- Validate the requested platform exists

```bash
platform="$1"

# Dynamically discover all available platforms
get_platforms() {
  (
    # Get platforms from pkg/types (exclude common dirs and files)
    cd pkg/types 2>/dev/null && ls -1d */ 2>/dev/null | sed 's|/$||' | \
      grep -v '^common$\|^conversion$\|^defaults$\|^dns$\|^featuregates$\|^validation$'

    # Get additional platforms from pkg/asset/machines
    cd pkg/asset/machines 2>/dev/null && ls -1d */ 2>/dev/null | sed 's|/$||'
  ) | sort -u
}

# Get full platform name
get_platform_fullname() {
  case "$1" in
    agent)       echo "Agent-based Installer" ;;
    aws)         echo "Amazon Web Services" ;;
    azure)       echo "Microsoft Azure" ;;
    baremetal)   echo "Bare Metal" ;;
    external)    echo "External Infrastructure" ;;
    gcp)         echo "Google Cloud Platform" ;;
    ibmcloud)    echo "IBM Cloud" ;;
    imagebased)  echo "Image-based Installer" ;;
    none)        echo "No Platform" ;;
    nutanix)     echo "Nutanix" ;;
    openstack)   echo "OpenStack" ;;
    ovirt)       echo "oVirt / Red Hat Virtualization" ;;
    powervc)     echo "IBM PowerVC" ;;
    powervs)     echo "IBM Power Virtual Server" ;;
    vsphere)     echo "VMware vSphere" ;;
    *)           echo "$1" ;;  # Fallback to short name for unknown platforms
  esac
}

# Show platform list if requested
if [ "$1" = "--list" ]; then
  echo "Available platforms in OpenShift Installer:"
  get_platforms | while read -r p; do
    fullname=$(get_platform_fullname "$p")
    echo "  - $p ($fullname)"
  done
  exit 0
fi

# Validate platform exists
if [ ! -d "pkg/types/${platform}" ] && [ ! -d "pkg/asset/machines/${platform}" ]; then
  echo "Error: Unknown platform '${platform}'"
  echo ""
  echo "Available platforms:"
  get_platforms | sed 's/^/  - /'
  echo ""
  echo "Use '/platform-experts --list' to see all platforms"
  exit 1
fi
```

### Step 2: Find relevant files
Search in multiple locations for comprehensive coverage:
- `pkg/types/${platform}/` - Platform type definitions
- `pkg/asset/machines/${platform}/` - Machine asset generation
- `pkg/asset/installconfig/${platform}/` - Install config validation
- `pkg/tfvars/${platform}/` - Terraform variable generation (if applicable)
- `data/data/${platform}/` - Platform-specific templates
- `upi/${platform}/` - User-Provisioned Infrastructure docs
- `cluster-api/providers/${platform}/` - Cluster API provider code

```bash
# Build list of directories to search (only existing ones)
paths=""
for dir in "pkg/types/${platform}" "pkg/asset/machines/${platform}" \
           "pkg/asset/installconfig/${platform}" "pkg/tfvars/${platform}" \
           "data/data/${platform}" "upi/${platform}" \
           "cluster-api/providers/${platform}"; do
  if [ -d "$dir" ]; then
    paths="$paths $dir"
  fi
done
```

### Step 3: Find Contributors
- Use `git log` and `git shortlog` for performance
- Filter out bot accounts: openshift-merge-bot, openshift-ci-robot, dependabot
- Consider commits from the last 12 months for recency
- Weight by both number of commits AND lines changed

```bash
count="${2:-2}"  # Default to 2 contributors
time_range="12 months ago"

# Get contributor statistics using git shortlog and numstat
# Exclude vendor directories from contribution analysis
git log --since="$time_range" --no-merges --pretty=format:"%an|%ae|%at" \
  --numstat -- $paths ':!vendor/' ':!*/vendor/' 2>/dev/null | \
  grep -v "openshift-merge-bot\|openshift-ci-robot\|dependabot" | \
  awk -F'|' '
    # Track lines added/removed and commit timestamps per author
    /^[0-9]/ {
      added+=$1; removed+=$2
      lines[author] += ($1 + $2)
    }
    /\|/ {
      author=$1; email=$2; timestamp=$3
      commits[author]++
      emails[author]=email
      if (timestamp > latest[author]) latest[author]=timestamp
    }
    END {
      for (a in commits) {
        # Score = commits + (lines/100) for balanced weighting
        score = commits[a] + (lines[a]/100)
        print score"|"latest[a]"|"a"|"emails[a]"|"commits[a]"|"lines[a]
      }
    }
  ' | \
  sort -t'|' -k1,1nr -k2,2nr | \
  head -n "$count"
```

### Step 4: Select recent contributors
- Parse the sorted output to get top N contributors
- Get recent commit messages for each contributor

```bash
# For each contributor, get their recent commit messages
while IFS='|' read -r score timestamp name email commit_count line_count; do
  # Get recent commit titles (last 3)
  recent_work=$(git log --since="$time_range" --author="$email" \
    --pretty=format:"%s" --no-merges -- $paths ':!vendor/' ':!*/vendor/' 2>/dev/null | head -3 | \
    sed 's/^/  - /')

  # Convert timestamp to readable date
  date=$(date -d "@$timestamp" "+%Y-%m-%d %H:%M:%S" 2>/dev/null || echo "recent")

  # Output formatted result
  echo "Contributor: $name"
  echo "  Email: $email"
  echo "  Last commit: $date"
  echo "  Commits: $commit_count, Lines: $line_count"
  echo "  Recent work:"
  echo "$recent_work"
  echo
done
```

### Step 5: Handle special cases
- If no contributors found in last 12 months, extend to 24 months
- If platform has fewer contributors than requested count, return all found contributors
- Suggest checking OWNERS files as fallback

```bash
# Count results
result_count=$(git log --since="12 months ago" --no-merges \
  --pretty=format:"%ae" -- $paths ':!vendor/' ':!*/vendor/' 2>/dev/null | \
  grep -v "openshift-merge-bot\|openshift-ci-robot\|dependabot" | \
  sort -u | wc -l)

if [ "$result_count" -eq 0 ]; then
  echo "No contributors found in last 12 months, trying 24 months..."
  time_range="24 months ago"
  # Re-run analysis with extended time range
fi
```

## Examples

1. **Basic Usage**

```
/platform-experts gcp
```

Output:
```
⏺ GCP Platform Contributors

Based on analysis of 105 GCP-related files across multiple directories, here are the 2 most recent contributors to the GCP platform (last 12 months):

1. user-1 (User 1)
   - Email: user-1@something.com
   - Most recent contribution: 2025-12-11 19:20:35
   - Total contributions: 2,955 lines across 42 files
   - Recent work: "Add support for GCP shared VPC", "Fix GCP machine validation"
   - GitHub: @user-1

2. user-2 (User 2)
   - Email: user-2@something.com
   - Most recent contribution: 2025-11-20 13:37:56
   - Total contributions: 598 lines across 18 files
   - Recent work: "Update GCP quota validation", "Refactor GCP defaults"
   - GitHub: @user-2

These contributors have the most recent commits to GCP platform files and would be excellent resources for questions related to GCP in the OpenShift Installer codebase.
```

2. **List Available Platforms**

```
/platform-experts --list
```

Output:
```
Available platforms in OpenShift Installer:
  - agent (Agent-based Installer)
  - aws (Amazon Web Services)
  - azure (Microsoft Azure)
  - baremetal (Bare Metal)
  - external (External Infrastructure)
  - gcp (Google Cloud Platform)
  - ibmcloud (IBM Cloud)
  - imagebased (Image-based Installer)
  - none (No Platform)
  - nutanix (Nutanix)
  - openstack (OpenStack)
  - ovirt (oVirt / Red Hat Virtualization)
  - powervc (IBM PowerVC)
  - powervs (IBM Power Virtual Server)
  - vsphere (VMware vSphere)

Usage: /platform-experts [platform-name] [count]
```

Note: The list is dynamically generated from the codebase and will automatically reflect new or removed platforms. Full names are provided for clarity.

3. **Custom Contributor Count**

```
/platform-experts aws 5
```

Output:
```
⏺ AWS Platform Contributors

Based on analysis of 287 AWS-related files, here are the 5 most recent contributors to the AWS platform (last 12 months):

1. contributor-1 (Contributor One)
   - Email: contributor-1@redhat.com
   - Most recent contribution: 2025-12-15 10:32:11
   - Total contributions: 4,521 lines across 89 files
   - Recent work: "Add AWS local zones support", "Fix VPC endpoint configuration"
   - GitHub: @contributor-1

2. contributor-2 (Contributor Two)
   ...

[Shows 5 contributors total]
```

## Arguments
- `platform-name`: The name of the platform to find relevant contributors (ex: aws).
- `count`: Number of contributors to return (optional, default: 2, max: 5).
- `--list`: List all available platforms without performing analysis.

## Performance Notes
- Analysis time scales with number of platform files
- Typical execution: 2-5 seconds for platforms with ~100 files
- Results are based on git history and may take longer on first run
- Using `git log` is more efficient than `git blame` for analyzing many files

## Troubleshooting

**No contributors found:**
- Check if the platform name is spelled correctly
- Some platforms may have been recently added with limited history
- Try checking OWNERS files: `cat OWNERS OWNERS_ALIASES`
- The command will automatically extend search to 24 months if no recent contributors found

**Git command fails:**
- Ensure you're in the installer repository root
- Verify git is installed and repository is initialized
- Check that you have access to git history (not a shallow clone)

**Unexpected results:**
- Bot accounts are filtered out automatically (openshift-merge-bot, openshift-ci-robot, dependabot)
- Results are based on file changes, not PR reviews
- Check pkg/types/${platform}/ manually if results seem off
- Some contributors may work on multiple platforms

**Unknown platform error:**
- Use `/platform-experts --list` to see all dynamically discovered platforms
- Platform names are case-sensitive and should be lowercase
- Some platforms may only exist in certain directories (e.g., `pkg/types/` or `pkg/asset/machines/`)
- The platform list is automatically generated from the codebase