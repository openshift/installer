#!/bin/bash

cat <<EOF > terraform.tfvars
vsphere_cluster = "${VSPHERE_CLUSTER:-devel}"

vsphere_network = "${VSPHERE_NETWORK:-VM Network}"

vm_template = "${VM_TEMPLATE:-rhcos-latest}"

ipam = "${IPAM:-139.178.89.254}"

ipam_token = "${IPAM_TOKEN}"

bootstrap_ignition_url = "${BOOTSTRAP_IGNITION_URL}"


## The rest of the variables are all taken from the installer. You should not need to change any of these.

cluster_name = $(jq '.["*installconfig.InstallConfig"].config.metadata.name' .openshift_install_state.json)
vsphere_server = $(jq '.["*installconfig.InstallConfig"].config.platform.vsphere.vCenter' .openshift_install_state.json)
vsphere_user = $(jq '.["*installconfig.InstallConfig"].config.platform.vsphere.username' .openshift_install_state.json)
vsphere_password = $(jq '.["*installconfig.InstallConfig"].config.platform.vsphere.password' .openshift_install_state.json)
vsphere_datacenter = $(jq '.["*installconfig.InstallConfig"].config.platform.vsphere.datacenter' .openshift_install_state.json)
vsphere_datastore = $(jq '.["*installconfig.InstallConfig"].config.platform.vsphere.defaultDatastore' .openshift_install_state.json)
machine_cidr = $(jq '.["*installconfig.InstallConfig"].config.networking.machineCIDR' .openshift_install_state.json)
base_domain = $(jq '.["*installconfig.InstallConfig"].config.baseDomain' .openshift_install_state.json)
control_plane_ignition = <<END_OF_CONTROL_PLANE_IGNITION
$(jq . master.ign)
END_OF_CONTROL_PLANE_IGNITION
compute_ignition = <<END_OF_COMPUTE_IGNITION
$(jq . worker.ign)
END_OF_COMPUTE_IGNITION
EOF
