data "openstack_images_image_v2" "masters_img" {
  name        = "${var.base_image}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "masters_flavor" {
  name = "${var.flavor_name}"
}

data "ignition_config" "master_ignition_config" {
  append {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.user_data_ign)}"
  }

  files = [
    "${data.ignition_file.master_ifcfg.id}",
  ]
}

data "ignition_file" "master_ifcfg" {
  filesystem = "root"
  mode       = "420"                                       // 0644
  path       = "/etc/sysconfig/network-scripts/ifcfg-eth0"

  content {
    content = <<EOF
DEVICE="eth0"
BOOTPROTO="dhcp"
ONBOOT="yes"
TYPE="Ethernet"
PERSISTENT_DHCLIENT="yes"
DNS1="${var.service_vm_fixed_ip}"
PEERDNS="no"
NM_CONTROLLED="yes"
EOF
  }
}

resource "openstack_compute_instance_v2" "master_conf" {
  name  = "${var.cluster_name}-master-${count.index}"
  count = "${var.instance_count}"

  flavor_id       = "${data.openstack_compute_flavor_v2.masters_flavor.id}"
  image_id        = "${data.openstack_images_image_v2.masters_img.id}"
  security_groups = ["${var.master_sg_ids}"]
  user_data       = "${data.ignition_config.master_ignition_config.rendered}"

  network = {
    port = "${var.master_port_ids[count.index]}"
  }

  metadata {
    Name               = "${var.cluster_name}-master"
    owned              = "kubernetes.io/cluster/${var.cluster_name}"
    openshiftClusterID = "${var.cluster_id}"
  }
}
