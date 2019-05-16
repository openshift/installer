resource "azurerm_network_security_group" "master" {
  name                = "${var.cluster_id}-controlplane-nsg"
  location            = var.region
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet_network_security_group_association" "master" {
  subnet_id                 = azurerm_subnet.master_subnet.id
  network_security_group_id = azurerm_network_security_group.master.id
}

resource "azurerm_network_security_group" "worker" {
  name                = "${var.cluster_id}-node-nsg"
  location            = var.region
  resource_group_name = var.resource_group_name
}

resource "azurerm_subnet_network_security_group_association" "worker" {
  subnet_id                 = azurerm_subnet.node_subnet.id
  network_security_group_id = azurerm_network_security_group.worker.id
}

resource "azurerm_network_security_rule" "apiserver_in" {
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
  network_security_group_name = azurerm_network_security_group.master.name
}

resource "azurerm_network_security_rule" "sint_in" {
  name                        = "sint_in"
  priority                    = 102
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22623"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = azurerm_network_security_group.master.name
}

resource "azurerm_network_security_rule" "ssh_in" {
  name                        = "ssh_in"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = var.resource_group_name
  network_security_group_name = azurerm_network_security_group.master.name
}

