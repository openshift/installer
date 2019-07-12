locals {
  nodes_cidr_block   = cidrsubnet(var.cidr_block, 1, 0)
  service_cidr_block = cidrsubnet(var.cidr_block, 1, 1)
}

data "openstack_networking_network_v2" "external_network" {
  name       = var.external_network
  network_id = var.external_network_id
  external   = true
}

resource "openstack_networking_network_v2" "openshift-private" {
  name           = "${var.cluster_id}-openshift"
  admin_state_up = "true"
  tags           = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "service" {
  name       = "${var.cluster_id}-service"
  cidr       = local.service_cidr_block
  ip_version = 4
  network_id = openstack_networking_network_v2.openshift-private.id
  tags       = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "nodes" {
  name            = "${var.cluster_id}-nodes"
  cidr            = local.nodes_cidr_block
  ip_version      = 4
  network_id      = openstack_networking_network_v2.openshift-private.id
  tags            = ["openshiftClusterID=${var.cluster_id}"]
  dns_nameservers = [openstack_networking_port_v2.service_port.all_fixed_ips[0]]

  # We reserve some space at the beginning of the CIDR to use for the VIPs
  # It would be good to make this more dynamic by calculating the number of
  # addresses in the provided CIDR. This currently assumes at least a /18.
  allocation_pool {
    start = cidrhost(local.nodes_cidr_block, 10)
    end   = cidrhost(local.nodes_cidr_block, 16000)
  }
}

resource "openstack_networking_port_v2" "masters" {
  name  = "${var.cluster_id}-master-port-${count.index}"
  count = var.masters_count

  admin_state_up     = "true"
  network_id         = openstack_networking_network_v2.openshift-private.id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.nodes.id
  }

  allowed_address_pairs {
    ip_address = var.api_int_ip
  }

  allowed_address_pairs {
    ip_address = var.node_dns_ip
  }
}

resource "openstack_networking_trunk_v2" "masters" {
  name  = "${var.cluster_id}-master-trunk-${count.index}"
  count = var.trunk_support ? var.masters_count : 0
  tags  = ["openshiftClusterID=${var.cluster_id}"]

  admin_state_up = "true"
  port_id        = openstack_networking_port_v2.masters[count.index].id
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "${var.cluster_id}-bootstrap-port"

  admin_state_up     = "true"
  network_id         = openstack_networking_network_v2.openshift-private.id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.nodes.id
  }
}

resource "openstack_networking_port_v2" "service_port" {
  name = "${var.cluster_id}-service-port"

  admin_state_up     = "true"
  network_id         = openstack_networking_network_v2.openshift-private.id
  security_group_ids = [openstack_networking_secgroup_v2.api.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.service.id
  }
}

resource "openstack_networking_floatingip_associate_v2" "service_fip" {
  count       = length(var.lb_floating_ip) == 0 ? 0 : 1
  port_id     = openstack_networking_port_v2.service_port.id
  floating_ip = var.lb_floating_ip
}

resource "openstack_networking_router_v2" "openshift-external-router" {
  name                = "${var.cluster_id}-external-router"
  admin_state_up      = true
  external_network_id = data.openstack_networking_network_v2.external_network.id
  tags                = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_router_interface_v2" "service_router_interface" {
  router_id = openstack_networking_router_v2.openshift-external-router.id
  subnet_id = openstack_networking_subnet_v2.service.id
}

resource "openstack_networking_router_interface_v2" "nodes_router_interface" {
  router_id = openstack_networking_router_v2.openshift-external-router.id
  subnet_id = openstack_networking_subnet_v2.nodes.id
}

