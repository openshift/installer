locals {
  # Common locals
  prefix    = var.cluster_id
  zones_all = distinct(concat(var.zones_master, var.zones_worker))

  # VPC locals
  vpc_id  = var.preexisting_vpc ? data.ibm_is_vpc.vpc[0].id : ibm_is_vpc.vpc[0].id
  vpc_crn = var.preexisting_vpc ? data.ibm_is_vpc.vpc[0].crn : ibm_is_vpc.vpc[0].crn

  # LB locals
  port_kubernetes_api   = 6443
  port_machine_config   = 22623
  control_plane_subnets = var.preexisting_vpc ? data.ibm_is_subnet.control_plane[*] : ibm_is_subnet.control_plane[*]
  compute_subnets       = var.preexisting_vpc ? data.ibm_is_subnet.compute[*] : ibm_is_subnet.compute[*]

  # SG locals
  subnet_cidr_blocks = concat(local.control_plane_subnets[*].ipv4_cidr_block, local.compute_subnets[*].ipv4_cidr_block)
}

data "ibm_is_vpc" "vpc" {
  count = var.preexisting_vpc ? 1 : 0

  name = var.cluster_vpc
}

data "ibm_is_subnet" "control_plane" {
  count = var.preexisting_vpc ? length(var.control_plane_subnets) : 0

  name = var.control_plane_subnets[count.index]
}

data "ibm_is_subnet" "compute" {
  count = var.preexisting_vpc ? length(var.compute_subnets) : 0

  name = var.compute_subnets[count.index]
}
