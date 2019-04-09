#!/usr/bin/env bash
set -eo pipefail

ARTIFACTS="${1:-/tmp/artifacts}"

echo "Gathering bootstrap journals ..."
mkdir -p "${ARTIFACTS}/bootstrap/journals"
for service in bootkube openshift kubelet crio
do
    journalctl --boot --no-pager --output=short --unit="${service}" > "${ARTIFACTS}/bootstrap/journals/${service}.log"
done

echo "Gathering bootstrap containers ..."
mkdir -p "${ARTIFACTS}/bootstrap/containers"
sudo crictl ps --all --quiet | while read -r container
do
    container_name="$(sudo crictl ps -a --id "${container}" -v | grep -oP "Name: \\K(.*)")"
    sudo crictl logs "${container}" >& "${ARTIFACTS}/bootstrap/containers/${container_name}.log" || true
    sudo crictl inspect "${container}" >& "${ARTIFACTS}/bootstrap/containers/${container_name}.inspect" || true
done
mkdir -p "${ARTIFACTS}/bootstrap/pods"
sudo podman ps --all --quiet | while read -r container
do
    sudo podman logs "${container}" >& "${ARTIFACTS}/bootstrap/pods/${container}.log"
    sudo podman inspect "${container}" >& "${ARTIFACTS}/bootstrap/pods/${container}.inspect"
done

# Collect cluster data
function queue() {
    local TARGET="${ARTIFACTS}/${1}"
    shift
    # shellcheck disable=SC2155
    local LIVE="$(jobs | wc -l)"
    while [[ "${LIVE}" -ge 45 ]]; do
        sleep 1
        LIVE="$(jobs | wc -l)"
    done
    # echo "${@}"
    if [[ -n "${FILTER}" ]]; then
        # shellcheck disable=SC2024
        sudo "${@}" | "${FILTER}" >"${TARGET}" &
    else
        # shellcheck disable=SC2024
        sudo "${@}" >"${TARGET}" &
    fi
}
mkdir -p "${ARTIFACTS}/control-plane" "${ARTIFACTS}/resources"

echo "Gathering cluster resources ..."
queue resources/nodes.list oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get nodes -o jsonpath --template '{range .items[*]}{.metadata.name}{"\n"}{end}'
queue resources/masters.list oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get nodes -o jsonpath -l 'node-role.kubernetes.io/master' --template '{range .items[*]}{.metadata.name}{"\n"}{end}'
# ShellCheck doesn't realize that $ns is for the Go template, not something we're trying to expand in the shell
# shellcheck disable=2016
queue resources/containers oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get pods --all-namespaces --template '{{ range .items }}{{ $name := .metadata.name }}{{ $ns := .metadata.namespace }}{{ range .spec.containers }}-n {{ $ns }} {{ $name }} -c {{ .name }}{{ "\n" }}{{ end }}{{ range .spec.initContainers }}-n {{ $ns }} {{ $name }} -c {{ .name }}{{ "\n" }}{{ end }}{{ end }}'
queue resources/api-pods oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get pods -l openshift.io/component=api --all-namespaces --template '{{ range .items }}-n {{ .metadata.namespace }} {{ .metadata.name }}{{ "\n" }}{{ end }}'

queue resources/apiservices.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get apiservices -o json
queue resources/clusteroperators.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get clusteroperators -o json
queue resources/clusterversion.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get clusterversion -o json
queue resources/configmaps.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get configmaps --all-namespaces -o json
queue resources/csr.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get csr -o json
queue resources/endpoints.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get endpoints --all-namespaces -o json
queue resources/events.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get events --all-namespaces -o json
queue resources/kubeapiserver.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get kubeapiserver -o json
queue resources/kubecontrollermanager.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get kubecontrollermanager -o json
queue resources/machineconfigpools.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get machineconfigpools -o json
queue resources/machineconfigs.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get machineconfigs -o json
queue resources/namespaces.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get namespaces -o json
queue resources/nodes.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get nodes -o json
queue resources/openshiftapiserver.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get openshiftapiserver -o json
queue resources/pods.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get pods --all-namespaces -o json
queue resources/rolebindings.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get rolebindings --all-namespaces -o json
queue resources/roles.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get roles --all-namespaces -o json
#queue resources/secrets.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get secrets --all-namespaces -o json
queue resources/services.json oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get services --all-namespaces -o json

FILTER=gzip queue resources/openapi.json.gz oc --config=/opt/openshift/auth/kubeconfig --request-timeout=5s get --raw /openapi/v2

echo "Waiting for logs ..."
wait

echo "Gather remote logs"
export MASTERS=()
if [ "$(stat --printf="%s" "${ARTIFACTS}/resources/masters.list")" -ne "0" ]
then
    # shellcheck disable=SC2030
    mapfile -t MASTERS < "${ARTIFACTS}/resources/masters.list"
else
    # Find out master IPs from etcd discovery record
    DOMAIN=$(sudo oc --config=/opt/openshift/auth/kubeconfig whoami --show-server | grep -oP "api.\\K([a-z\\.]*)")
    # shellcheck disable=SC2031
    mapfile -t MASTERS < "$(dig -t SRV "_etcd-server-ssl._tcp.${DOMAIN}" +short | cut -f 4 -d ' ' | sed 's/.$//')"
fi

for master in "${MASTERS[@]}"
do
  echo "Collecting info from ${master}"
  scp -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null /usr/local/bin/installer-masters-gather.sh "core@${master}:" || true
  mkdir -p "${ARTIFACTS}/control-plane/${master}"
  ssh -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null "core@${master}" -C 'sudo ./installer-masters-gather.sh' </dev/null  || true
  ssh -o PreferredAuthentications=publickey -o StrictHostKeyChecking=false -o UserKnownHostsFile=/dev/null "core@${master}" -C 'sudo tar c -C /tmp/artifacts/ .' </dev/null | tar -x -C "${ARTIFACTS}/control-plane/${master}/" || true
done
tar cz -C /tmp/artifacts . > ~/log-bundle.tar.gz
echo "Log bundle written to ~/log-bundle.tar.gz"
