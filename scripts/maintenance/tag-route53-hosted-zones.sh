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
                    default tag value of tomorrow's date. This script tags resources
                    with 'expirationDate: some-date-string', where some-date-string
                    is replaced with either tomorrow's date or date-override.

EOF
}

force=
date_string=

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
      date_string="${2:-}"
      shift
    ;;
    *)
      echo "Flag '$1' is not supported." >&2
      exit 1
    ;;
  esac
  shift
done

if ! command -V jq >/dev/null || ! command -V aws >/dev/null; then
  echo "Missing required dependencies" >&2
  exit 1
fi

set -e

# Tag all Route53 hosted zones that do not already have a tag with the same keys,
# in this case 'expirationDate', with tomorrow's date as default, or
# with the --date-override value. Format YYYY-MM-DD.
if [ -z "$date_string" ]; then
  date_string="$(date -d tomorrow '+%Y-%m-%d')"
fi

tags="[{\"Key\":\"expirationDate\",\"Value\":\"$date_string\"}]"

echo "Tagging hosted zones with the following tags:"
echo "$tags"

if [ ! $force ]; then
  read -rp "Proceed tagging these resources? [y/N]: " yn
  if [ "$yn" != "y" ]; then
    echo "Aborting tagging and cleaning up." >&2
    exit 1
  fi
fi

private_zones=$(aws route53 list-hosted-zones | \
                jq ".HostedZones[] | select(.Config.PrivateZone == true) | .Id" | \
                sed "s@\"@@g")

for key in $(echo "$tags" | jq ".[].Key"); do
  for zone in $private_zones; do
    zone="${zone##*/}"
    is_not_tagged=$(aws route53 list-tags-for-resource \
    --resource-type hostedzone \
    --resource-id "$zone" | \
    jq ".ResourceTagSet | select(.Tags[]? | .Key == $key) | .ResourceId")
    if [ -z "$is_not_tagged" ]; then
      if aws route53 change-tags-for-resource \
      --resource-type hostedzone \
      --add-tags "$tags" \
      --resource-id "${zone##*/}"; then
        echo "Tagged hosted zone ${zone##*/}"
      else
        echo "Error tagging hosted zone ${zone##*/}" >&2
      fi
    fi
  done
done

set +e
