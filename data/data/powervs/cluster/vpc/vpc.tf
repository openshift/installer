data "ibm_resource_group" "rg" {
  name = var.resource_group
}

resource "ibm_is_vpc" "new_vpc" {
  count          = var.vpc_name == "" ? 1 : 0
  name           = "vpc-${var.cluster_id}"
  classic_access = false
  resource_group = data.ibm_resource_group.rg.id
}

resource "time_sleep" "wait_for_vpc" {
  count           = var.vpc_subnet_name == "" ? 1 : 0
  depends_on      = [ibm_is_vpc.new_vpc]
  create_duration = var.wait_for_vpc
}

resource "ibm_is_public_gateway" "dns_vm_gateway" {
  name = "${var.cluster_id}-gateway"
  vpc  = ibm_is_vpc.new_vpc[0].id
  zone = var.vpc_zone
}

resource "ibm_is_subnet" "new_vpc_subnet" {
  count                    = var.vpc_subnet_name == "" ? 1 : 0
  depends_on               = [time_sleep.wait_for_vpc]
  name                     = "vpc-subnet-${var.cluster_id}"
  vpc                      = ibm_is_vpc.new_vpc[0].id
  resource_group           = data.ibm_resource_group.rg.id
  total_ipv4_address_count = 256
  zone                     = var.vpc_zone
  public_gateway           = ibm_is_public_gateway.dns_vm_gateway.id
  tags                     = [var.cluster_id]
}

data "ibm_is_vpc" "ocp_vpc" {
  name = var.vpc_name == "" ? ibm_is_vpc.new_vpc[0].name : var.vpc_name
}

data "ibm_is_subnet" "ocp_vpc_subnet" {
  name = var.vpc_subnet_name == "" ? ibm_is_subnet.new_vpc_subnet[0].name : var.vpc_subnet_name
}

