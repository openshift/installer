data "ibm_resource_group" "rg" {
  name = var.resource_group
}

resource "ibm_is_vpc" "ocp_vpc" {
  name           = "vpc-${var.cluster_id}"
  classic_access = false
  resource_group = data.ibm_resource_group.rg.id
}

resource "ibm_is_subnet" "ocp_vpc_subnet" {
  name                     = "vpc-subnet-${var.cluster_id}"
  vpc                      = ibm_is_vpc.ocp_vpc.id
  resource_group           = data.ibm_resource_group.rg.id
  total_ipv4_address_count = 256
  zone                     = var.vpc_zone
  tags                     = [var.cluster_id]
}
