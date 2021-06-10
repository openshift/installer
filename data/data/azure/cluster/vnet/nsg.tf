data "azurerm_network_security_group" "cluster" {
  name                = "${var.cluster_id}-nsg"
  resource_group_name = var.resource_group_name
}
