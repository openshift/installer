#!/usr/bin/env bash

# shellcheck disable=SC1091  # using path on bootstrap machine
. /usr/local/bin/bootstrap-service-record.sh

KUBECONFIG="${1}"

# Attributes shared by every kubelet client CSR that the machine-config-operator
# bootstrapper issues on behalf of a node joining the cluster.
readonly CLIENT_SIGNER="kubernetes.io/kube-apiserver-client-kubelet"
readonly CLIENT_REQUESTER="system:serviceaccount:openshift-machine-config-operator:node-bootstrapper"
readonly CLIENT_GROUPS="system:authenticated,system:serviceaccounts,system:serviceaccounts:openshift-machine-config-operator"
readonly CLIENT_USAGES="client auth,digital signature"

# Attributes shared by every kubelet serving CSR, which a node requests for
# itself once it holds the client certificate issued above.
readonly SERVING_SIGNER="kubernetes.io/kubelet-serving"
readonly SERVING_GROUPS="system:authenticated,system:nodes"
readonly SERVING_USAGES="digital signature,server auth"

readonly NODE_ORG="system:nodes"

# Networks a kubelet cannot be serving on, and what to call each one in the
# report. Written as CIDRs so that ipcalc decides membership and nothing here
# has to parse an address.
readonly REFUSED_NETWORKS="
0.0.0.0/32 unspecified
::/128 unspecified
127.0.0.0/8 loopback
::1/128 loopback
169.254.0.0/16 link-local
fe80::/10 link-local
224.0.0.0/4 multicast
ff00::/8 multicast
"

# ip_in_cidr reports whether an address falls inside a CIDR. The address is
# masked to the prefix the CIDR carries and compared with the network the CIDR
# itself describes, so a prefix that does not land on a nibble boundary
# behaves like one that does, and an address openssl printed expanded matches
# a CIDR the installer wrote compressed. ipcalc refuses a prefix the family
# cannot hold, so an address never matches a CIDR of the other family.
#   1 - address   2 - CIDR
ip_in_cidr() {
    local cidr="${2}" prefix masked network
    case "${cidr}" in */*) ;; *) return 1 ;; esac
    prefix="${cidr#*/}"
    masked="$(ipcalc --network "${1}/${prefix}" 2>/dev/null)" || return 1
    network="$(ipcalc --network "${cidr}" 2>/dev/null)" || return 1
    [ -n "${masked}" ] && [ "${masked}" = "${network}" ]
}

# subnet_check_works reports whether ip_in_cidr answers correctly, for an
# address inside the network and one outside it, in both families. That is not
# the same question as whether ipcalc is installed: more than one program goes
# by that name, and they do not agree on which options exist. Every refusal
# below rests on a call that failed, so a check that always fails is a check
# that always approves.
subnet_check_works() {
    ip_in_cidr 10.1.2.3 10.0.0.0/8 &&
        ! ip_in_cidr 192.168.1.1 10.0.0.0/8 &&
        ip_in_cidr fd00::1 fd00::/8 &&
        ! ip_in_cidr 2600:1f14::a3 fd00::/8
}

# Every address check below is the subnet check deciding whether an address
# falls in a network, so without a working one nothing could be judged and an
# address a node cannot hold would read like one it can. Stop rather than
# approve on a guess.
if ! subnet_check_works; then
    echo "approve-csr: no working subnet check, refusing to approve any CSR" >&2
    exit 1
fi

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

# request_sans prints the subject alternative names a PEM encoded certificate
# request carries, or nothing at all when it carries none.
#   1 - PEM encoded certificate request
request_sans() {
    printf '%s' "${1}" |
        openssl req -noout -text 2>/dev/null |
        sed -n '/X509v3 Subject Alternative Name/{n;p}'
}

# san_field prints the subject alternative names of one type carried by a PEM
# encoded certificate request, comma joined.
#   1 - SAN type as openssl prints it: DNS, "IP Address", email or URI
#   2 - PEM encoded certificate request
san_field() {
    request_sans "${2}" |
        tr ',' '\n' |
        sed -n "s/^ *${1}://p" |
        paste -sd, -
}

# ip_address_class prints the class of an address that a kubelet cannot be
# serving on, or nothing at all when the address is an ordinary unicast one.
# An address ipcalc cannot read is called malformed rather than left to pass
# as unicast, since nothing can be said about it.
#   1 - address as openssl prints it
ip_address_class() {
    local address="${1}" canonical network class

    # openssl prints IPv6 expanded and in upper case, so ::1 arrives as
    # 0:0:0:0:0:0:0:1. ipcalc reads that and prints the canonical form, which
    # is what the networks below are matched against.
    canonical="$(ipcalc --address "${address}" 2>/dev/null)"
    canonical="${canonical#ADDRESS=}"
    if [ -z "${canonical}" ]; then
        printf 'malformed'
        return 0
    fi

    while read -r network class; do
        [ -n "${network}" ] || continue
        if ip_in_cidr "${canonical}" "${network}"; then
            printf '%s' "${class}"
            return 0
        fi
    done <<< "${REFUSED_NETWORKS}"
}

