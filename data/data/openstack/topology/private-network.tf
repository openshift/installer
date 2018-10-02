locals {
  new_master_cidr_range = "${cidrsubnet(var.cidr_block, 1, 0)}"
  new_worker_cidr_range = "${cidrsubnet(var.cidr_block, 1, 1)}"
}

resource "openstack_networking_network_v2" "openshift-private" {
  name           = "openshift"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "masters" {
  name       = "masters"
  cidr       = "${local.new_master_cidr_range}"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.openshift-private.id}"
}

resource "openstack_networking_subnet_v2" "workers" {
  name       = "worker"
  cidr       = "${local.new_worker_cidr_range}"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.openshift-private.id}"
}

resource "openstack_networking_port_v2" "masters" {
  name  = "master-port-${count.index}"
  count = "${var.masters_count}"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.masters.id}"
  }
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "bootstrap-port"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.masters.id}"
  }
}

data "openstack_networking_network_v2" "external_network" {
  name     = "${var.external_network}"
  external = true
}

resource "openstack_networking_router_v2" "openshift-external-router" {
  name                = "openshift-external-router"
  admin_state_up      = true
  external_network_id = "${data.openstack_networking_network_v2.external_network.id}"
}

resource "openstack_networking_router_interface_v2" "masters_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.masters.id}"
}

resource "openstack_networking_router_interface_v2" "workers_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.workers.id}"
}
