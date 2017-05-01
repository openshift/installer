#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"

GITHUB_API_URL="https://api.github.com/repos/coreos-inc/tectonic/releases"

echo "Creating new release on GitHub"
curl \
    --fail \
    -u "$GITHUB_CREDENTIALS" \
    -H "Content-Type: application/json" \
    -d "{\"tag_name\":\"$VERSION\",\"body\":\"Release tarball is available at $TECTONIC_RELEASE_TARBALL_URL.\"}" \
    $GITHUB_API_URL

