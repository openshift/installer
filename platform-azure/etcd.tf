resource "azurerm_virtual_machine" "etcd_node" {
  name            = "${var.tectonic_cluster_name}_etcd_node_${count.index}"
  security_groups = ["${azurerm_network_security_group.etcd_group.name}"]

  metadata {
    role = "etcd"
  }

  user_data    = "${ignition_config.etcd.rendered}"
  config_drive = false
}

resource "azurerm_network_security_group" "etcd_group" {
  name        = "${var.tectonic_cluster_name}_etcd_group"

  security_rule {
    source_port_range      = 22
    destination_port_range = 22
    protocol               = "tcp"
    source_address_prefix  = "0.0.0.0/0"
  }

  security_rule {
    source_port_range      = 2379
    destination_port_range = 2380
    protocol               = "tcp"
    source_address_prefix  = "0.0.0.0/0"
  }

  security_rule {
    source_port_range      = -1
    destination_port_range = -1
    protocol               = "icmp"
    source_address_prefix  = "0.0.0.0/0"
  }
}
