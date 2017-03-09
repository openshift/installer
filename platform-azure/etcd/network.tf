resource "azurerm_network_interface" "etcd_nic" {
  name                      = "${var.cluster_name}_etcd_nic"
  location                  = "${var.location}"
  network_security_group_id = "${azurerm_network_security_group.etcd_group.id}"
  resource_group_name       = "${var.resource_group_name}"

  ip_configuration {
    name                                    = "tectonic_etcd_configuration"
    subnet_id                               = "${azurerm_subnet.etcd_subnet.id}"
    private_ip_address_allocation           = "Dynamic"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.etcd-lb.id}"]
  }
}

resource "azurerm_subnet" "etcd_subnet" {
  name                 = "${var.cluster_name}_etcd_subnet"
  resource_group_name  = "${var.resource_group_name}"
  virtual_network_name = "${azurerm_virtual_network.etcd_vnet.name}"
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_virtual_network" "etcd_vnet" {
  name                = "${var.cluster_name}_etcd_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"
}

resource "azurerm_network_security_group" "etcd_group" {
  name                = "${var.cluster_name}_etcd_group"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  security_rule {
    name                       = "rule1"
    source_port_range          = 22
    destination_port_range     = 22
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "100"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule2"
    source_port_range          = 2379
    destination_port_range     = 2380
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "101"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule3"
    source_port_range          = "*"
    destination_port_range     = "*"
    protocol                   = "*"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "103"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "rule4"
    source_port_range          = "*"
    destination_port_range     = "*"
    protocol                   = "*"
    destination_address_prefix = "Internet"
    source_address_prefix      = "0.0.0.0/0"
    access                     = "Allow"
    priority                   = "104"
    direction                  = "Outbound"
  }
}
