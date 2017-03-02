#!/bin/sh

function main {
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
        "openstack"|"aws")
            local cloud_formation="${2}"
            local az_count=$(jq '.Resources.AutoScaleController.Properties.AvailabilityZones|length' "${cloud_formation}")
            local worker_count=$(jq -r .Resources.AutoScaleWorker.Properties.MinSize "${cloud_formation}")
            local master_count=$(jq -r .Resources.AutoScaleController.Properties.MinSize "${cloud_formation}")

            local master_type=$(jq -r .Resources.LaunchConfigurationController.Properties.InstanceType "${cloud_formation}")
            local worker_type=$(jq -r .Resources.LaunchConfigurationWorker.Properties.InstanceType "${cloud_formation}")

            local tectonic_domain=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}")
            local cluster_name=$(echo "${tectonic_domain}" | cut -d '.' -f 1)
            local base_domain=$(echo "${tectonic_domain}" | cut -d '.' -f 2-)

            cat <<EOF
az_count = ${az_count}

worker_count = ${worker_count}

master_count = ${master_count}

master_ec2_type = "${master_type}"

worker_ec2_type = "${worker_type}"

base_domain = "${base_domain}"

cluster_name = "${cluster_name}"

EOF
            ;;
        *)
            echo "unsupported platform $1"
            ;;
    esac
}

function assets {
    local assets="${2}"

    case "$1" in
        "openstack")
            local cloud_formation="${assets}/cloud-formation.json"
            local tectonic_domain=$(jq -r .Resources.TectonicDomain.Properties.Name "${cloud_formation}")

            for f in "${assets}/manifests/kube-apiserver.yaml" "${assets}/manifests/kube-controller-manager.yaml"; do
                sed -i '/--cloud-provider=aws/d' $f
            done

            for f in $(find "${assets}" -type f -name "*.yaml"); do
                sed -i "s/https:\/\/${tectonic_domain}/https:\/\/${tectonic_domain}:32000/g" $f
            done
            ;;
        *)
            echo "unsupported platform $1"
            ;;
    esac
}

main $@
