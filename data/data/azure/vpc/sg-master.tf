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
    destination_address_prefix = "*"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-master-nsg",
  ), var.tags)}"
}

resource "azurerm_subnet_network_security_group_association" "master" {
  subnet_id                 = "${azurerm_subnet.public_subnet.id}"
  network_security_group_id = "${azurerm_network_security_group.public_subnet.id}"
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
  destination_address_prefix  = "${data.aws_vpc.cluster_vpc.cidr_block}"
  resource_group_name         = "${var.rg_name}"
  network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
  security_group_id = "${azurerm_network_security_group.master.id}"  
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
  destination_address_prefix  = "${data.aws_vpc.cluster_vpc.cidr_block}"
  resource_group_name         = "${var.rg_name}"
  network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
  security_group_id = "${azurerm_network_security_group.master.id}" 
}

// Not Needed. same vnet?
# resource "azurerm_network_security_rule" "master_ingress_vxlan" {
#   name                        = "master_ingress_vxlan"
#   priority                    = 300
#   direction                   = "Inbound"
#   access                      = "Allow"
#   protocol                    = "Udp"
#   source_port_range           = "4789"
#   destination_port_range      = "4789"
#   source_address_prefix       = "VirtualNetwork"
#   destination_address_prefix  = "${data.aws_vpc.cluster_vpc.cidr_block}"
#   resource_group_name         = "${var.rg_name}"
#   network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
#   security_group_id = "${azurerm_network_security_group.master.id}"  
# }

# //NOT NEEDED? (same vnet)
# resource "azurerm_network_security_rule" "master_ingress_internal" {
#   name                        = "master_ingress_internal"
#   priority                    = 300
#   direction                   = "Inbound"
#   access                      = "Allow"
#   protocol                    = "*"
#   source_port_range           = "9000"
#   destination_port_range      = "9999"
#   source_address_prefix       = "VirtualNetwork"
#   destination_address_prefix  = "VirtualNetwork"
#   resource_group_name         = "${var.rg_name}"
#   network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
#   security_group_id = "${azurerm_network_security_group.master.id}"  
# }

resource "azurerm_network_security_rule" "master_ingress_kube_scheduler" {
  name                        = "master_ingress_kube_scheduler"
  priority                    = 300
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "10251"
  destination_port_range      = "10251"
  source_address_prefix       = "*"
  destination_address_prefix  = "VirtualNetwork" //TODO : restrict to masters
  resource_group_name         = "${var.rg_name}"
  network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
  security_group_id = "${azurerm_network_security_group.master.id}"  
}
//NOT NEEDED (same vnet)
# resource "azurerm_network_security_rule" "master_ingress_kube_controller_manager" {
#   name                        = "master_ingress_kube_controller_manager"
#   priority                    = 300
#   direction                   = "Inbound"
#   access                      = "Allow"
#   protocol                    = "Tcp"
#   source_port_range           = "10252"
#   destination_port_range      = "10252"
#   source_address_prefix       = "VirtualNetwork"
#   destination_address_prefix  = "VirtualNetwork"
#   resource_group_name         = "${var.rg_name}"
#   network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
#   security_group_id = "${azurerm_network_security_group.master.id}"  
# }

// NOT NEEDED (within same vnet)
# resource "azurerm_network_security_rule" "master_ingress_kube_controller_manager" {
#   name                        = "master_ingress_kube_controller_manager"
#   priority                    = 300
#   direction                   = "Inbound"
#   access                      = "Allow"
#   protocol                    = "Tcp"
#   source_port_range           = "10250"
#   destination_port_range      = "10250"
#   source_address_prefix       = "VirtualNetwork"
#   destination_address_prefix  = "VirtualNetwork"
#   resource_group_name         = "${var.rg_name}"
#   network_security_group_name = "${azurerm_network_security_group.public_subnet.name}"
#   security_group_id = "${azurerm_network_security_group.master.id}"
# }

// plus more for etcd that are not needed in Azure