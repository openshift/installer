locals {
  nodes_cidr_block = var.cidr_block
  nodes_subnet_id  = var.machines_subnet_id != "" ? var.machines_subnet_id : openstack_networking_subnet_v2.nodes[0].id
  nodes_network_id = var.machines_network_id != "" ? var.machines_network_id : openstack_networking_network_v2.openshift-private[0].id
}

data "openstack_networking_network_v2" "external_network" {
  name       = var.external_network
  network_id = var.external_network_id
  external   = true
}

resource "openstack_networking_network_v2" "openshift-private" {
  count          = var.machines_subnet_id == "" ? 1 : 0
  name           = "${var.cluster_id}-openshift"
  admin_state_up = "true"
  tags           = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_subnet_v2" "nodes" {
  count           = var.machines_subnet_id == "" ? 1 : 0
  name            = "${var.cluster_id}-nodes"
  cidr            = local.nodes_cidr_block
  ip_version      = 4
  network_id      = local.nodes_network_id
  tags            = ["openshiftClusterID=${var.cluster_id}"]
  dns_nameservers = var.external_dns

  # We reserve some space at the beginning of the CIDR to use for the VIPs
  # It would be good to make this more dynamic by calculating the number of
  # addresses in the provided CIDR. This currently assumes at least a /18.
  # FIXME(mandre) if we let the ports pick up VIPs automatically, we don't have
  # to do any of this.
  allocation_pool {
    start = cidrhost(local.nodes_cidr_block, 10)
    end   = cidrhost(local.nodes_cidr_block, 16000)
  }
}

resource "openstack_networking_port_v2" "masters" {
  name  = "${var.cluster_id}-master-port-${count.index}"
  count = var.masters_count

  admin_state_up     = "true"
  network_id         = local.nodes_network_id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  extra_dhcp_option {
    name  = "domain-search"
    value = var.cluster_domain
  }

  fixed_ip {
    subnet_id = local.nodes_subnet_id
  }

  allowed_address_pairs {
    ip_address = var.api_int_ip
  }

  allowed_address_pairs {
    ip_address = var.node_dns_ip
  }

  allowed_address_pairs {
    ip_address = var.ingress_ip
  }

  depends_on = [openstack_networking_port_v2.api_port, openstack_networking_port_v2.ingress_port]
}

resource "openstack_networking_port_v2" "api_port" {
  name = "${var.cluster_id}-api-port"

  admin_state_up     = "true"
  network_id         = local.nodes_network_id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id  = local.nodes_subnet_id
    ip_address = var.api_int_ip
  }
}

resource "openstack_networking_port_v2" "ingress_port" {
  name = "${var.cluster_id}-ingress-port"

  admin_state_up     = "true"
  network_id         = local.nodes_network_id
  security_group_ids = [openstack_networking_secgroup_v2.worker.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  fixed_ip {
    subnet_id  = local.nodes_subnet_id
    ip_address = var.ingress_ip
  }
}

resource "openstack_networking_trunk_v2" "masters" {
  name  = "${var.cluster_id}-master-trunk-${count.index}"
  count = var.trunk_support ? var.masters_count : 0
  tags  = ["openshiftClusterID=${var.cluster_id}"]

  admin_state_up = "true"
  port_id        = openstack_networking_port_v2.masters[count.index].id
}

// Assign the floating IP to one of the masters.
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
// If the floating IP is not provided, the installer will time out waiting for
// bootstrapping to complete, but the OpenShift cluster itself should come up
// as expected.
resource "openstack_networking_floatingip_associate_v2" "api_fip" {
  count       = length(var.lb_floating_ip) == 0 ? 0 : 1
  port_id     = openstack_networking_port_v2.api_port.id
  floating_ip = var.lb_floating_ip
  depends_on  = [openstack_networking_router_interface_v2.nodes_router_interface]
}

resource "openstack_networking_router_v2" "openshift-external-router" {
  name                = "${var.cluster_id}-external-router"
  admin_state_up      = true
  external_network_id = data.openstack_networking_network_v2.external_network.id
  tags                = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_router_interface_v2" "nodes_router_interface" {
  router_id = openstack_networking_router_v2.openshift-external-router.id
  subnet_id = local.nodes_subnet_id
}
