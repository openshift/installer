resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "${var.cluster_id}-bootstrap-port"

  admin_state_up     = "true"
  network_id         = var.private_network_id
  security_group_ids = [var.master_sg_id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  extra_dhcp_option {
    name  = "domain-search"
    value = var.cluster_domain
  }

  fixed_ip {
    subnet_id = var.nodes_subnet_id
  }

  allowed_address_pairs {
    ip_address = var.api_int_ip
  }

  allowed_address_pairs {
    ip_address = var.node_dns_ip
  }

  depends_on = [var.master_port_ids]
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = var.flavor_name
}

resource "openstack_compute_instance_v2" "bootstrap" {
  name      = "${var.cluster_id}-bootstrap"
  flavor_id = data.openstack_compute_flavor_v2.bootstrap_flavor.id
  image_id  = var.base_image_id

  user_data = var.bootstrap_shim_ignition

  network {
    port = openstack_networking_port_v2.bootstrap_port.id
  }

  metadata = {
    Name = "${var.cluster_id}-bootstrap"
    # "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    openshiftClusterID = var.cluster_id
  }
}

resource "openstack_networking_floatingip_v2" "bootstrap_fip" {
  description = "${var.cluster_id}-bootstrap-fip"
  pool        = var.external_network
  port_id     = openstack_networking_port_v2.bootstrap_port.id
  tags        = ["openshiftClusterID=${var.cluster_id}"]

  depends_on = [openstack_compute_instance_v2.bootstrap]
}
