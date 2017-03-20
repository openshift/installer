resource "azurerm_virtual_network" "tectonic_vnet" {
  name                = "${var.tectonic_cluster_name}"
  resource_group_name = "${var.resource_group_name}"
  address_space       = ["${var.vnet_cidr_block}"]
  location            = "${var.location}"

  tags {
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }
}

resource "azurerm_subnet" "master_subnet" {
  name                 = "${var.tectonic_cluster_name}_master_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${azurerm_virtual_network.tectonic_vnet.name}"
  address_prefix       = "${cidrsubnet(azurerm_virtual_network.tectonic_vnet.address_space[0], 4, 0)}"
}

resource "azurerm_subnet" "worker_subnet" {
  name                 = "${var.tectonic_cluster_name}_worker_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${azurerm_virtual_network.tectonic_vnet.name}"
  address_prefix       = "${cidrsubnet(azurerm_virtual_network.tectonic_vnet.address_space[0], 4, 1)}"
}
