resource "azurerm_network_security_group" "cluster" {
  name                = "${var.cluster_id}-nsg"
  location            = var.azure_region
  resource_group_name = data.azurerm_resource_group.main.name
}

resource "azurerm_subnet_network_security_group_association" "master" {
  count = var.azure_preexisting_network ? 0 : 1

  subnet_id                 = azurerm_subnet.master_subnet[0].id
  network_security_group_id = azurerm_network_security_group.cluster.id
}

resource "azurerm_subnet_network_security_group_association" "worker" {
  count = var.azure_preexisting_network ? 0 : 1

  subnet_id                 = azurerm_subnet.worker_subnet[0].id
  network_security_group_id = azurerm_network_security_group.cluster.id
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
  resource_group_name         = data.azurerm_resource_group.main.name
  network_security_group_name = azurerm_network_security_group.cluster.name
  description                 = local.description
}
