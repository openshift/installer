resource "azurerm_virtual_network" "tectonic_vnet" {
  count               = "${var.external_vnet_name == "" ? 1 :0 }"
  name                = "${var.tectonic_cluster_name}"
  resource_group_name = "${var.resource_group_name}"
  address_space       = ["${var.vnet_cidr_block}"]
  location            = "${var.location}"

  tags {
    "kubernetes.io/cluster/${var.tectonic_cluster_name}" = "owned"
  }
}

resource "azurerm_subnet" "master_subnet" {
  count                = "${var.external_vnet_name == "" ? 1 : 0}"
  name                 = "${var.tectonic_cluster_name}_master_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${var.external_vnet_name == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : var.external_vnet_name }"
  address_prefix       = "${cidrsubnet(var.vnet_cidr_block, 4, 0)}"
}

resource "azurerm_subnet" "worker_subnet" {
  count                = "${var.external_vnet_name == "" ? 1 : 0}"
  name                 = "${var.tectonic_cluster_name}_worker_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${var.external_vnet_name == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : var.external_vnet_name }"
  address_prefix       = "${cidrsubnet(var.vnet_cidr_block, 4, 1)}"
}
