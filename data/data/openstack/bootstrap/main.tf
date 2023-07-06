locals {
  description = "Created By OpenShift Installer"
}

provider "openstack" {
  auth_url            = var.openstack_credentials_auth_url
  cert                = var.openstack_credentials_cert
  cloud               = var.openstack_credentials_cloud
  domain_id           = var.openstack_credentials_domain_id
  domain_name         = var.openstack_credentials_domain_name
  endpoint_type       = var.openstack_credentials_endpoint_type
  insecure            = var.openstack_credentials_insecure
  key                 = var.openstack_credentials_key
  password            = var.openstack_credentials_password
  project_domain_id   = var.openstack_credentials_project_domain_id
  project_domain_name = var.openstack_credentials_project_domain_name
  region              = var.openstack_credentials_region
  swauth              = var.openstack_credentials_swauth
  tenant_id           = var.openstack_credentials_tenant_id
  tenant_name         = var.openstack_credentials_tenant_name
  token               = var.openstack_credentials_token
  use_octavia         = var.openstack_credentials_use_octavia
  user_domain_id      = var.openstack_credentials_user_domain_id
  user_domain_name    = var.openstack_credentials_user_domain_name
  user_id             = var.openstack_credentials_user_id
  user_name           = var.openstack_credentials_user_name
}

data "openstack_images_image_v2" "base_image" {
  name = var.openstack_base_image_name
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = var.openstack_master_flavor_name
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

  dynamic "fixed_ip" {
    for_each = var.nodes_default_port.fixed_ips

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

  depends_on = [var.master_port_ids]
}

resource "openstack_blockstorage_volume_v3" "bootstrap_volume" {
  name        = "${var.cluster_id}-bootstrap"
  count       = var.openstack_master_root_volume_size == null ? 0 : 1
  description = local.description

  size        = var.openstack_master_root_volume_size
  volume_type = var.openstack_master_root_volume_types[0]
  image_id    = data.openstack_images_image_v2.base_image.id

  availability_zone = var.openstack_master_root_volume_availability_zones[0]
}

resource "openstack_compute_instance_v2" "bootstrap" {
  name              = "${var.cluster_id}-bootstrap"
  flavor_id         = data.openstack_compute_flavor_v2.bootstrap_flavor.id
  image_id          = var.openstack_master_root_volume_size == null ? data.openstack_images_image_v2.base_image.id : null
  availability_zone = var.openstack_master_availability_zones[0]

  user_data = var.openstack_bootstrap_shim_ignition

  dynamic "block_device" {
    for_each = var.openstack_master_root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.bootstrap_volume[0].id]
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
    for_each = var.openstack_additional_network_ids

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
  count       = var.openstack_external_network != "" ? 1 : 0
  description = "${var.cluster_id}-bootstrap-fip"
  pool        = var.openstack_external_network
  port_id     = openstack_networking_port_v2.bootstrap_port.id
  tags        = ["openshiftClusterID=${var.cluster_id}"]

  depends_on = [openstack_compute_instance_v2.bootstrap]
}
