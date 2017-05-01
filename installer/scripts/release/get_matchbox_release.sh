#!/bin/bash -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
source "$DIR/common.env.sh"

mkdir -p "$TMP_DIR" "$MATCHBOX_TMP_DIR" "$MATCHBOX_RELEASE_DIR"

set -x
curl -L -o "$TMP_DIR/matchbox.tar.gz" "$MATCHBOX_RELEASE_URL"
tar -xzf "$TMP_DIR/matchbox.tar.gz" \
    --strip-components=1 \
    -C "$MATCHBOX_RELEASE_DIR" \
    $MATCHBOX_ARCHIVE_TOP_DIR/LICENSE \
    $MATCHBOX_ARCHIVE_TOP_DIR/matchbox \
    $MATCHBOX_ARCHIVE_TOP_DIR/scripts \
    $MATCHBOX_ARCHIVE_TOP_DIR/contrib
