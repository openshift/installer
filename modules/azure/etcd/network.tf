resource "azurerm_network_interface" "etcd_nic" {
  count                     = "${var.etcd_count}"
  name                      = "${var.cluster_name}-etcd-nic-${count.index}"
  location                  = "${var.location}"
  network_security_group_id = "${azurerm_network_security_group.etcd_group.id}"
  resource_group_name       = "${var.resource_group_name}"

  ip_configuration {
    name                                    = "tectonic_etcd_configuration"
    subnet_id                               = "${var.subnet}"
    private_ip_address_allocation           = "dynamic"
    load_balancer_backend_address_pools_ids = ["${azurerm_lb_backend_address_pool.etcd-lb.id}"]
  }
}

resource "azurerm_network_security_group" "etcd_group" {
  name                = "${var.cluster_name}-etcd"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  security_rule {
    name                       = "ssh"
    source_port_range          = "*"
    destination_port_range     = 22
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "VirtualNetwork"
    access                     = "Allow"
    priority                   = "100"
    direction                  = "Inbound"
  }

  security_rule {
    name                       = "etcd-client-perr"
    source_port_range          = "*"
    destination_port_range     = "2379-2380"
    protocol                   = "Tcp"
    destination_address_prefix = "0.0.0.0/0"
    source_address_prefix      = "VirtualNetwork"
    access                     = "Allow"
    priority                   = "101"
    direction                  = "Inbound"
  }

  # security_rule {
  #   name                       = "all-in"
  #   source_port_range          = "*"
  #   destination_port_range     = "*"
  #   protocol                   = "*"
  #   destination_address_prefix = "0.0.0.0/0"
  #   source_address_prefix      = "0.0.0.0/0"
  #   access                     = "Allow"
  #   priority                   = "103"
  #   direction                  = "Inbound"
  # }

  security_rule {
    name                       = "all-out"
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
