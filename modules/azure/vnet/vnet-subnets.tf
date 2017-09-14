resource "azurerm_virtual_network" "tectonic_vnet" {
  count               = "${var.external_vnet_id == "" ? 1 :0 }"
  name                = "${var.cluster_name}"
  resource_group_name = "${var.resource_group_name}"
  address_space       = ["${var.vnet_cidr_block}"]
  location            = "${var.location}"

  tags = "${merge(map(
    "Name", "${var.cluster_name}_vnet",
    "tectonicClusterID", "${var.cluster_id}"),
    var.extra_tags)}"
}

resource "azurerm_subnet" "master_subnet" {
  count                = "${var.external_master_subnet_id == "" ? 1 : 0}"
  name                 = "${var.cluster_name}_master_subnet"
  resource_group_name  = "${var.external_vnet_id == "" ? var.resource_group_name : replace(var.external_vnet_id, "${var.const_id_to_group_name_regex}", "$1")}"
  virtual_network_name = "${var.external_vnet_id == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : replace(var.external_vnet_id, "${var.const_id_to_group_name_regex}", "$2")}"
  address_prefix       = "${cidrsubnet(var.vnet_cidr_block, 4, 0)}"
}

resource "azurerm_subnet" "worker_subnet" {
  count                = "${var.external_worker_subnet_id == "" ? 1 : 0}"
  name                 = "${var.cluster_name}_worker_subnet"
  resource_group_name  = "${var.external_vnet_id == "" ? var.resource_group_name : replace(var.external_vnet_id, "${var.const_id_to_group_name_regex}", "$1")}"
  virtual_network_name = "${var.external_vnet_id == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : replace(var.external_vnet_id, "${var.const_id_to_group_name_regex}", "$2") }"
  address_prefix       = "${cidrsubnet(var.vnet_cidr_block, 4, 1)}"
}
