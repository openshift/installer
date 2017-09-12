#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") deletes AWS resources tagged with tags specified in a tag file.

AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environmental variables must be set.

Options:

  --force           Override user input prompts. Useful for automation.

  --grafiti-version Either the semver release version, ex. v0.1.1, or sha commit
                    hash of a grafiti image hosted in quay.io.

  --aws-region      The AWS region you wish to query for taggable resources. This
                    flag is optional if AWS_REGION is set. AWS_REGION overrides
                    values passed in by this flag.

  --config-file     A grafiti configuration file. See an example at
                    https://github.com/coreos/grafiti/blob/master/config.toml.

  --tag-file        A file containing a TagFilter list. See the AWS Resource Group
                    Tagging API 'TagFilter' documentation for file structure.

  --date-override   (optional) Date of the format YYYY-MM-DD that overrides the
                    default tag value of today's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either the following days' date or date-override.
                    Only use if --tag-file is not used.

  --workspace-dir   (optional) Parent directory for a temporary directory. /tmp is
                    used by default.

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
workspace=
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
    --workspace-dir)
      workspace="${2:-}"
      shift
    ;;
    --dry-run)
      dry_run="$1"
    ;;
    *)
      echo "Flag '$2' is not supported."
      exit
    ;;
  esac
  shift
done

if [ -n "$AWS_REGION" ]; then
  region="${AWS_REGION:-}"
fi

if [ -z "$version" ]; then
  echo "Grafiti image version required."
  exit 1
fi

if [ -z "$region" ]; then
  echo "Must provide an AWS region, set the AWS_REGION, or set a region in your ~/.aws/config}"
  exit 1
fi

if [ -n "$tag_file" ] && [ -n "$date_override" ]; then
  echo "Cannot use both --tag-file and --date-override flags simultaneously."
  exit 1
fi

set -e

tmp_dir="/tmp/config"
if [ -n "$workspace" ]; then
  tmp_dir="$(readlink -m "${workspace}/config")"
fi
mkdir -p "$tmp_dir"
trap 'rm -rf "$tmp_dir"; exit' EXIT

if [ -z "$config_file" ]; then
  config_file="$(mktemp -p "$tmp_dir" --suffix=.toml)"
  echo "maxNumRequestRetries = 11" > "$config_file"
fi

if [ -z "$tag_file" ]; then
  tag_file="$(mktemp -p "$tmp_dir")"

  date_string="$(date "+%Y-%m-%d" -d "-1 day")\",\"$(date "+%Y-%-m-%-d" -d "-1 day")\",
  \"$(date "+%m-%-d-%-Y" -d "-1 day")\",\"$(date "+%-m-%-d-%-Y" -d "-1 day")\",\"$(date "+%d-%m-%-Y" -d "-1 day")\",
  \"$(date "+%d-%-m-%-Y" -d "-1 day")\",\"$(date +%m-%d-%Y)\",\"$(date +%d-%m-%Y)\",
  \"$(date +%d-%-m-%Y)\",\"$(date +%Y-%m-%d)\",\"$(date +%Y-%-m-%-d)"
  if [ -n "$date_override" ]; then
  	date_string="$date_override"
  fi

  cat <<EOF > "$tag_file"
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
    echo "Aborting deletion and cleaning up."
    exit 1
  fi
fi

trap 'docker stop grafiti-deleter && docker rm grafiti-deleter; exit' EXIT

docker run -t --rm --name grafiti-deleter \
	-v "$tmp_dir":/tmp/config:z \
	-e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
	-e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
  -e AWS_REGION="$region" \
	-e CONFIG_FILE="/tmp/config/$(basename "$config_file")" \
	-e TAG_FILE="/tmp/config/$(basename "$tag_file")" \
	quay.io/coreos/grafiti:"${version}" \
	bash -c "grafiti $dry_run --config \"\$CONFIG_FILE\" --ignore-errors delete --all-deps --delete-file \"\$TAG_FILE\""

set +e
