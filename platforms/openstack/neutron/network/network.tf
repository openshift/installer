resource "openstack_networking_router_v2" "router" {
  name             = "${var.cluster_name}_router"
  admin_state_up   = "true"
  external_gateway = "${var.external_gateway_id}"
}

resource "openstack_networking_network_v2" "network" {
  name           = "${var.cluster_name}_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet" {
  name       = "${var.cluster_name}_subnet"
  network_id = "${openstack_networking_network_v2.network.id}"
  cidr       = "${var.subnet_cidr}"
  ip_version = 4

  dns_nameservers = ["${var.dns_nameservers}"]
}

resource "openstack_networking_router_interface_v2" "interface" {
  router_id = "${openstack_networking_router_v2.router.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet.id}"
}

# master

resource "openstack_networking_port_v2" "master" {
  count          = "${var.master_count}"
  name           = "${var.cluster_name}_port_master_${count.index}"
  network_id     = "${openstack_networking_network_v2.network.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet.id}"
  }

  allowed_address_pairs {
    ip_address = "${var.service_cidr}"
  }

  allowed_address_pairs {
    ip_address = "${var.cluster_cidr}"
  }
}

resource "openstack_networking_floatingip_v2" "master" {
  count = "${var.master_count}"
  pool  = "${var.floatingip_pool}"
}

# worker

resource "openstack_networking_port_v2" "worker" {
  count          = "${var.worker_count}"
  name           = "${var.cluster_name}_port_worker_${count.index}"
  network_id     = "${openstack_networking_network_v2.network.id}"
  admin_state_up = "true"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet.id}"
  }

  allowed_address_pairs {
    ip_address = "${var.service_cidr}"
  }

  allowed_address_pairs {
    ip_address = "${var.cluster_cidr}"
  }
}

resource "openstack_networking_floatingip_v2" "worker" {
  count = "${var.worker_count}"
  pool  = "${var.floatingip_pool}"
}
