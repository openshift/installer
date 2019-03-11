# resource "azurerm_network_security_group" "worker" {
#   name                = "${var.cluster_id}-worker-nsg"
#   location            = "${var.region}"
#   resource_group_name = "${var.rg_name}"

#   security_rule {
#     name                       = "AllowSSH"
#     priority                   = 100
#     direction                  = "Inbound"
#     access                     = "Allow"
#     protocol                   = "Tcp"
#     source_port_range          = "22"
#     destination_port_range     = "22"
#     source_address_prefix      = "*"
#     destination_address_prefix = "VirtualNetwork"
#   }

#   tags = "${merge(map(
#     "Name", "${var.cluster_id}-worker-nsg",
#   ), var.tags)}"
# }

# resource "azurerm_subnet_network_security_group_association" "worker" {
#   subnet_id                 = "${azurerm_subnet.public_subnet.id}"
#   network_security_group_id = "${azurerm_network_security_group.public_subnet.id}"
# }
