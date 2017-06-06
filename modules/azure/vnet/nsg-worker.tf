resource "azurerm_network_security_group" "worker" {
  count               = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                = "${var.tectonic_cluster_name}-worker"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_network_security_rule" "worker_egress" {
  count                       = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-worker_egress"
  priority                    = 2010
  direction                   = "Outbound"
  access                      = "Allow"
  protocol                    = "*"
  source_port_range           = "*"
  destination_port_range      = "*"
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_ssh" {
  count                       = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-worker_ingress_ssh"
  priority                    = 600
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "${var.ssh_network_internal}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

# TODO: Add external SSH rule
resource "azurerm_network_security_rule" "worker_ingress_ssh_admin" {
  count                       = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-worker_ingress_ssh_admin"
  priority                    = 605
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "22"
  source_address_prefix       = "${var.ssh_network_external}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_services" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_services"
  priority               = 610
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "30000-32767"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_services_from_console" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_services_from_console"
  priority               = 615
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "30000-32767"

  # TODO: Need to allow traffic from console
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_flannel" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_flannel"
  priority               = 620
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "udp"
  source_port_range      = "*"
  destination_port_range = "4789"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_flannel_from_master" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_flannel_from_master"
  priority               = 625
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "udp"
  source_port_range      = "*"
  destination_port_range = "4789"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_kubelet_insecure" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_kubelet_insecure"
  priority               = 630
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "10250"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_kubelet_insecure_from_master" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_kubelet_insecure_from_master"
  priority               = 635
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "10250"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_kubelet_secure" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_kubelet_secure"
  priority               = 640
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "10255"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_kubelet_secure_from_master" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_kubelet_secure_from_master"
  priority               = 645
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "10255"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_node_exporter" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_node_exporter"
  priority               = 650
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "9100"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_node_exporter_from_master" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_node_exporter_from_master"
  priority               = 655
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "9100"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_heapster" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_heapster"
  priority               = 660
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "4194"

  # TODO: Need to allow traffic from self
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

resource "azurerm_network_security_rule" "worker_ingress_heapster_from_master" {
  count                  = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                   = "${var.tectonic_cluster_name}-worker_ingress_heapster_from_master"
  priority               = 665
  direction              = "Inbound"
  access                 = "Allow"
  protocol               = "tcp"
  source_port_range      = "*"
  destination_port_range = "4194"

  # TODO: Need to allow traffic from master
  source_address_prefix       = "${var.vnet_cidr_block}"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

# TODO: Add rules for self-hosted etcd (etcd-operator)

# TODO: Review NSG
resource "azurerm_network_security_rule" "worker_ingress_http" {
  count                       = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-worker_ingress_http"
  priority                    = 670
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "80"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}

# TODO: Review NSG
resource "azurerm_network_security_rule" "worker_ingress_https" {
  count                       = "${var.external_nsg_worker == "" ? 1 : 0}"
  name                        = "${var.tectonic_cluster_name}-worker_ingress_https"
  priority                    = 675
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "tcp"
  source_port_range           = "*"
  destination_port_range      = "443"
  source_address_prefix       = "VirtualNetwork"
  destination_address_prefix  = "*"
  resource_group_name         = "${var.external_resource_group == "" ? var.resource_group_name : var.external_resource_group}"
  network_security_group_name = "${azurerm_network_security_group.worker.name}"
}
