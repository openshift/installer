locals {
  new_private_cidr_range = "${cidrsubnet(var.cidr_block,1,1)}"
  new_public_cidr_range  = "${cidrsubnet(var.cidr_block,1,0)}"
}

resource "azurerm_virtual_network" "new_vnet" {
  resource_group_name = "${var.rg_name}"
  location            = "${var.region}"
  address_space       = "${var.cidr_block}"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-vnet",
  ), var.tags)}"
}

resource "azurerm_route_table" "route_table" {
  name                = "${var.cluster_id}-route-table"
  location            = "${var.region}"
  resource_group_name = "${var.rg_name}"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-route-table",
  ), var.tags)}"
}