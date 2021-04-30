locals {
  prefix     = var.cluster_id
  zone_count = length(var.zone_list)
}

############################################
# VPC
############################################

resource "ibm_is_vpc" "vpc" {
  name           = "${local.prefix}-vpc"
  resource_group = var.resource_group_id
}

############################################
# Public gateways
############################################

resource "ibm_is_public_gateway" "public_gateway" {
  count          = local.zone_count

  name           = "${local.prefix}-public-gateway-${var.zone_list[count.index]}"
  resource_group = var.resource_group_id
  vpc            = ibm_is_vpc.vpc.id
  zone           = var.zone_list[count.index]
}

############################################
# Subnets
############################################

resource "ibm_is_subnet" "control_plane" {
  count                    = local.zone_count

  name                     = "${local.prefix}-subnet-control-plane-${var.zone_list[count.index]}"
  resource_group           = var.resource_group_id
  vpc                      = ibm_is_vpc.vpc.id
  zone                     = var.zone_list[count.index]
  public_gateway           = ibm_is_public_gateway.public_gateway[count.index].id
  total_ipv4_address_count = "256"
}

resource "ibm_is_subnet" "compute" {
  count                    = local.zone_count

  name                     = "${local.prefix}-subnet-compute-${var.zone_list[count.index]}"
  resource_group           = var.resource_group_id
  vpc                      = ibm_is_vpc.vpc.id
  zone                     = var.zone_list[count.index]
  public_gateway           = ibm_is_public_gateway.public_gateway[count.index].id
  total_ipv4_address_count = "256"
}