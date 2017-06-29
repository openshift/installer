#!/usr/bin/env bash

usage() {
  cat <<EOF

$(basename "$0") tags AWS Route53 Hosted Zones with an 'expirationDate' of tomorrow.
Requires that both 'jq' and the AWS CLI are installed.

Either the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environmental variables
must be set, or ~/.aws/credentials must contain valid AWS credentials.

Options:

  --force           Override user input prompts. Useful for automation.

  --date-override   (optional) Date of the format YYYY-MM-DD that overrides the
                    default tag value of today's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either the following days' date or date-override.

EOF
}

force=
date_override=

while [ $# -gt 0 ]; do
  case $1 in
    --help)
      usage
      exit
    ;;
    --force)
      force=true
    ;;
    --date-override)
      date_override="${2:-}"
      shift
    ;;
    *)
      echo "Flag '$2' is not supported."
      exit
    ;;
  esac
  shift
done

if ! command -v jq > /dev/null || ! command -v aws > /dev/null; then
  "Dependencies not installed."
  exit 1
fi

set -e

# Tag all Route53 hosted zones that do not already have a tag with the same keys,
# in this case 'expirationDate', with today's date as default, or
# with the DATE_VALUE_OVERRIDE value. Format YYYY-MM-DD.
date_string="$(date "+%Y-%m-%d")"
if [ -n "$date_override" ]; then
	date_string="${date_override}"
fi

tags="[{\"Key\":\"expirationDate\",\"Value\":\"$date_string\"}]"

echo "Tagging hosted zones with the following tags:"
echo "$tags"

if [ ! $force ]; then
  read -rp "Proceed tagging these resources? [y/N]: " yn
  if [ "$yn" != "y" ]; then
    echo "Aborting tagging and cleaning up."
    exit 1
  fi
fi

private_zones=$(aws route53 list-hosted-zones | \
                jq ".HostedZones[] | select(.Config.PrivateZone == true) | .Id" | \
                sed "s@\"@@g")

for key in $(echo -e "$tags" | jq ".[].Key"); do
  for zone in $private_zones; do
    zone="${zone##*/}"
    is_not_tagged=$(aws route53 list-tags-for-resource \
    --resource-type hostedzone \
    --resource-id "$zone" | \
    jq ".ResourceTagSet | select(.Tags[]? | .Key == $key) | .ResourceId")
    if [ -z "$is_not_tagged" ]; then
      if aws route53 change-tags-for-resource \
      --resource-type hostedzone \
      --add-tags "$(echo -e "$tags")" \
      --resource-id "${zone##*/}"; then
        echo "Tagged hosted zone ${zone##*/}"
      else
        echo "Error tagging hosted zone ${zone##*/}"
      fi
    fi
  done
done

set +e
