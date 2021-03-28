resource "azurestack_network_security_group" "cluster" {
  name                = "${var.cluster_id}-nsg"
  location            = var.region
  resource_group_name = var.resource_group_name
}

// TODO: We can add this association when the subnet is created. However, if the subnet is user-supplied, then we don't
// have that flexibility.
//resource "azurestack_subnet_network_security_group_association" "master" {
//  count = var.preexisting_network ? 0 : 1
//
//  subnet_id                 = azurestack_subnet.master_subnet[0].id
//  network_security_group_id = azurestack_network_security_group.cluster.id
//}
//
//resource "azurestack_subnet_network_security_group_association" "worker" {
//  count = var.preexisting_network ? 0 : 1
//
//  subnet_id                 = azurestack_subnet.worker_subnet[0].id
//  network_security_group_id = azurestack_network_security_group.cluster.id
//}

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
