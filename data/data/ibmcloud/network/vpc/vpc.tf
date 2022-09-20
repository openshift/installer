############################################
# VPC
############################################

resource "ibm_is_vpc" "vpc" {
  count          = var.preexisting_vpc ? 0 : 1
  name           = "${local.prefix}-vpc"
  resource_group = var.network_resource_group_id
  tags           = var.tags
}

############################################
# Public gateways
############################################

resource "ibm_is_public_gateway" "public_gateway" {
  count = var.preexisting_vpc ? 0 : length(local.zones_all)

  name           = "${local.prefix}-public-gateway-${local.zones_all[count.index]}"
  resource_group = var.network_resource_group_id
  tags           = var.tags
  vpc            = ibm_is_vpc.vpc[0].id
  zone           = local.zones_all[count.index]
}

############################################
# Subnets
############################################

resource "ibm_is_subnet" "control_plane" {
  count = var.preexisting_vpc ? 0 : length(var.zones_master)

  name                     = "${local.prefix}-subnet-control-plane-${var.zones_master[count.index]}"
  resource_group           = var.network_resource_group_id
  tags                     = var.tags
  vpc                      = ibm_is_vpc.vpc[0].id
  zone                     = var.zones_master[count.index]
  public_gateway           = ibm_is_public_gateway.public_gateway[index(ibm_is_public_gateway.public_gateway.*.zone, var.zones_master[count.index])].id
  total_ipv4_address_count = "256"
}

resource "ibm_is_subnet" "compute" {
  count = var.preexisting_vpc ? 0 : length(var.zones_worker)

  name                     = "${local.prefix}-subnet-compute-${var.zones_worker[count.index]}"
  resource_group           = var.network_resource_group_id
  tags                     = var.tags
  vpc                      = ibm_is_vpc.vpc[0].id
  zone                     = var.zones_worker[count.index]
  public_gateway           = ibm_is_public_gateway.public_gateway[index(ibm_is_public_gateway.public_gateway.*.zone, var.zones_worker[count.index])].id
  total_ipv4_address_count = "256"
}
