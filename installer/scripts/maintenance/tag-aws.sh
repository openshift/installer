#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") tags AWS resources with 'expirationDate: some-date-string',
defaulting to the following days' date, and excludes all resources tagged with
tag keys/values specified in an 'exclude' file. Requires that both 'jq' and the
AWS CLI are installed.

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

  --exclude-file    A file containing a JSON array of Key/Value pair objects.

  --start-hour      Integer hour to start looking at CloudTrail logs. Defaults to 8.

  --end-hour        Integer hour to end looking at CloudTrail logs. Defaults to 1.

  --date-override   (optional) Date of the format YYYY-MM-DD that overrides the
                    default tag value of today's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either the following days' date or date-override.

  --workspace-dir   (optional) Parent directory for a temporary directory. /tmp is
                    used by default.

  --dry-run         (optional) If set, grafiti will only do a dry run, i.e. not tag
                    any resources.

EOF
}

force=
version=
region=
config_file=
exclude_file=
date_override=
workspace=
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

if [ -z "$start_hour" ] || [ -z "$end_hour" ]; then
  echo "Start hour and end hour must be specified."
  exit 1
fi

if [ -z "$region" ]; then
  echo "Must provide an AWS region, set the AWS_REGION, or set a region in your ~/.aws/config}"
  exit 1
fi

set -e

# Tag all resources present in CloudTrail over the specified time period with the
# following day's date as default, or with the DATE_VALUE_OVERRIDE value.
# Format YYYY-MM-DD.
tmp_dir="/tmp/config"
if [ -n "$workspace" ]; then
  tmp_dir="$(readlink -m "${workspace}/config")"
fi
mkdir -p "$tmp_dir"
trap 'rm -rf "$tmp_dir"; exit' EXIT

date_string='now|strftime(\"%Y-%m-%d\")'
if [ -n "$date_override" ]; then
	date_string='\"'"${date_override}"'\"'
fi

# Configure grafiti to tag all resources created between START_HOUR and END_HOUR's
# ago
if [ -z "$config_file" ]; then
  config_file="$(mktemp -p "$tmp_dir" --suffix=.toml)"
  cat <<EOF > "$config_file"
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
if [ -z "$exclude_file" ]; then
  exclude_file="$(mktemp -p "$tmp_dir")"
  echo '{"TagFilters":[{"Key":"expirationDate","Values":[]}]}' > "$exclude_file"
fi

echo "Tagging resources with the following configuration:"
cat "$config_file"

if [ -n "$dry_run" ]; then
  echo "Dry run flag set. Not tagging any resources."
fi

if [ ! $force ]; then
  read -rp "Proceed tagging these resources? [y/N]: " yn
  if [ "$yn" != "y" ]; then
    echo "Aborting tagging and cleaning up."
    exit 1
  fi
fi

trap 'docker stop grafiti-tagger && docker rm grafiti-tagger; exit' EXIT

docker run -t --rm --name grafiti-tagger \
	-v "$tmp_dir":/tmp/config:z \
	-e AWS_ACCESS_KEY_ID="$AWS_ACCESS_KEY_ID" \
	-e AWS_SECRET_ACCESS_KEY="$AWS_SECRET_ACCESS_KEY" \
  -e AWS_REGION="$region" \
	-e CONFIG_FILE="/tmp/config/$(basename "$config_file")" \
	-e TAG_FILE="/tmp/config/$(basename "$exclude_file")" \
	quay.io/coreos/grafiti:"${version}" \
  bash -c "grafiti --config \"\$CONFIG_FILE\" parse | \
	grafiti --config \"\$CONFIG_FILE\" filter --ignore-file \"\$TAG_FILE\" | \
	grafiti $dry_run --config \"\$CONFIG_FILE\" tag"

set +e
