locals {
  new_master_cidr_range = "${cidrsubnet(var.cidr_block, 1, 0)}"
  new_worker_cidr_range = "${cidrsubnet(var.cidr_block, 1, 1)}"
}

resource "openstack_networking_network_v2" "openshift-private" {
  name           = "openshift"
  admin_state_up = "true"
  tags           = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]
}

#resource "openstack_networking_subnet_v2" "service" {
#  name        = "service"
#  cidr        = "10.3.0.0/17"
#  ip_version  = 4
#  enable_dhcp = "true"
#  gateway_ip  = "10.3.0.254"
#  network_id  = "${openstack_networking_network_v2.openshift-private.id}"
#  tags        = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
#}

resource "openstack_networking_subnet_v2" "masters" {
  name       = "masters"
  cidr       = "${local.new_master_cidr_range}"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.openshift-private.id}"
  tags       = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "workers" {
  name       = "worker"
  cidr       = "${local.new_worker_cidr_range}"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.openshift-private.id}"
  tags       = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]
}

#resource "openstack_networking_port_v2" "api_service_port" {
#  name  = "api-service-port-${count.index}"
#  count = "${var.masters_count}"
#
#  admin_state_up     = "true"
#  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
#  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
#  tags               = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
#
#  fixed_ip {
#    "subnet_id" = "${openstack_networking_subnet_v2.service.id}"
#    "ip_address" = "10.3.0.${count.index+1}"
#  }
#}

#resource "openstack_networking_port_v2" "api_service_router_port" {
#  name  = "api-service-router-port"
#
#  admin_state_up     = "true"
#  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
#  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
#  tags               = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
#
#  fixed_ip {
#    "subnet_id" = "${openstack_networking_subnet_v2.service.id}"
#    "ip_address" = "10.3.0.253"
#  }
#}

resource "openstack_networking_port_v2" "masters" {
  name  = "master-port-${count.index}"
  count = "${var.masters_count}"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
  tags               = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.masters.id}"
  }
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "bootstrap-port"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
  tags               = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    "subnet_id" = "${openstack_networking_subnet_v2.masters.id}"
  }
}

resource "openstack_networking_port_v2" "lb_port" {
  name = "lb-port"

  admin_state_up     = "true"
  network_id         = "${openstack_networking_network_v2.openshift-private.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.master.id}"]
  tags               = ["${format("tectonicClusterID=%s", var.cluster_id)}"]

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
  tags                = ["tectonicClusterID=${var.cluster_id}", "openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_router_interface_v2" "masters_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.masters.id}"
}

resource "openstack_networking_router_interface_v2" "workers_router_interface" {
  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
  subnet_id = "${openstack_networking_subnet_v2.workers.id}"
}

#resource "openstack_networking_router_interface_v2" "service_router_interface" {
#  router_id = "${openstack_networking_router_v2.openshift-external-router.id}"
#  port_id = "${openstack_networking_port_v2.api_service_router_port.id}"
#}
