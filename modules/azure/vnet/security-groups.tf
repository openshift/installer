resource "azurerm_network_security_group" "cluster_default" {
  name                = "${var.tectonic_cluster_name}"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  # TODO: enable all inbound traffic to make debugging easier
  security_rule {
    name                       = "cluster_default_ingress"
    priority                   = 4000
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  # Allow horizontal traffic in the vnet
  security_rule {
    name                       = "cluster_default_internal"
    priority                   = 4050
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "VirtualNetwork"
    destination_address_prefix = "VirtualNetwork"
  }

  # Allow internet outbound for all machines
  security_rule {
    name                       = "cluster_default_egress"
    priority                   = 4051
    direction                  = "Outbound"
    access                     = "Allow"
    protocol                   = "*"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}
