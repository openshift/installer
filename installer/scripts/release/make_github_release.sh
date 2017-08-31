#!/bin/bash -e

# USAGE:
#  export GITHUB_CREDENTIALS=username:personal-access-token
#  export TECTONIC_RELEASE_TARBALL_URL=url-of-tarball
#  export VERSION=w.x.y-tectonic.z
#  export PRE_RELEASE=true

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# shellcheck disable=SC1090
source "$DIR/common.env.sh"

GITHUB_API_URL="https://api.github.com/repos/coreos/tectonic-installer/releases"

# shellcheck disable=SC2028
echo "Creating new release on GitHub ${#GITHUB_CREDENTIALS} \n\n {\"tag_name\":\"$VERSION\", \"name\": \"$VERSION\", \"prerelease\":$PRE_RELEASE,\"body\":\"Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL.\"}"
curl \
    --fail \
    -u "$GITHUB_CREDENTIALS" \
    -H "Content-Type: application/json" \
    -d "{\"tag_name\":\"$VERSION\",\"name\": \"$VERSION\",\"prerelease\":$PRE_RELEASE,\"body\":\"Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL.\"}" \
    $GITHUB_API_URL
