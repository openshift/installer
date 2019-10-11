data "openstack_images_image_v2" "masters_img" {
  name        = var.base_image
  most_recent = true
}

data "openstack_compute_flavor_v2" "masters_flavor" {
  name = var.flavor_name
}

data "ignition_file" "hostname" {
  count      = var.instance_count
  filesystem = "root"
  mode       = "420" // 0644
  path       = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-master-${count.index}
EOF
  }
}

data "ignition_systemd_unit" "etcd_ramdisk" {
  name = "var-lib-etcd.mount"

  content = <<EOF
[Unit]
Description=Mount etcd as a ramdisk
Before=local-fs.target

[Mount]
What=none
Where=/var/lib/etcd
Type=tmpfs
Options=size=2G

[Install]
WantedBy=local-fs.target
EOF
}

data "ignition_config" "master_ignition_config" {
  count = var.instance_count

  append {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.user_data_ign)}"
  }

  files = [
    element(data.ignition_file.hostname.*.id, count.index)
  ]

  systemd = [
    data.ignition_systemd_unit.etcd_ramdisk.id,
  ]
}

resource "openstack_blockstorage_volume_v3" "master_volume" {
  name = "${var.cluster_id}-master-${count.index}"
  count = var.root_volume_size == null ? 0 : var.instance_count

  size = var.root_volume_size
  volume_type = var.root_volume_type
  image_id = data.openstack_images_image_v2.masters_img.id
}

resource "openstack_compute_instance_v2" "master_conf" {
  name = "${var.cluster_id}-master-${count.index}"
  count = var.instance_count

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = var.root_volume_size == null ? data.openstack_images_image_v2.masters_img.id : null
  security_groups = var.master_sg_ids
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

  metadata = {
    # FIXME(mandre) shouldn't it be "${var.cluster_id}-master-${count.index}" ?
    Name = "${var.cluster_id}-master"
    # "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    openshiftClusterID = var.cluster_id
  }
}
