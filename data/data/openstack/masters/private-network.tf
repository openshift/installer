locals {
  # Create subnet for the first MachineNetwork CIDR if we need to
  nodes_cidr_block = var.machine_v4_cidrs[0]
  nodes_default_port = var.openstack_default_machines_port != null ? var.openstack_default_machines_port : {
    network_id = openstack_networking_network_v2.openshift-private[0].id,
    fixed_ips  = [{ subnet_id = openstack_networking_subnet_v2.nodes[0].id, ip_address = "" }],
  }
  nodes_ports   = [for port in var.openstack_machines_ports : port != null ? port : local.nodes_default_port]
  create_router = (var.openstack_external_network != "" && var.openstack_default_machines_port == null) ? 1 : 0
}

data "openstack_networking_network_v2" "external_network" {
  count      = var.openstack_external_network != "" ? 1 : 0
  name       = var.openstack_external_network
  network_id = var.openstack_external_network_id
  external   = true
}

resource "openstack_networking_network_v2" "openshift-private" {
  count          = var.openstack_default_machines_port == null ? 1 : 0
  name           = "${var.cluster_id}-openshift"
  admin_state_up = "true"
  description    = local.description
  tags           = ["openshiftClusterID=${var.cluster_id}", "${var.cluster_id}-primaryClusterNetwork"]
}

resource "openstack_networking_subnet_v2" "nodes" {
  count           = var.openstack_default_machines_port == null ? 1 : 0
  name            = "${var.cluster_id}-nodes"
  description     = local.description
  cidr            = local.nodes_cidr_block
  ip_version      = 4
  network_id      = openstack_networking_network_v2.openshift-private[0].id
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
  network_id     = local.nodes_ports[count.index].network_id
  security_group_ids = concat(
    var.openstack_master_extra_sg_ids,
    [openstack_networking_secgroup_v2.master.id],
  )
  tags = ["openshiftClusterID=${var.cluster_id}"]

  extra_dhcp_option {
    name  = "domain-search"
    value = var.cluster_domain
  }

  dynamic "fixed_ip" {
    for_each = local.nodes_ports[count.index].fixed_ips

    content {
      subnet_id  = fixed_ip.value["subnet_id"]
      ip_address = fixed_ip.value["ip_address"]
    }
  }

  dynamic "allowed_address_pairs" {
    for_each = var.openstack_user_managed_load_balancer ? [] : var.openstack_api_int_ips
    content {
      ip_address = allowed_address_pairs.value
    }
  }

  dynamic "allowed_address_pairs" {
    for_each = var.openstack_user_managed_load_balancer ? [] : var.openstack_ingress_ips
    content {
      ip_address = allowed_address_pairs.value
    }
  }

  depends_on = [openstack_networking_port_v2.api_port, openstack_networking_port_v2.ingress_port,
  data.openstack_networking_port_ids_v2.api_ports, data.openstack_networking_port_ids_v2.ingress_ports]
}

# Port needs to be created by the user when using dual-stack since SLAAC or Stateless
# does not allow specification of fixed-ips during Port creation.
data "openstack_networking_port_ids_v2" "api_ports" {
  fixed_ip   = var.openstack_api_int_ips[0]
  network_id = local.nodes_default_port.network_id
}

# Port needs to be created by the user when using dual-stack since SLAAC or Stateless
# does not allow specification of fixed-ips during Port creation.
data "openstack_networking_port_ids_v2" "ingress_ports" {
  fixed_ip   = var.openstack_ingress_ips[0]
  network_id = local.nodes_default_port.network_id
}

resource "openstack_networking_port_secgroup_associate_v2" "api_port_sg" {
  count              = (! var.openstack_user_managed_load_balancer && var.use_ipv6) ? 1 : 0
  port_id            = data.openstack_networking_port_ids_v2.api_ports.ids[0]
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  depends_on         = [data.openstack_networking_port_ids_v2.api_ports]
}

resource "openstack_networking_port_secgroup_associate_v2" "ingress_port_sg" {
  count              = (! var.openstack_user_managed_load_balancer && var.use_ipv6) ? 1 : 0
  port_id            = data.openstack_networking_port_ids_v2.ingress_ports.ids[0]
  security_group_ids = [openstack_networking_secgroup_v2.worker.id]
  depends_on         = [data.openstack_networking_port_ids_v2.ingress_ports]
}

resource "openstack_networking_port_v2" "api_port" {
  count       = var.openstack_user_managed_load_balancer || var.use_ipv6 ? 0 : 1
  name        = "${var.cluster_id}-api-port"
  description = local.description

  admin_state_up     = "true"
  network_id         = local.nodes_default_port.network_id
  security_group_ids = [openstack_networking_secgroup_v2.master.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  dynamic "fixed_ip" {
    for_each = local.nodes_default_port.fixed_ips

    content {
      subnet_id  = fixed_ip.value["subnet_id"]
      ip_address = var.openstack_api_int_ips[0]
    }
  }
}

resource "openstack_networking_port_v2" "ingress_port" {
  count       = var.openstack_user_managed_load_balancer || var.use_ipv6 ? 0 : 1
  name        = "${var.cluster_id}-ingress-port"
  description = local.description

  admin_state_up     = "true"
  network_id         = local.nodes_default_port.network_id
  security_group_ids = [openstack_networking_secgroup_v2.worker.id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  dynamic "fixed_ip" {
    for_each = local.nodes_default_port.fixed_ips

    content {
      subnet_id  = fixed_ip.value["subnet_id"]
      ip_address = var.openstack_ingress_ips[0]
    }
  }
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
  count       = (var.openstack_user_managed_load_balancer || length(var.openstack_api_floating_ip) == 0) ? 0 : 1
  port_id     = var.use_ipv6 ? data.openstack_networking_port_ids_v2.api_ports.ids[0] : openstack_networking_port_v2.api_port[0].id
  floating_ip = var.openstack_api_floating_ip
  depends_on  = [openstack_networking_router_interface_v2.nodes_router_interface]
}

resource "openstack_networking_floatingip_associate_v2" "ingress_fip" {
  count       = (var.openstack_user_managed_load_balancer || length(var.openstack_ingress_floating_ip) == 0) ? 0 : 1
  port_id     = var.use_ipv6 ? data.openstack_networking_port_ids_v2.ingress_ports.ids[0] : openstack_networking_port_v2.ingress_port[0].id
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
  subnet_id = openstack_networking_subnet_v2.nodes[0].id
}
