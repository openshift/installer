locals {
  nodes_cidr_block   = "${cidrsubnet(var.cidr_block, 1, 0)}"
  service_cidr_block = "${cidrsubnet(var.cidr_block, 1, 1)}"
}

resource "openstack_networking_network_v2" "openshift-private" {
  name           = "openshift"
  admin_state_up = "true"
  tags           = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "service" {
  name       = "service"
  cidr       = "${local.service_cidr_block}"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.openshift-private.id}"
  tags       = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "nodes" {
  name            = "nodes"
  cidr            = "${local.nodes_cidr_block}"
  ip_version      = 4
  network_id      = "${openstack_networking_network_v2.openshift-private.id}"
  tags            = ["openshiftClusterID=${var.cluster_id}"]
  dns_nameservers = ["${openstack_networking_port_v2.service_port.all_fixed_ips[0]}"]
}

resource "openstack_networking_port_v2" "masters" {
  name  = "master-port-${count.index}"
  count = "${var.masters_count}"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.nodes.id}"
  }
}

resource "openstack_networking_trunk_v2" "masters" {
  name  = "master-trunk-${count.index}"
  count = "${var.trunk_support ? var.masters_count : 0}"
  tags  = ["openshiftClusterID=${var.cluster_id}"]

  admin_state_up = "true"
  port_id        = "${openstack_networking_port_v2.masters.*.id[count.index]}"
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "bootstrap-port"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.nodes.id}"
  }
}

resource "openstack_networking_port_v2" "service_port" {
  name = "service-port"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.api.id}"]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.service.id}"
  }
}

data "openstack_networking_network_v2" "external_network" {
  name     = "${var.external_network}"
  external = true
}

resource "openstack_networking_floatingip_v2" "service_fip" {
  pool    = "${var.external_network}"
  port_id = "${openstack_networking_port_v2.service_port.id}"
}

resource "openstack_networking_router_v2" "openshift-external-router" {
  name                = "openshift-external-router"
  admin_state_up      = true
  external_network_id = "${data.openstack_networking_network_v2.external_network.id}"
  tags                = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_router_interface_v2" "service_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.service.id}"
}

resource "openstack_networking_router_interface_v2" "nodes_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.nodes.id}"
}
