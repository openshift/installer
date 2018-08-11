resource "openstack_compute_servergroup_v2" "master_group" {
  name     = "${var.cluster_name}-master-group"
  policies = ["anti-affinity"]
}

resource "openstack_compute_instance_v2" "master_node" {
  count = "${var.instance_count}"
  name  = "${var.cluster_name}-master-${count.index}.${var.base_domain}"

  flavor_name     = "${var.flavor_name}"
  image_name      = "${var.image_name}"
  key_pair        = "${var.key_pair}"
  security_groups = ["${var.master_sg_ids}"]
  user_data       = "${var.user_data_ign}"

  metadata {
    role              = "master"
    tectonicClusterID = "${var.cluster_id}"
  }
}
