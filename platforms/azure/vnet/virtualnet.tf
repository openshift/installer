
data "azurerm_virtual_network" "cluster_vpc" {
  id = "${var.tectonic_azure_external_vnet_id == "" ? aws_vpc.new_vpc.id : var.tectonic_azure_external_vnet_id }"
}

resource "azurerm_virtual_network" "new_vnet" {
  count               = "${length(var.tectonic_azure_external_vnet_id) > 0 ? 0 : 1}"
  name                = "${var.tectonic_cluster_name}"
  resource_group_name = "${var.resource_group_name}"
  address_space       = ["${var.tectonic_azure_vnet_cidr_block}"]
  location            = "${var.location}"

  subnet {
    name           = "cluster_default"
    address_prefix = "10.0.0.0/16"
    security_group = "${azurerm_network_security_group.cluster_default.id}"
  }

  tags {
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }
}
