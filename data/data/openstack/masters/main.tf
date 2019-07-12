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

data "ignition_config" "master_ignition_config" {
  count = var.instance_count

  append {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.user_data_ign)}"
  }

  files = [
    element(data.ignition_file.hostname.*.id, count.index)
  ]
}

resource "openstack_compute_instance_v2" "master_conf" {
  name = "${var.cluster_id}-master-${count.index}"
  count = var.instance_count

  flavor_id = data.openstack_compute_flavor_v2.masters_flavor.id
  image_id = data.openstack_images_image_v2.masters_img.id
  security_groups = var.master_sg_ids
  user_data = element(
    data.ignition_config.master_ignition_config.*.rendered,
    count.index,
  )

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

