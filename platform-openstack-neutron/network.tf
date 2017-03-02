resource "openstack_networking_router_v2" "router" {
  name = "${var.cluster_name}_router"
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
  cidr       = "192.168.1.0/24"
  ip_version = 4

  # TOOD make this configurable
  dns_nameservers = [ "8.8.8.8", "8.8.4.4" ]
}

resource "openstack_networking_router_interface_v2" "interface" {
  router_id = "${openstack_networking_router_v2.router.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet.id}"
}

resource "openstack_compute_floatingip_v2" "master" {
  count = "${var.master_count}"
  pool  = "public"
}

resource "openstack_compute_floatingip_v2" "worker" {
  count = "${var.worker_count}"
  pool  = "public"
}
