# https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/tg_gateway
resource "ibm_tg_gateway" "transit_gateway" {
  count          = var.attached_transit_gateway == "" ? 1 : 0
  name           = "tg-${var.cluster_id}"
  location       = var.vpc_region
  global         = true
  resource_group = data.ibm_resource_group.rg_pvs_ipi_rg.id
}

# https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/data-sources/resource_group
data "ibm_resource_group" "rg_pvs_ipi_rg" {
  name = var.resource_group
}

# https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/tg_connection
resource "ibm_tg_connection" "tg_connection_vpc" {
  count        = var.tg_connection_vpc_id == "" ? 1 : 0
  gateway      = var.attached_transit_gateway == "" ? resource.ibm_tg_gateway.transit_gateway[0].id : var.attached_transit_gateway
  network_type = "vpc"
  name         = "tg-${var.cluster_id}-conn-vpc"
  network_id   = var.vpc_crn
}

# https://registry.terraform.io/providers/IBM-Cloud/ibm/latest/docs/resources/tg_connection
resource "ibm_tg_connection" "tg_connection_pvs" {
  count        = var.attached_transit_gateway == "" ? 1 : 0
  gateway      = resource.ibm_tg_gateway.transit_gateway[count.index].id
  network_type = "power_virtual_server"
  name         = "tg-${var.cluster_id}-conn-pvs"
  network_id   = var.service_instance_crn
}
