#!/usr/bin/env bash

SED=""
BASE64=""

function main {
    local PLATFORM=$(uname)

    case "${PLATFORM}" in
        "Darwin")
            SED=('sed' '-i' '' '-e')
            BASE64=('base64' '-D')
            ;;
        "Linux")
            SED=('sed' '-i')
            BASE64=('base64' '-d')
            ;;
    esac

    case "$1" in
        "tfvars") tfvars "${2}" "${3}";;
        "assets") assets "${2}" "${3}";;
        *)
            usage
            exit 2
            ;;
    esac
}

function usage {
    echo -e "USAGE ${0##*/} <folder>
Commands:
tfvars <platform> <cloud-formation-path>\tCreate terraform.tvfars file from the given cloud-formation.json path.
assets <platform> <assets-dir>\Modify the given assets for the given platform.
"
}

function tfvars {
    case "$1" in
        openstack*|aws*)
            local cloud_formation="${2}"
            local tectonic_aws_az_count=$(jq '.Resources.AutoScaleController.Properties.AvailabilityZones|length' "${cloud_formation}")
            local tectonic_worker_count=$(jq -r .Resources.AutoScaleWorker.Properties.MinSize "${cloud_formation}")
            local tectonic_master_count=$(jq -r .Resources.AutoScaleController.Properties.MinSize "${cloud_formation}")

            local master_type=$(jq -r .Resources.LaunchConfigurationController.Properties.InstanceType "${cloud_formation}")
            local worker_type=$(jq -r .Resources.LaunchConfigurationWorker.Properties.InstanceType "${cloud_formation}")

            local tectonic_domain=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}")
            local tectonic_cluster_name=$(echo "${tectonic_domain}" | cut -d '.' -f 1)
            local tectonic_base_domain=$(echo "${tectonic_domain}" | cut -d '.' -f 2-)
            local tectonic_kube_version=$(tectonic_kube_version "${cloud_formation}")
            local tectonic_aws_ssh_key=$(jq -r .Resources.LaunchConfigurationController.Properties.KeyName "${cloud_formation}")
            local tectonic_dns_name=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}" | cut -d '.' -f 1)
            cat <<EOF
tectonic_aws_az_count = ${tectonic_aws_az_count}
tectonic_worker_count = ${tectonic_worker_count}
tectonic_master_count = ${tectonic_master_count}
tectonic_aws_master_ec2_type = "${master_type}"
tectonic_aws_worker_ec2_type = "${worker_type}"
tectonic_aws_etcd_ec2_type = ""
tectonic_aws_external_vpc_id = ""
tectonic_aws_vpc_cidr_block = ""
tectonic_base_domain = "${tectonic_base_domain}"
tectonic_cluster_name = "${tectonic_cluster_name}"
tectonic_kube_version = "${tectonic_kube_version}"
tectonic_assets_dir = "$(dirname ${2})"

tectonic_admin_email = ""
tectonic_admin_password_hash = ""
tectonic_ca_cert = ""
tectonic_ca_key = ""
tectonic_etcd_servers = [ "" ]
tectonic_license = ""
tectonic_pull_secret = ""
tectonic_aws_ssh_key = "${tectonic_aws_ssh_key}"
tectonic_cl_channel = "stable"
tectonic_dns_name = "${tectonic_dns_name}"
EOF
            ;;
        azure)
            local cloud_formation="${2}"
            local tectonic_domain=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}")
            local tectonic_cluster_name=$(echo "${tectonic_domain}" | cut -d '.' -f 1)
            local tectonic_base_domain=$(echo "${tectonic_domain}" | cut -d '.' -f 2-)
            local tectonic_kube_version=$(tectonic_kube_version "${cloud_formation}")
            local tectonic_worker_count=$(jq -r .Resources.AutoScaleWorker.Properties.MinSize "${cloud_formation}")
            local tectonic_master_count=$(jq -r .Resources.AutoScaleController.Properties.MinSize "${cloud_formation}")
            cat <<EOF
tectonic_worker_count = ${tectonic_worker_count}
tectonic_master_count = ${tectonic_master_count}
tectonic_base_domain = "${tectonic_base_domain}"
tectonic_cluster_name = "${tectonic_cluster_name}"
tectonic_kube_version = "${tectonic_kube_version}"
EOF
            ;;
        *)
            echo "ignoring unsupported platform $1"
            ;;
    esac
}

function assets {
    local assets="${2}"

    case "$1" in
        openstack*)
            local cloud_formation="${assets}/cloud-formation.json"
            local tectonic_domain=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}")

            for f in "${assets}/manifests/kube-apiserver.yaml" "${assets}/manifests/kube-controller-manager.yaml"; do
                "${SED[@]}" '/--cloud-provider=aws/d' $f
            done

            for f in $(find "${assets}" -type f -name "*.yaml"); do
                "${SED[@]}" "s/https:\/\/${tectonic_domain}/https:\/\/${tectonic_domain}:32000/g" $f
            done
            ;;
        azure)
            local cloud_formation="${assets}/cloud-formation.json"
            
            for f in "${assets}/manifests/kube-apiserver.yaml" "${assets}/manifests/kube-controller-manager.yaml"; do
                "${SED[@]}" '/--cloud-provider=aws/d' $f
            done
            ;;
        *)
            echo "ignoring unsupported platform $1"
            ;;
    esac
}

function tectonic_kube_version {
    local cloud_formation="${1}"
    local kube_env_encoded=$(jq -r .Resources.LaunchConfigurationController.Properties.UserData "${cloud_formation}" \
                                 | "${BASE64[@]}" \
                                 | jq -r '.storage.files[]|select(.path=="/etc/kubernetes/kubelet.env").contents.source' \
                                 | cut -d, -f2)
    local kube_env_decoded=$(printf '%b' "${kube_env_encoded//%/\\x}")
    echo $kube_env_decoded | cut -d\  -f2 | cut -d=  -f2
}

main $@
