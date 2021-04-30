locals {
  # NOTE: Defined in ./vpc.tf
  # prefix = var.cluster_id
}

############################################
# Security group
############################################

# TODO: Adjust security group policies to be more restrictive

resource "ibm_is_security_group" "control_plane" {
  name           = "${local.prefix}-security-group-control-plane"
  resource_group = var.resource_group_id
  vpc            = ibm_is_vpc.vpc.id
}

resource "ibm_is_security_group_rule" "control_plane_inbound" {
  group = ibm_is_security_group.control_plane.id
  direction = "inbound"
  remote = "0.0.0.0/0"
}

resource "ibm_is_security_group_rule" "control_plane_outbound" {
  group = ibm_is_security_group.control_plane.id
  direction = "outbound"
  remote = "0.0.0.0/0"
}