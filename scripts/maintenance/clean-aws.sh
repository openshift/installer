#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") deletes AWS resources tagged with tags specified in a tag file.
Requires that 'docker' and 'jq' are installed.

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

  --date-override   (optional) Date of the format YYYY-MM-DD that overrides the
                    default tag value of today's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either the following days' date or date-override.
                    Only use if --tag-file is not used.

  --dry-run         (optional) If set, grafiti will only do a dry run, i.e. not
                    delete any resources.

EOF
}

force=
version=
region=
config_file=
tag_file=
date_override=
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
      date_override="${2:-}"
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

if ! command -V docker >/dev/null || ! command -V jq >/dev/null; then
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

if [ -n "$tag_file" ] && [ -n "$date_override" ]; then
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
  tag_file="$(mktemp -p "$tmp_dir")"

  date_string="$(date "+%Y-%m-%d" -d "-1 day")\",\"$(date "+%Y-%-m-%-d" -d "-1 day")\",
  \"$(date "+%m-%-d-%-Y" -d "-1 day")\",\"$(date "+%-m-%-d-%-Y" -d "-1 day")\",\"$(date "+%d-%m-%-Y" -d "-1 day")\",
  \"$(date "+%d-%-m-%-Y" -d "-1 day")\",\"$(date +%m-%d-%Y)\",\"$(date +%d-%m-%Y)\",
  \"$(date +%d-%-m-%Y)\",\"$(date +%Y-%m-%d)\",\"$(date +%Y-%-m-%-d)"
  if [ -n "$date_override" ]; then
  	date_string="$date_override"
  fi

  cat <<EOF >"$tmp_dir/tag.json"
{"TagFilters":[{"Key":"expirationDate","Values":["${date_string}"]}]}
EOF
fi

echo "Deleting resources with the following tags:"
jq '.' "$tag_file"

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

trap 'docker stop grafiti-deleter; exit' EXIT

docker run -t --rm --name grafiti-deleter \
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
