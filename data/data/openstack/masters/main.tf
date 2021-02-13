data "openstack_compute_flavor_v2" "masters_flavor" {
  name = var.flavor_name
}

data "ignition_file" "hostname" {
  count = var.instance_count
  mode  = "420" // 0644
  path  = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-master-${count.index}
EOF
  }
}

data "ignition_config" "master_ignition_config" {
  count = var.instance_count

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.user_data_ign)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, count.index)
  ]
}

resource "openstack_blockstorage_volume_v3" "master_volume" {
  name = "${var.cluster_id}-master-${count.index}"
  count = var.root_volume_size == null ? 0 : var.instance_count

  size = var.root_volume_size
  volume_type = var.root_volume_type
  image_id = var.base_image_id
}

resource "openstack_compute_servergroup_v2" "master_group" {
  name = var.server_group_name
  policies = ["soft-anti-affinity"]
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
  name = "${var.cluster_id}-master-0"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.root_volume_size == null ? var.base_image_id : null
  security_groups = var.master_sg_ids
  availability_zone = var.zones[0 % length(var.zones)]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    0,
  )

  dynamic block_device {
    for_each = var.root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[0].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = var.master_port_ids[0]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic network {
    for_each = var.additional_network_ids

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
  name = "${var.cluster_id}-master-1"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.root_volume_size == null ? var.base_image_id : null
  security_groups = var.master_sg_ids
  availability_zone = var.zones[1 % length(var.zones)]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    1,
  )

  dynamic block_device {
    for_each = var.root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[1].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = var.master_port_ids[1]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic network {
    for_each = var.additional_network_ids

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
  name = "${var.cluster_id}-master-2"

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.root_volume_size == null ? var.base_image_id : null
  security_groups = var.master_sg_ids
  availability_zone = var.zones[2 % length(var.zones)]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    2,
  )

  dynamic block_device {
    for_each = var.root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[2].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = var.master_port_ids[2]
  }

  scheduler_hints {
    group = openstack_compute_servergroup_v2.master_group.id
  }

  dynamic network {
    for_each = var.additional_network_ids

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
