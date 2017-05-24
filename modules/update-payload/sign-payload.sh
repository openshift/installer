#!/bin/bash

set -e -o pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# GPG binary may be named gpg2 or gpg
if ! gpg=$(type -P gpg2 || type -P gpg); then
    echo "Cannot find gpg2 or gpg binary!" >&2
    exit 1
fi

while ! serial=$(${gpg} --card-status 2>/dev/null \
        | grep ^Serial | sed -e 's/^.*: //'); do
    echo  "Please insert CoreOS App Signing Yubikey..."
    sleep 5
done

case "${serial}" in
04612233)
    echo "Found Tectonic Yubikey"
    subkey="BEDDBA18"
    ;;
04149341)
    echo "Found Berlin Yubikey" # aka rkt
    subkey="3F1B2C87"
    ;;
*)
    echo "Unknown Yubikey ${serial}" >&2
    exit 1
esac

${gpg} -u "${subkey}!" --detach-sign "${DIR}/payload.json"
${gpg} --verify "${DIR}/payload.json.sig"

# TODO: In the future when we rotate away from the above keys we must
# sign payloads with both the old and new keys. To do so, sign the
# payload twice, writing each signature to different files, and then
# concatenating them together.
