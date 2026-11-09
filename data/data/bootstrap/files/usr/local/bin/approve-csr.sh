#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/bootstrap-service-record.sh

KUBECONFIG="${1}"

# Attributes shared by every kubelet client CSR that the machine-config-operator
# bootstrapper issues on behalf of a node joining the cluster. A request that
# carries all of them can only ever become a node identity, because the
# kube-apiserver-client-kubelet signer refuses any other subject.
readonly CLIENT_SIGNER="kubernetes.io/kube-apiserver-client-kubelet"
readonly CLIENT_REQUESTER="system:serviceaccount:openshift-machine-config-operator:node-bootstrapper"
readonly CLIENT_GROUPS="system:authenticated,system:serviceaccounts,system:serviceaccounts:openshift-machine-config-operator"
readonly CLIENT_USAGES="client auth,digital signature"
readonly NODE_ORG="system:nodes"

# subject_field prints the values a PEM encoded certificate request carries for
# one subject field, comma joined. Joining means a second value of the same
# field cannot hide behind the first, which matters for the organization: every
# organization becomes a group once the certificate is used.
# esc_ctrl makes sure one field holding a newline is not read as two fields.
#   1 - field name, e.g. CN or O   2 - PEM encoded certificate request
subject_field() {
    printf '%s' "${2}" |
        openssl req -noout -subject -nameopt sep_multiline,esc_ctrl,utf8 2>/dev/null |
        sed -n "s/^ *${1}=//p" |
        paste -sd, -
}

# is_node_client_request reports whether a PEM encoded certificate request
# carries exactly the subject a kubelet client certificate is allowed to have:
# a system:node:<name> common name, the system:nodes organization on its own,
# and no subject alternative names at all.
#   1 - PEM encoded certificate request
is_node_client_request() {
    local pem="${1}" text common_name

    # A request openssl cannot parse leaves every field below empty, which
    # would otherwise read as a pass. Reject it instead.
    text="$(printf '%s' "${pem}" | openssl req -noout -text 2>/dev/null)"
    [ -n "${text}" ] || return 1

    case "${text}" in
        *"X509v3 Subject Alternative Name"*) return 1 ;;
        *) ;;
    esac

    [ "$(subject_field O "${pem}")" = "${NODE_ORG}" ] || return 1

    common_name="$(subject_field CN "${pem}")"
    case "${common_name}" in
        # esc_ctrl turns a control character into a backslash escape. A node
        # name is a DNS name and never holds one, and the prefix test alone
        # would pass a name whose tail is a newline.
        *\\*) return 1 ;;
        "system:node:"?*) ;;
        *) return 1 ;;
    esac
}

echo "Approving node client CSR requests until bootstrapping is complete..."
# Names of the requests already reported as left pending, so that each one is
# described once rather than on every pass.
logged=" "
while [ ! -f /opt/openshift/.bootkube.done ]
do
    # Approve only the kubelet client requests that match every expected
    # attribute. Anything else is left pending for cluster-machine-approver
    # rather than approved here.
    while IFS=$'\t' read -r name attributes signer requester request
    do
        [ -n "${name}" ] || continue
        [ "${request}" = "<none>" ] && request=""

        if [ "${attributes}" = "client" ] &&
            is_node_client_request "$(printf '%s' "${request}" | base64 --decode 2>/dev/null)"; then
            echo "Approving node client CSR ${name}"
            oc --kubeconfig="$KUBECONFIG" adm certificate approve "${name}"
            continue
        fi

        case "${logged}" in
            *" ${name} "*) continue ;;
            *) logged="${logged}${name} " ;;
        esac
        echo "Leaving CSR ${name} pending: signer=${signer} requester=${requester}"
    done <<< "$(oc --kubeconfig="$KUBECONFIG" get csr -o json 2>/dev/null | jq -r \
        --arg signer "${CLIENT_SIGNER}" \
        --arg requester "${CLIENT_REQUESTER}" \
        --arg groups "${CLIENT_GROUPS}" \
        --arg usages "${CLIENT_USAGES}" '
        # Tab is an IFS whitespace character, so an empty field would be
        # swallowed by read and shift every field after it. Never emit one.
        def or_none: if . == null or . == "" then "<none>" else . end;

        # A CSR is pending when nothing has yet added an Approved or Denied
        # condition to it.
        .items[]?
        | select(
            ((.status.conditions // []) | map(.type)) as $types
            | ($types | index("Approved")) == null and ($types | index("Denied")) == null
          )
        | [
            .metadata.name,
            (if .spec.signerName == $signer
                and .spec.username == $requester
                and (((.spec.groups // []) | sort | join(",")) == $groups)
                and (((.spec.usages // []) | sort | join(",")) == $usages)
             then "client" else "other" end),
            (.spec.signerName | or_none),
            (.spec.username | or_none),
            (.spec.request | or_none)
          ]
        | @tsv')"

    sleep 20
done
