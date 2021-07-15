#!/usr/bin/env bash

API_SERVER_ARTIFACTS_DIR="/tmp/artifacts-api-server-temp"
function queue() {
    local TARGET="${ARTIFACTS_TEMP}/${1}"
    shift
    # shellcheck disable=SC2155
    local LIVE="$(jobs | wc -l)"
    while [[ "${LIVE}" -ge 45 ]]; do
        sleep 1
        LIVE="$(jobs | wc -l)"
    done
    if [[ -n "${FILTER-}" ]]; then
        # shellcheck disable=SC2024
        sudo KUBECONFIG="${GATHER_KUBECONFIG}" "${@}" | "${FILTER}" >"${TARGET}" &
    else
        # shellcheck disable=SC2024
        sudo KUBECONFIG="${GATHER_KUBECONFIG}" "${@}" >"${TARGET}" &
    fi
}

function cluster_bootstrap_gather() {
    GATHER_KUBECONFIG="/opt/openshift/auth/kubeconfig"

    ALTERNATIVE_KUBECONFIG="/etc/kubernetes/bootstrap-secrets/kubeconfig"
    if [[ -f ${ALTERNATIVE_KUBECONFIG} ]]; then
        GATHER_KUBECONFIG=${ALTERNATIVE_KUBECONFIG}
    fi

    echo "Using ${GATHER_KUBECONFIG} as KUBECONFIG"

    ARTIFACTS_TEMP="$(mktemp -d)"

    mkdir -p "${ARTIFACTS_TEMP}/resources"

    echo "Gathering cluster resources ..."
    queue resources/nodes.list oc --request-timeout=5s get nodes -o jsonpath --template '{range .items[*]}{.metadata.name}{"\n"}{end}'
    queue resources/masters.list oc --request-timeout=5s get nodes -o jsonpath -l 'node-role.kubernetes.io/master' --template '{range .items[*]}{.metadata.name}{"\n"}{end}'
    # ShellCheck doesn't realize that $ns is for the Go template, not something we're trying to expand in the shell
    # shellcheck disable=2016
    queue resources/containers oc --request-timeout=5s get pods --all-namespaces --template '{{ range .items }}{{ $name := .metadata.name }}{{ $ns := .metadata.namespace }}{{ range .spec.containers }}-n {{ $ns }} {{ $name }} -c {{ .name }}{{ "\n" }}{{ end }}{{ range .spec.initContainers }}-n {{ $ns }} {{ $name }} -c {{ .name }}{{ "\n" }}{{ end }}{{ end }}'
    queue resources/api-pods oc --request-timeout=5s get pods -l apiserver=true --all-namespaces --template '{{ range .items }}-n {{ .metadata.namespace }} {{ .metadata.name }}{{ "\n" }}{{ end }}'

    queue resources/apiservices.json oc --request-timeout=5s get apiservices -o json
    queue resources/clusteroperators.json oc --request-timeout=5s get clusteroperators -o json
    queue resources/clusterversion.json oc --request-timeout=5s get clusterversion -o json
    queue resources/configmaps.json oc --request-timeout=5s get configmaps --all-namespaces -o json
    queue resources/csr.json oc --request-timeout=5s get csr -o json
    queue resources/endpoints.json oc --request-timeout=5s get endpoints --all-namespaces -o json
    queue resources/events.json oc --request-timeout=5s get events --all-namespaces -o json
    queue resources/kubeapiserver.json oc --request-timeout=5s get kubeapiserver -o json
    queue resources/kubecontrollermanager.json oc --request-timeout=5s get kubecontrollermanager -o json
    queue resources/machineconfigpools.json oc --request-timeout=5s get machineconfigpools -o json
    queue resources/machineconfigs.json oc --request-timeout=5s get machineconfigs -o json
    queue resources/namespaces.json oc --request-timeout=5s get namespaces -o json
    queue resources/nodes.json oc --request-timeout=5s get nodes -o json
    queue resources/openshiftapiserver.json oc --request-timeout=5s get openshiftapiserver -o json
    queue resources/pods.json oc --request-timeout=5s get pods --all-namespaces -o json
    queue resources/rolebindings.json oc --request-timeout=5s get rolebindings --all-namespaces -o json
    queue resources/roles.json oc --request-timeout=5s get roles --all-namespaces -o json
    # this just lists names and number of keys
    queue resources/secrets-names.txt oc --request-timeout=5s get secrets --all-namespaces
    # this adds annotations, but strips out the SA tokens and dockercfg secrets which are noisy and may contain secrets in the annotations
    queue resources/secrets-names-with-annotations.txt oc --request-timeout=5s get secrets --all-namespaces -o=custom-columns=NAMESPACE:.metadata.namespace,NAME:.metadata.name,TYPE:.type,ANNOTATIONS:.metadata.annotations | grep -v -- '-token-' | grep -v -- '-dockercfg-'
    queue resources/services.json oc --request-timeout=5s get services --all-namespaces -o json

    FILTER=gzip queue resources/openapi.json.gz oc --request-timeout=5s get --raw /openapi/v2

    echo "Waiting for logs ..."
    while wait -n; do jobs; done

    if (( $(stat -c%s "${ARTIFACTS_TEMP}/resources/openapi.json.gz") <= 20 ))
    then
        rm -rf "${ARTIFACTS_TEMP}"
    else
        rm -rf "${API_SERVER_ARTIFACTS_DIR}"
        mkdir -p "${API_SERVER_ARTIFACTS_DIR}"
        mv "${ARTIFACTS_TEMP}/resources" "${API_SERVER_ARTIFACTS_DIR}"
    fi
}
