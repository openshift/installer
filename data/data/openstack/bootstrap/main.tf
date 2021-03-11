locals {
  description = "Created By OpenShift Installer"
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name        = "${var.cluster_id}-bootstrap-port"
  description = local.description

  admin_state_up     = "true"
  network_id         = var.private_network_id
  security_group_ids = var.master_sg_ids
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

  depends_on = [var.master_port_ids]
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = var.flavor_name
}

resource "openstack_blockstorage_volume_v3" "bootstrap_volume" {
  name        = "${var.cluster_id}-bootstrap"
  count       = var.root_volume_size == null ? 0 : 1
  description = local.description

  size        = var.root_volume_size
  volume_type = var.root_volume_type
  image_id    = var.base_image_id
}

resource "openstack_compute_instance_v2" "bootstrap" {
  name              = "${var.cluster_id}-bootstrap"
  flavor_id         = data.openstack_compute_flavor_v2.bootstrap_flavor.id
  image_id          = var.root_volume_size == null ? var.base_image_id : null
  availability_zone = var.zone

  user_data = var.bootstrap_shim_ignition

  dynamic "block_device" {
    for_each = var.root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.bootstrap_volume[0].id]
    content {
      uuid                  = block_device.value
      source_type           = "volume"
      boot_index            = 0
      destination_type      = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = openstack_networking_port_v2.bootstrap_port.id
  }

  dynamic "network" {
    for_each = var.additional_network_ids

    content {
      uuid = network.value
    }
  }

  tags = ["openshiftClusterID=${var.cluster_id}"]

  metadata = {
    Name               = "${var.cluster_id}-bootstrap"
    openshiftClusterID = var.cluster_id
  }
}

resource "openstack_networking_floatingip_v2" "bootstrap_fip" {
  count       = var.external_network != "" ? 1 : 0
  description = "${var.cluster_id}-bootstrap-fip"
  pool        = var.external_network
  port_id     = openstack_networking_port_v2.bootstrap_port.id
  tags        = ["openshiftClusterID=${var.cluster_id}"]

  depends_on = [openstack_compute_instance_v2.bootstrap]
}
