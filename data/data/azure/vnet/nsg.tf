resource "azurerm_network_security_group" "master" {
  name                = "${var.cluster_id}-master-nsg"
  location            = "${var.region}"
  resource_group_name = "${var.rg_name}"

  security_rule {
    name                       = "AllowSSH"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "22"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "VirtualNetwork"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-master-nsg",
  ), var.tags)}"
}

resource "azurerm_subnet_network_security_group_association" "master" {
  subnet_id                 = "${azurerm_subnet.public_subnet.id}"
  network_security_group_id = "${azurerm_network_security_group.master.id}"
}

resource "azurerm_network_security_rule" "master_mcs" {
  name                        = "master_mcs"
  priority                    = 200
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "22623"
  destination_port_range      = "22623"
  source_address_prefix       = "Internet" //?
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = "${var.rg_name}"
  
  network_security_group_name = "${azurerm_network_security_group.master.name}"  
}

resource "azurerm_network_security_rule" "master_ingress_https" {
  name                        = "master_ingress_https"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range           = "6443"
  destination_port_range      = "6443"
  source_address_prefix       = "Internet"
  destination_address_prefix  = "VirtualNetwork"
  resource_group_name         = "${var.rg_name}"
  network_security_group_name = "${azurerm_network_security_group.master.name}"
}

resource "azurerm_network_security_rule" "master_ingress_kube_scheduler" {
  name                        = "master_ingress_kube_scheduler"
  priority                    = 400
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "10251"
  destination_port_range      = "10251"
  source_address_prefix       = "*"
  destination_address_prefix  = "VirtualNetwork" //TODO : restrict to masters
  resource_group_name         = "${var.rg_name}"
  network_security_group_name = "${azurerm_network_security_group.master.name}" 
}
