## Network
## These are be optional arguments in the install-config (e.g. Platform)
## so that users may specify them. Have them be "hidden" in that the survey doesn't ask for them
## unless the OCP leads disagree.
## And since they're optional, use the count = construct to conditionally create them if the tf
## vars are unset.

## Note, the following are incomplete placeholders to be tested and reviewed later when the TF
## support for these has been added to the ibmcloud terraform provider (which is now forked into
## https://github.com/openshift/terraform-provider-ibm)

#resource "ibm_direct_link" "ocp_direct_link" {
# TODO
#}

#resource "ibm_pi_network" "ocp_network" {
#  provider             = ibm.powervs
#  count                = var.powervs_network_name == "" ? 1 : 0
#  pi_network_name      = "pvs-net-${var.cluster_id}"
#  pi_cloud_instance_id = "powervs_cloud_instance_id"
#  pi_network_type      = "dhcp"
#  pi_cidr              = "192.168.0.0/24"
#  pi_dns               = [<"DNS Servers">]
#}

#resource "ibm_is_vpc" "ocp_vpc" {
#  provider       = ibm.vpc
#  count          = var.powervs_vpc == "" ? 1 : 0
#  name           = "vpc_${var.cluster_id}"
#  classic_access = false
#  resource_group = var.powervs_resource_group
#}

#resource "ibm_is_subnet" "ocp_vpc_subnet" {
#  provider        = ibm.vpc
#  count           = var.powervs_vpc_subnet == "" ? 1 : 0
#  name            = "vpc_subnet_${var.cluster_id}"
#  vpc             = ibm_is_vpc..id
#  ipv4_cidr_block = "192.168.0.0/1"
#}
