#!/bin/bash -e

# USAGE:
#  export GITHUB_CREDENTIALS=username:personal-access-token
#  export TECTONIC_RELEASE_TARBALL_URL=url-of-tarball
#  export VERSION=w.x.y-tectonic.z

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"

GITHUB_API_URL="https://api.github.com/repos/coreos/tectonic-installer/releases"

echo "Creating new release on GitHub"
curl \
    --fail \
    -u "$GITHUB_CREDENTIALS" \
    -H "Content-Type: application/json" \
    -d "{\"tag_name\":\"$VERSION\",\"prerelease\":true,\"body\":\"Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL.\"}" \
    $GITHUB_API_URL

