locals {
  nodes_cidr_block = var.machine_v4_cidrs[0]
  nodes_subnet_id  = var.openstack_machines_subnet_id != "" ? var.openstack_machines_subnet_id : openstack_networking_subnet_v2.nodes[0].id
  nodes_network_id = var.openstack_machines_network_id != "" ? var.openstack_machines_network_id : openstack_networking_network_v2.openshift-private[0].id
  create_router    = (var.openstack_external_network != "" && var.openstack_machines_subnet_id == "") ? 1 : 0
}

data "openstack_networking_network_v2" "external_network" {
  count      = var.openstack_external_network != "" ? 1 : 0
  name       = var.openstack_external_network
  network_id = var.openstack_external_network_id
  external   = true
}

resource "openstack_networking_network_v2" "openshift-private" {
  count          = var.openstack_machines_subnet_id == "" ? 1 : 0
  name           = "${var.cluster_id}-openshift"
  admin_state_up = "true"
  description    = local.description
  tags           = ["openshiftClusterID=${var.cluster_id}", "${var.cluster_id}-primaryClusterNetwork"]
}

resource "openstack_networking_subnet_v2" "nodes" {
  count           = var.openstack_machines_subnet_id == "" ? 1 : 0
  name            = "${var.cluster_id}-nodes"
  description     = local.description
  cidr            = local.nodes_cidr_block
  ip_version      = 4
  network_id      = local.nodes_network_id
  tags            = ["openshiftClusterID=${var.cluster_id}"]
  dns_nameservers = var.openstack_external_dns

  # We reserve some space at the beginning of the CIDR to use for the VIPs
  # FIXME(mandre) if we let the ports pick up VIPs automatically, we don't have
  # to do any of this.
  allocation_pool {
    start = cidrhost(local.nodes_cidr_block, 10)
    end   = cidrhost(local.nodes_cidr_block, pow(2, (32 - split("/", local.nodes_cidr_block)[1])) - 2)
  }
}

resource "openstack_networking_port_v2" "masters" {
  name        = "${var.cluster_id}-master-${count.index}"
  count       = var.master_count
  description = local.description

  admin_state_up = "true"
  network_id     = local.nodes_network_id
  security_group_ids = concat(
    var.openstack_master_extra_sg_ids,
    [openstack_networking_secgroup_v2.master.id],
  )
  tags = ["openshiftClusterID=${var.cluster_id}"]

  extra_dhcp_option {
    name  = "domain-search"
    value = var.cluster_domain
  }

  fixed_ip {
    subnet_id = local.nodes_subnet_id
  }

  allowed_address_pairs {
    ip_address = var.openstack_api_int_ip
  }

  allowed_address_pairs {
    ip_address = var.openstack_ingress_ip
  }

  depends_on = [openstack_networking_port_v2.api_port, openstack_networking_port_v2.ingress_port]
}

resource "openstack_networking_port_v2" "api_port" {
  name        = "${var.cluster_id}-api-port"
  description = local.description

  admin_state_up     = "true"
  network_id         = local.nodes_network_id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id  = local.nodes_subnet_id
    ip_address = var.openstack_api_int_ip
  }
}

resource "openstack_networking_port_v2" "ingress_port" {
  name        = "${var.cluster_id}-ingress-port"
  description = local.description

  admin_state_up     = "true"
  network_id         = local.nodes_network_id
  security_group_ids = [openstack_networking_secgroup_v2.worker.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id  = local.nodes_subnet_id
    ip_address = var.openstack_ingress_ip
  }
}

resource "openstack_networking_trunk_v2" "masters" {
  name        = "${var.cluster_id}-master-trunk-${count.index}"
  count       = var.openstack_trunk_support ? var.master_count : 0
  description = local.description
  tags        = ["openshiftClusterID=${var.cluster_id}"]

  admin_state_up = "true"
  port_id        = openstack_networking_port_v2.masters[count.index].id
}

// If external network is defined, assign the floating IP to one of the masters.
//
// Strictly speaking, this is not required to finish the installation. We
// support environments without floating IPs. However, since the installer
// is running outside of the nodes subnet (often outside of the OpenStack
// cluster itself), it needs a floating IP to monitor the progress.
//
// This IP address is not expected to be the final solution for providing HA.
// It is only here to let the installer finish without any errors. Configuring
// a load balancer and providing external connectivity is a post-installation
// step that can't always be automated (we need to support OpenStack clusters)
// that do not have or do not want to use Octavia.
//
// If an external network has not been defined then a floating IP
// will not be provided or assigned to the masters.
//
// If the floating IP is not provided, the installer will time out waiting for
// bootstrapping to complete, but the OpenShift cluster itself should come up
// as expected.

resource "openstack_networking_floatingip_associate_v2" "api_fip" {
  count       = length(var.openstack_api_floating_ip) == 0 ? 0 : 1
  port_id     = openstack_networking_port_v2.api_port.id
  floating_ip = var.openstack_api_floating_ip
  depends_on  = [openstack_networking_router_interface_v2.nodes_router_interface]
}

resource "openstack_networking_floatingip_associate_v2" "ingress_fip" {
  count       = length(var.openstack_ingress_floating_ip) == 0 ? 0 : 1
  port_id     = openstack_networking_port_v2.ingress_port.id
  floating_ip = var.openstack_ingress_floating_ip
  depends_on  = [openstack_networking_router_interface_v2.nodes_router_interface]
}

resource "openstack_networking_router_v2" "openshift-external-router" {
  count               = local.create_router
  description         = local.description
  name                = "${var.cluster_id}-external-router"
  admin_state_up      = true
  external_network_id = join("", data.openstack_networking_network_v2.external_network.*.id)
  tags                = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_router_interface_v2" "nodes_router_interface" {
  count     = local.create_router
  router_id = join("", openstack_networking_router_v2.openshift-external-router.*.id)
  subnet_id = local.nodes_subnet_id
}
