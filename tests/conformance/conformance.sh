#!/bin/bash
set -exo pipefail

PROJECT=/go/src/github.com/coreos/tectonic-installer

export PLATFORM=aws
export CLUSTER="tf-${PLATFORM}-${BUILD_ID}"
export TF_VAR_tectonic_pull_secret_path=${TF_VAR_tectonic_pull_secret_path}
export TF_VAR_tectonic_license_path=${TF_VAR_tectonic_license_path}
export TECTONIC_BUILDER=quay.io/coreos/tectonic-builder:v1.36
export KUBE_CONFORMANCE=quay.io/coreos/kube-conformance:v1.7.1_coreos.0

# Create an env var file
# shellcheck disable=SC2154
{
cat <<EOF > env.list
PLATFORM=aws
CLUSTER="tf-${PLATFORM}-${BUILD_ID}"
TF_VAR_tectonic_cluster_name=$(echo "${CLUSTER}" | awk '{print tolower($0)}')
AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
TF_VAR_tectonic_aws_region="us-west-2"
TF_VAR_tectonic_pull_secret_path=${TF_VAR_tectonic_pull_secret_path}
TF_VAR_tectonic_license_path=${TF_VAR_tectonic_license_path}
TF_VAR_tectonic_aws_ssh_key="jenkins"
TF_VAR_tectonic_admin_email=${TF_VAR_tectonic_admin_email}
TF_VAR_tectonic_admin_password_hash=${TF_VAR_tectonic_admin_password_hash}
EOF
}

# Mount secret command for docker run
MNT_SECRETS="-v ${TF_VAR_tectonic_license_path}:${TF_VAR_tectonic_license_path} -v ${TF_VAR_tectonic_pull_secret_path}:${TF_VAR_tectonic_pull_secret_path}"

function cleanup() {
    # shellcheck disable=SC2086
    docker run --env-file ./env.list -i -v ${WORKSPACE}:${PROJECT} ${MNT_SECRETS} ${TECTONIC_BUILDER} /bin/bash <<EOF
cd ${PROJECT}
make destroy
EOF
}
# Trap to ensure cluster clean up before exit
trap cleanup EXIT

export KUBECTL=${WORKSPACE}/kubectl
curl -o "${KUBECTL}" https://storage.googleapis.com/kubernetes-release/release/v1.6.3/bin/linux/amd64/kubectl && chmod +x "${KUBECTL}"

function kubectl() {
    local i=0

    echo "Executing kubectl" "$@"
    while true; do
        (( i++ )) && (( i == 100 )) && echo "kubectl failed, giving up" && exit 1

        set +e
        out=$($KUBECTL "$@" 2>&1)
        status=$?
        set -e

        if [[ "$out" == *"AlreadyExists"* ]]; then
            echo "$out, skipping"
            return
        fi

        echo "$out"
        if [ "$status" -eq 0 ]; then
            return
        fi

        echo "kubectl failed, retrying in 5 seconds"
        sleep 5
    done
}

function wait_for_pods() {
    set +e
    echo "Waiting for pods in namespace $1"
    while true; do

        out=$($KUBECTL -n "$1" get po -o custom-columns=STATUS:.status.phase,NAME:.metadata.name)
        status=$?
        echo "$out"

        if [ "$status" -ne "0" ]; then
           echo "kubectl command failed, retrying in 5 seconds"
           sleep 5
           continue
        fi

        # make sure kubectl does not return "no resources found"
        if [ "$(echo "$out" | tail -n +2 | grep -c '^')" -eq 0 ]; then
           echo "no nodes were found, retrying in 5 seconds"
           sleep 5
           continue
        fi

        stat=$( echo "$out"| tail -n +2 | grep -v '^Running')
        if [[ "$stat" == "" ]]; then
            return
        fi

        echo "Pods not available yet, waiting for 5 seconds"
        sleep 5
    done
    set -e
}

# shellcheck disable=SC2086
docker run --env-file ./env.list -i -v ${WORKSPACE}:${PROJECT} ${MNT_SECRETS} -w ${PROJECT} ${TECTONIC_BUILDER} /bin/bash <<EOF

mkdir -p ${PROJECT}/build/${CLUSTER}/
ln -sf ${PROJECT}/tests/smoke/aws/vars/aws.tfvars ${PROJECT}/build/${CLUSTER}/terraform.tfvars

make plan
make apply
EOF

export KUBECONFIG=${WORKSPACE}/build/${CLUSTER}/generated/auth/kubeconfig

wait_for_pods kube-system
echo "API Server UP!"

kubectl get pods --all-namespaces
kubectl get nodes

# shellcheck disable=SC2086
docker pull ${KUBE_CONFORMANCE}
# shellcheck disable=SC2086
docker run -v ${KUBECONFIG}:/kubeconfig ${KUBE_CONFORMANCE}