# node_common_name validates the parts of the subject that both kinds of
# kubelet request must carry and prints the common name. It fails when the
# request does not parse, carries anything but system:nodes on its own as the
# organization, names more than one node, or does not name a node.
#   1 - PEM encoded certificate request
node_common_name() {
    local pem="${1}" common_name

    # A request openssl cannot parse leaves every field below empty, which
    # would otherwise read as a pass. Reject it instead.
    [ -n "$(printf '%s' "${pem}" | openssl req -noout -subject 2>/dev/null)" ] || return 1

    [ "$(subject_field O "${pem}")" = "${NODE_ORG}" ] || return 1

    common_name="$(subject_field CN "${pem}")"
    case "${common_name}" in
        # subject_field joins repeated values with a comma, so a comma here
        # means the subject carries more than one common name and there is no
        # single name to act on. A node name never contains one, so the value
        # is refused rather than guessed at: the prefix below would otherwise
        # pass on the first name while a second travelled alongside it.
        *,*) return 1 ;;
        # esc_ctrl turns a control character into a backslash escape. A node
        # name is a DNS name and never holds one, and the prefix test alone
        # would pass a name whose tail is a newline.
        *\\*) return 1 ;;
        "system:node:"?*) printf '%s' "${common_name}" ;;
        *) return 1 ;;
    esac
}

# is_node_client_request reports whether a PEM encoded certificate request
# carries exactly the subject a kubelet client certificate is allowed to have.
#   1 - PEM encoded certificate request
is_node_client_request() {
    local pem="${1}" common_name

    # Declared separately: `local name=$(...)` would return the status of
    # local rather than of the command substitution.
    common_name="$(node_common_name "${pem}")" || return 1

    # A client certificate authenticates the kubelet, it never names a host,
    # so it has no use for subject alternative names of any kind.
    [ -z "$(request_sans "${pem}")" ] || return 1
}

# is_node_serving_request reports whether a PEM encoded certificate request
# carries exactly the subject a kubelet serving certificate is allowed to have,
# for the node that asked for it.
#   1 - PEM encoded certificate request   2 - requesting username
#   3 - CSR name, used only to report an address worth looking at
is_node_serving_request() {
    local pem="${1}" requester="${2}" name="${3}" common_name address class

    common_name="$(node_common_name "${pem}")" || return 1

    # A node may only ask for a serving certificate that names itself.
    [ "${common_name}" = "${requester}" ] || return 1

    # A serving certificate legitimately names hosts and addresses. The other
    # kinds of subject alternative name carry identities a kubelet never needs.
    [ -z "$(san_field email "${pem}")$(san_field URI "${pem}")" ] || return 1

    # A kubelet serves on the addresses of the node it runs on, so an address
    # that cannot be one of those names something other than that node.
    for address in $(san_field 'IP Address' "${pem}" | tr ',' ' '); do
        class="$(ip_address_class "${address}")"
        case "${class}" in
            '') ;;
            link-local)
                # Reported rather than refused: the cluster network itself uses
                # addresses out of 169.254.0.0/16, and it is not yet confirmed
                # that none of them can reach a node's address list.
                echo "CSR ${name} carries a ${class} address ${address}"
                ;;
            *) return 1 ;;
        esac
    done
}

echo "Approving kubelet CSR requests until bootstrapping is complete..."
# Names of the requests already reported as left pending, so that each one is
# described once rather than on every pass.
logged=" "
while [ ! -f /opt/openshift/.bootkube.done ]
do
    # Approve only the kubelet client and serving requests that match every
    # expected attribute. Anything else is left pending for
    # cluster-machine-approver rather than approved here.
    while IFS=$'\t' read -r name attributes signer requester request
    do
        [ -n "${name}" ] || continue
        [ "${request}" = "<none>" ] && request=""
        pem="$(printf '%s' "${request}" | base64 --decode 2>/dev/null)"

        case "${attributes}" in
            client)
                if is_node_client_request "${pem}"; then
                    echo "Approving node client CSR ${name}"
                    oc --kubeconfig="$KUBECONFIG" adm certificate approve "${name}"
                    continue
                fi
                ;;
            serving)
                if is_node_serving_request "${pem}" "${requester}" "${name}"; then
                    echo "Approving node serving CSR ${name}"
                    oc --kubeconfig="$KUBECONFIG" adm certificate approve "${name}"
                    continue
                fi
                ;;
            *) ;;
        esac

        case "${logged}" in
            *" ${name} "*) continue ;;
            *) logged="${logged}${name} " ;;
        esac
        echo "Leaving CSR ${name} pending: signer=${signer} requester=${requester}"
    done <<< "$(oc --kubeconfig="$KUBECONFIG" get csr -o json 2>/dev/null | jq -r \
        --arg client_signer "${CLIENT_SIGNER}" \
        --arg client_requester "${CLIENT_REQUESTER}" \
        --arg client_groups "${CLIENT_GROUPS}" \
        --arg client_usages "${CLIENT_USAGES}" \
        --arg serving_signer "${SERVING_SIGNER}" \
        --arg serving_groups "${SERVING_GROUPS}" \
        --arg serving_usages "${SERVING_USAGES}" '
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
            (if .spec.signerName == $client_signer
                and .spec.username == $client_requester
                and (((.spec.groups // []) | sort | join(",")) == $client_groups)
                and (((.spec.usages // []) | sort | join(",")) == $client_usages)
             then "client"
             elif .spec.signerName == $serving_signer
                and ((.spec.username // "") | startswith("system:node:"))
                and (((.spec.groups // []) | sort | join(",")) == $serving_groups)
                and (((.spec.usages // []) | sort | join(",")) == $serving_usages)
             then "serving"
             else "other" end),
            (.spec.signerName | or_none),
            (.spec.username | or_none),
            (.spec.request | or_none)
          ]
        | @tsv')"

    sleep 20
done
