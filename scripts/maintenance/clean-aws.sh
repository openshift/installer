#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") deletes AWS resources tagged with tags specified in a tag file.
Requires that 'podman' and 'jq' are installed.

AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environmental variables must be set.

Options:

  --force           Override user input prompts. Useful for automation.

  --grafiti-version Either the semver release version, ex. v0.1.1, or sha commit
                    hash of a grafiti image hosted in quay.io.

  --aws-region      The AWS region you wish to query for taggable resources. This
                    flag is optional if AWS_REGION is set.  You can also set a
                    default region for the default profile in your ~/.aws
                    configuration files, although for this you must have the 'aws'
                    command installed).

  --config-file     A grafiti configuration file. See an example at
                    https://github.com/coreos/grafiti/blob/master/config.toml.

  --tag-file        A file containing a TagFilter list. See the AWS Resource Group
                    Tagging API 'TagFilter' documentation for file structure.

  --date-override   (optional) Date of the format YYYY-MM-DD to delete resources
                    tagged with 'expirationDate: some-date-string'.  By default,
                    this script deletes resources which expired yesterday or
                    today.  Not compatible with --tag-file.

  --dry-run         (optional) If set, grafiti will only do a dry run, i.e. not
                    delete any resources.

EOF
}

force=
version=
region=
config_file=
tag_file=
date_string=
dry_run=

while [ $# -gt 0 ]; do
  case $1 in
    --help)
      usage
      exit
    ;;
    --force)
      force=true
    ;;
    --grafiti-version)
      version="${2:-}"
      shift
    ;;
    --aws-region)
      region="${2:-}"
      shift
    ;;
    --config-file)
      config_file="${2:-}"
      shift
    ;;
    --tag-file)
      tag_file="${2:-}"
      shift
    ;;
    --date-override)
      date_string="[\"${2:-}\"]"
      shift
    ;;
    --dry-run)
      dry_run="$1"
    ;;
    *)
      echo "Flag '$1' is not supported." >&2
      exit 1
    ;;
  esac
  shift
done

if ! command -V podman >/dev/null || ! command -V jq >/dev/null; then
  echo "Missing required dependencies" >&2
  exit 1
fi

if [ -z "$region" ]; then
  if [ -n "$AWS_REGION" ]; then
    region="${AWS_REGION:-}"
  elif ! command -V aws >/dev/null; then
    echo "Without the 'aws' command, you must set either --aws-region or \$AWS_REGION" >&2
    exit 1
  else
    region="$(aws configure get region)"
    if [ -z "$region" ]; then
      echo "Must provide an AWS region, set the AWS_REGION, or set a region in your ~/.aws/config" >&2
      exit 1
    fi
  fi
fi

if [ -z "$version" ]; then
  echo "Grafiti image version required." >&2
  exit 1
fi

if [ -n "$tag_file" ] && [ -n "$date_string" ]; then
  echo "Cannot use both --tag-file and --date-override flags simultaneously." >&2
  exit 1
fi

set -e

tmp_dir="$(readlink -m "$(mktemp -d clean-aws-XXXXXXXXXX)")"
mkdir -p "$tmp_dir"
trap 'rm -rf "$tmp_dir"; exit' EXIT

if [ -n "$config_file" ]; then
  cat "$config_file" >"$tmp_dir/config.toml"
else
  echo "maxNumRequestRetries = 11" >"$tmp_dir/config.toml"
fi

if [ -n "$tag_file" ]; then
  cat "$tag_file" >"$tmp_dir/tag.json"
else
  if [ -z "$date_string" ]; then
    date_string="$(jq --null-input '[["%Y-%m-%d", "%Y-%-m-%-d", "%m-%d-%Y", "%m-%-d-%-Y", "%-m-%-d-%-Y", "%d-%m-%Y", "%d-%-m-%-Y"][] | . as $format | [now, now - 24*60*60][] | strftime($format)]')"
  fi

  cat <<EOF >"$tmp_dir/tag.json"
{"TagFilters":[{"Key":"expirationDate","Values":${date_string}}]}
EOF
fi

echo "Deleting resources with the following tags:"
jq '.' "$tmp_dir/tag.json"

if [ -n "$dry_run" ]; then
  echo "Dry run flag set. Not deleting any resources."
fi

if [ ! $force ]; then
  read -rp "Proceed deleting these resources? [y/N]: " yn
  if [ "$yn" != "y" ]; then
    echo "Aborting deletion and cleaning up." >&2
    exit 1
  fi
fi

trap 'podman stop grafiti-deleter; exit' EXIT

podman run -t --rm --name grafiti-deleter \
	-v "$tmp_dir":/tmp/config:z \
	-e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
	-e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
  -e AWS_SESSION_TOKEN="$AWS_SESSION_TOKEN" \
  -e AWS_REGION="$region" \
	-e CONFIG_FILE="/tmp/config/config.toml" \
	-e TAG_FILE="/tmp/config/tag.json" \
	quay.io/coreos/grafiti:"${version}" \
	bash -c "grafiti $dry_run --config \"\$CONFIG_FILE\" --ignore-errors delete --all-deps --delete-file \"\$TAG_FILE\""

set +e
