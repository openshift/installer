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

resource "openstack_compute_instance_v2" "master_conf" {
  name = "${var.cluster_id}-master-${count.index}"
  count = var.instance_count

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.root_volume_size == null ? var.base_image_id : null
  security_groups = var.master_sg_ids
  availability_zone = var.zones[count.index % length(var.zones)]
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    count.index,
  )

  dynamic block_device {
    for_each = var.root_volume_size == null ? [] : [openstack_blockstorage_volume_v3.master_volume[count.index].id]
    content {
      uuid = block_device.value
      source_type = "volume"
      boot_index = 0
      destination_type = "volume"
      delete_on_termination = true
    }
  }

  network {
    port = var.master_port_ids[count.index]
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
