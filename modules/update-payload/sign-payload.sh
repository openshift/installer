#!/bin/bash

set -e -o pipefail

TEST_GPG_ID="AEDE60EDD314382875F4B2C279C51CC59F11DE78"
TEST_GPG_KEY="
-----BEGIN PGP PRIVATE KEY BLOCK-----

lQOYBFkBC5cBCAC+L2HdvKPHVglOFuBpE3rKEAy7K3ccRvd386ESgn4SnYZBkrGP
57zwaZM3tbWusdttALwJCRkHsO0knFgfF3Tavso4ekrdMT+W97TWKbKYjS38aq6E
FbY/cX4SDdeY8x95VK5hk7a15lXLPGaX5u3au32kUMb9D03FCGMQWphV+SFbMW/v
GxaDprJNTcO4aj4v2IDhu2/UJAyf/CKJlEYRpeDel18jqlVzqdf7zRpFA+RuTPWa
6GRXy821wtzo33GdnZJu6SDlY7A/HyIMqGARge4i3zL+eywAIXd8JrQ9evCK/V0U
41X7jJY8tY/MqNrxsuwuiBU9OHWDm1WDW/JvABEBAAEAB/9bCvmxmVVZ3nDz5fWl
t4KHaml9zE/nxH1A+er4nmlV5fzYNS41Mn17JT8pYx5HM7NUGg7p0GYgDW1ookN5
FtSExuKtaLWf76x+S7RQ9YMeji8eb799UZt+AYmVgPTdrj056vTTl0/XAh29/fsq
6oygsjJyT40CpanYEazCrmvQc3A+EJqNE36XJLB8GERkBKfxp71dusOpklgIf0Q2
rRsiBqbcs/KIpdA6qWLBpHAO1FIIjTlCsfw04NFbLVKenP2wYAa53QRopDfIaIrT
14XQe5uRyA/AkcrUBlHHYIRF9Eqs9bul+RhOPATyT7xqVqiaJ+rTcOjdcN2OXe+j
9BtJBADUNlywu5eu1tmX9mKtmfkQkgGwMro7OBDXuAasSuiyvFlpbaXbFpMnjzv1
7AMpJwSMxER/R9kiz/jEMBynw9RN+5JpFIMWTs1gdjIJ01EZ1glgPEp5Z7h94bWw
OKBM2uf0IjBeS2HT7003q+fVcIEGjYu5vEt7/S+TLX8ySDThswQA5W16z9E+U32L
hmLkEJ//BZC4ef9SRU9pskzguMoDkKnxYwngo6ZF77fhITw52VQ1ZxylCRrmugU/
7c2QnkLF7LKaMf7OVrkESX8kT4U3gKgyqhsy6o8zMYscBGI+m57riClg+lm9+G/Y
K31/+5sVgMdGbuh1kYUV+Px+ev/p9lUEAMLQ1JVUGHUTIS0mGKha6bMnObOvHhsb
myQJja2rWuY3YVtOT7wDiGNcBP/dkNjrjHIElV+LhSeq3srw6w+gPs4smMdYrGPx
F4yhfJCeZIT2gpJZyCkzhnD2sUX11lor1DRd7ZgoLQhx/MSRqUrBMY4y5ldAZXnb
k91rThC4fTx4Rv+0IlRDTyBUZXN0IEtleSA8dGVzdGtleUBleGFtcGxlLmNvbT6J
AU4EEwEIADgWIQSu3mDt0xQ4KHX0ssJ5xRzFnxHeeAUCWQELlwIbAwULCQgHAgYV
CAkKCwIEFgIDAQIeAQIXgAAKCRB5xRzFnxHeeCk9CACaM0DxZO0W3La1GBj3rkkL
j24ylxIvZbp5hYjK98M2fxshQ+i2R+sqt0SApcyIFOQ27vrvGPTV+lHAKFIugIaE
km02aeEY8zDZxVQHRjytPe9vYEhtdsk9bSeiyTKXEpFakF9CiVp/jqOYz8q3Rc0s
D52wq+GZNmxL+12cFXAN1rqEWygHz2KlmPriC6cZC9WlPvlC/YKFIU85/tRv8QL5
yBA4rauCi4z1fjbc44H7/yWplk3Q0r2awLjNNyhzeQRHzy7kG1VSVax/fzamXpWT
NFUGWIFj+uedImuj2Fzd1X4oK9eT5WyjdjHQBPFHOe7wEO9ZEQ5yNJSHSuuDSGCW
=o+gW
-----END PGP PRIVATE KEY BLOCK-----
"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# GPG binary may be named gpg2 or gpg
if ! gpg=$(type -P gpg2 || type -P gpg); then
  echo "Cannot find gpg2 or gpg binary!" >&2
  exit 1
fi

function sign_with_testkey() {
  echo "Signing with test key"

  GNUPGHOME=$(mktemp -d)
  trap 'rm -rf ${GNUPGHOME}' EXIT
  export GNUPGHOME

  ${gpg} --quiet --import <<<"${TEST_GPG_KEY}"
  ${gpg} -u "${TEST_GPG_ID}" --detach-sign "${DIR}/payload.json"
  ${gpg} --verify "${DIR}/payload.json.sig"
}

function sign_with_yubikey() {
  echo "Signing with yubikey"

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
}

function copy_existing_sig() {
  if [ -f "${DIR}/payload.json.sig" ]; then
      cp "${DIR}/payload.json.sig" "${DIR}/payload.json.sig.old"
  fi
}

copy_existing_sig

if [[ "${SIGN_WITH_TEST_KEY}" == "true" ]]; then
  sign_with_testkey
else
  sign_with_yubikey
fi
