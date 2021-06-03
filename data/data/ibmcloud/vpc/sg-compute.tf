locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id
}

############################################
# Security group
############################################

# TODO: Adjust security group policies to be more restrictive

resource "ibm_is_security_group" "compute" {
  name           = "${local.prefix}-security-group-compute"
  resource_group = var.resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc.id
}

resource "ibm_is_security_group_rule" "compute_inbound" {
  group = ibm_is_security_group.compute.id
  direction = "inbound"
  remote = "0.0.0.0/0"
}

resource "ibm_is_security_group_rule" "compute_outbound" {
  group = ibm_is_security_group.compute.id
  direction = "outbound"
  remote = "0.0.0.0/0"
}