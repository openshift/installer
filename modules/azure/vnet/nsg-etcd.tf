resource "azurerm_network_security_group" "etcd" {
  count               = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                = "${var.tectonic_cluster_name}-etcd"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_network_security_rule" "etcd_egress" {
  count                       = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-etcd_egress"
  priority                    = 2000
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_ssh" {
  count                       = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-etcd_ingress_ssh"
  priority                    = 400
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "${var.ssh_network_internal}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

# TODO: Add external SSH rule
resource "azurerm_network_security_rule" "etcd_ingress_ssh_admin" {
  count                       = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-etcd_ingress_ssh_admin"
  priority                    = 405
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "${var.ssh_network_external}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_ssh_self" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_ssh_self"
  priority               = 410
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "22"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_ssh_from_master" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_services_from_console"
  priority               = 415
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "22"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_client_self" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_client_self"
  priority               = 420
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "2379"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_client_master" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_client_master"
  priority               = 425
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "2379"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "*"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_client_worker" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_client_worker"
  priority               = 430
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "2379"

  # TODO: Need to allow traffic from workers
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}

resource "azurerm_network_security_rule" "etcd_ingress_peer" {
  count                  = "${var.external_nsg_etcd == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-etcd_ingress_peer"
  priority               = 435
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "2380"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.etcd.name}"
}
