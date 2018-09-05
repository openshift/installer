#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") tags AWS resources with 'expirationDate: some-date-string',
defaulting to tomorrow's date, and excludes all resources tagged with
tag keys/values specified in an 'exclude' file. Requires that 'podman' is
installed.

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

  --exclude-file    A file containing a JSON array of Key/Value pair objects.

  --start-hour      Integer hour to start looking at CloudTrail logs. Defaults to 8.

  --end-hour        Integer hour to end looking at CloudTrail logs. Defaults to 1.

  --date-override   (optional) Date of the format YYYY-MM-DD that overrides the
                    default tag value of tomorrow's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either tomorrow's date or date-override.

  --dry-run         (optional) If set, grafiti will only do a dry run, i.e. not tag
                    any resources.

EOF
}

force=
version=
region=
config_file=
exclude_file=
date_string=
start_hour=8
end_hour=1
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
    --exclude-file)
      exclude_file="${2:-}"
      shift
    ;;
    --start-hour)
      start_hour="${2:-}"
      shift
    ;;
    --end-hour)
      end_hour="${2:-}"
      shift
    ;;
    --date-override)
      date_string="\\\"${2:-}\\\""
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

if ! command -V podman >/dev/null; then
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

if [ -z "$start_hour" ] || [ -z "$end_hour" ]; then
  echo "Start hour and end hour must be specified." >&2
  exit 1
fi

set -e

# Tag all resources present in CloudTrail over the specified time period with the
# today's date as default, or with the --date-override value.
# Format YYYY-MM-DD.
tmp_dir="$(readlink -m "$(mktemp -d tag-aws-XXXXXXXXXX)")"
trap 'rm -rf "$tmp_dir"; exit' EXIT

if [ -z "$date_string" ]; then
	date_string='(now + 24*60*60|strftime(\"%Y-%m-%d\"))'
fi

# Configure grafiti to tag all resources created between START_HOUR and END_HOUR's
# ago
if [ -n "$config_file" ]; then
  cat "$config_file" >"$tmp_dir/config.toml"
else
  cat <<EOF >"$tmp_dir/config.toml"
endHour = -${end_hour}
startHour = -${start_hour}
includeEvent = false
tagPatterns = [
	"{expirationDate: ${date_string}}"
]
EOF
fi

# Exclusion file prevents tagging of resources that already have tags with the key
# "expirationDate"
if [ -n "$exclude_file" ]; then
  cat "$exclude_file" >"$tmp_dir/exclude"
else
  echo '{"TagFilters":[{"Key":"expirationDate","Values":[]}]}' >"$tmp_dir/exclude"
fi

echo "Tagging resources with the following configuration:"
cat "$tmp_dir/config.toml"

if [ -n "$dry_run" ]; then
  echo "Dry run flag set. Not tagging any resources."
fi

if [ ! $force ]; then
  read -rp "Proceed tagging these resources? [y/N]: " yn
  if [ "$yn" != "y" ]; then
    echo "Aborting tagging and cleaning up." >&2
    exit 1
  fi
fi

trap 'podman stop grafiti-tagger; exit' EXIT

podman run -t --rm --name grafiti-tagger \
	-v "$tmp_dir":/tmp/config:z \
	-e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
	-e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
  -e AWS_SESSION_TOKEN="$AWS_SESSION_TOKEN" \
  -e AWS_REGION="$region" \
	-e CONFIG_FILE="/tmp/config/config.toml" \
	-e TAG_FILE="/tmp/config/exclude" \
	quay.io/coreos/grafiti:"${version}" \
  bash -c "grafiti --config \"\$CONFIG_FILE\" parse | \
	grafiti --config \"\$CONFIG_FILE\" filter --ignore-file \"\$TAG_FILE\" | \
	grafiti $dry_run --config \"\$CONFIG_FILE\" tag"

set +e
