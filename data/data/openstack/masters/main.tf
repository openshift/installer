data "openstack_images_image_v2" "masters_img" {
  name        = "${var.base_image}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "masters_flavor" {
  name = "${var.flavor_name}"
}

resource "openstack_compute_instance_v2" "master_conf" {
  name  = "${var.cluster_name}-master-${count.index}"
  count = "${var.instance_count}"

  flavor_id       = "${data.openstack_compute_flavor_v2.masters_flavor.id}"
  image_id        = "${data.openstack_images_image_v2.masters_img.id}"
  security_groups = ["${var.master_sg_ids}"]
  user_data       = "${var.user_data_igns[count.index]}"

  network = {
    port = "${var.subnet_ids[count.index]}"
  }

  metadata {
    Name              = "${var.cluster_name}-master"
    owned             = "kubernetes.io/cluster/${var.cluster_name}"
    tectonicClusterID = "${var.cluster_id}"
  }
}
