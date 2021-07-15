resource "azurestack_network_security_group" "cluster" {
  name                = "${var.cluster_id}-nsg"
  location            = var.region
  resource_group_name = var.resource_group_name
}
resource "azurestack_network_security_rule" "apiserver_in" {
  name                        = "apiserver_in"
  priority                    = 101
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "6443"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = azurestack_network_security_group.cluster.name
  description                 = local.description
}
