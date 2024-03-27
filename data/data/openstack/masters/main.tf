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

data "openstack_compute_flavor_v2" "masters_flavor" {
  name = var.openstack_master_flavor_name
}

data "ignition_file" "hostname" {
  count = var.master_count
  mode  = "420" // 0644
  path  = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-master-${count.index}
EOF
  }
}

data "ignition_config" "master_ignition_config" {
  count = var.master_count

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition_master)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, count.index)
  ]
}

resource "openstack_blockstorage_volume_v3" "master_volume" {
  name = "${var.cluster_id}-master-${count.index}"
  description = local.description
  count = var.openstack_master_root_volume_size == null ? 0 : var.master_count

  size = var.openstack_master_root_volume_size
  volume_type = var.openstack_master_root_volume_types[count.index]
  image_id = data.openstack_images_image_v2.base_image.id

  availability_zone = var.openstack_master_root_volume_availability_zones[count.index]
}

resource "openstack_compute_servergroup_v2" "master_group" {
  name = var.openstack_master_server_group_name
  policies = [var.openstack_master_server_group_policy]
}

# The master servers are created in three separate resource definition blocks,
# rather than a single block with a "count" meta-property, because they need to
# be created sequentially rather than concurrently by Terraform.
#
# The reason why they need to be created one at a time is that OpenStack's
# Compute module is currently unable to honour the "soft-anti-affinity" policy
# when the servers are created concurrently.
#
# We chose to unroll the loop into three instances, because three is the
# minimum number of required Control plane nodes, as stated in the
# documentation[1].
#
# The expectation is that machine-api-operator will take care of creating any
# other requested nodes as soon as the deployment is effective, and that a
# similar workaround is applied for day-2 operations.
#
# [1]: https://github.com/openshift/installer/tree/master/docs/user/openstack#master-nodes
resource "openstack_compute_instance_v2" "master_conf_0" {
  count = var.master_count > 0 ? 1 : 0
  name = "${var.cluster_id}-master-0"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.openstack_master_root_volume_size == null ? data.openstack_images_image_v2.base_image.id : null
  security_groups = local.master_sg_ids
  availability_zone = var.openstack_master_availability_zones[0]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    0,
  )

  dynamic "block_device" {
    for_each = var.openstack_master_root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[0].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = local.master_port_ids[0]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic "network" {
    for_each = [for port in openstack_networking_port_v2.master_0_failuredomain : port.id]

    content {
      port = network.value
    }
  }

  dynamic "network" {
    for_each = var.openstack_additional_network_ids

    content {
      uuid = network.value
    }
  }

  tags = ["openshiftClusterID=${var.cluster_id}"]

  metadata = {
    Name = "${var.cluster_id}-master"
    openshiftClusterID = var.cluster_id
  }
}

resource "openstack_compute_instance_v2" "master_conf_1" {
  count = var.master_count > 1 ? 1 : 0
  name = "${var.cluster_id}-master-1"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.openstack_master_root_volume_size == null ? data.openstack_images_image_v2.base_image.id : null
  security_groups = local.master_sg_ids
  availability_zone = var.openstack_master_availability_zones[1]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    1,
  )

  dynamic "block_device" {
    for_each = var.openstack_master_root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[1].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = local.master_port_ids[1]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic "network" {
    for_each = [for port in openstack_networking_port_v2.master_1_failuredomain : port.id]

    content {
      port = network.value
    }
  }

  dynamic "network" {
    for_each = var.openstack_additional_network_ids

    content {
      uuid = network.value
    }
  }

  tags = ["openshiftClusterID=${var.cluster_id}"]

  metadata = {
    Name = "${var.cluster_id}-master"
    openshiftClusterID = var.cluster_id
  }

  depends_on = [openstack_compute_instance_v2.master_conf_0]
}

resource "openstack_compute_instance_v2" "master_conf_2" {
  count = var.master_count > 2 ? 1 : 0
  name = "${var.cluster_id}-master-2"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.openstack_master_root_volume_size == null ? data.openstack_images_image_v2.base_image.id : null
  security_groups = local.master_sg_ids
  availability_zone = var.openstack_master_availability_zones[2]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    2,
  )

  dynamic "block_device" {
    for_each = var.openstack_master_root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[2].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = local.master_port_ids[2]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic "network" {
    for_each = [for port in openstack_networking_port_v2.master_2_failuredomain : port.id]

    content {
      port = network.value
    }
  }

  dynamic "network" {
    for_each = var.openstack_additional_network_ids

    content {
      uuid = network.value
    }
  }

  tags = ["openshiftClusterID=${var.cluster_id}"]

  metadata = {
    Name = "${var.cluster_id}-master"
    openshiftClusterID = var.cluster_id
  }

  depends_on = [openstack_compute_instance_v2.master_conf_1]
}

# Pre-create server groups for the Compute MachineSets, with the given policy.
resource "openstack_compute_servergroup_v2" "server_groups" {
  for_each = var.openstack_worker_server_group_names
  name = each.key
  policies = [var.openstack_worker_server_group_policy]
}

resource "openstack_networking_port_v2" "master_0_failuredomain" {
  count = var.master_count > 0 ? length(var.openstack_additional_ports[0]) : 0

  name = "${var.cluster_id}-master-0-${count.index}"
  description = local.description
  network_id = var.openstack_additional_ports[0][count.index].network_id
  security_group_ids = concat(var.openstack_master_extra_sg_ids, [openstack_networking_secgroup_v2.master.id])
  tags = ["openshiftClusterID=${var.cluster_id}"]

  dynamic "fixed_ip" {
    for_each = var.openstack_additional_ports[0][count.index].fixed_ips

    content {
      subnet_id = fixed_ip.value["subnet_id"]
      ip_address = fixed_ip.value["ip_address"]
    }
  }
}

resource "openstack_networking_port_v2" "master_1_failuredomain" {
  count = var.master_count > 1 ? length(var.openstack_additional_ports[1]) : 0

  name = "${var.cluster_id}-master-1-${count.index}"
  description = local.description
  network_id = var.openstack_additional_ports[1][count.index].network_id
  security_group_ids = concat(var.openstack_master_extra_sg_ids, [openstack_networking_secgroup_v2.master.id])
  tags = ["openshiftClusterID=${var.cluster_id}"]

  dynamic "fixed_ip" {
    for_each = var.openstack_additional_ports[1][count.index].fixed_ips

    content {
      subnet_id = fixed_ip.value["subnet_id"]
      ip_address = fixed_ip.value["ip_address"]
    }
  }
}

resource "openstack_networking_port_v2" "master_2_failuredomain" {
  count = var.master_count > 2 ? length(var.openstack_additional_ports[2]) : 0

  name = "${var.cluster_id}-master-2-${count.index}"
  description = local.description
  network_id = var.openstack_additional_ports[2][count.index].network_id
  security_group_ids = concat(var.openstack_master_extra_sg_ids, [openstack_networking_secgroup_v2.master.id])
  tags = ["openshiftClusterID=${var.cluster_id}"]

  dynamic "fixed_ip" {
    for_each = var.openstack_additional_ports[2][count.index].fixed_ips

    content {
      subnet_id = fixed_ip.value["subnet_id"]
      ip_address = fixed_ip.value["ip_address"]
    }
  }
}
