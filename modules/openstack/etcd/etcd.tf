resource "openstack_compute_servergroup_v2" "etcd_group" {
  name     = "${var.cluster_name}-etcd-group"
  policies = ["anti-affinity"]
}

resource "openstack_compute_instance_v2" "etcd_node" {
  count = "${var.instance_count}"
  name  = "${var.cluster_name}-etcd-${count.index}.${var.base_domain}"

  flavor_name     = "${var.flavor_name}"
  image_name      = "${var.image_name}"
  key_pair        = "${var.key_pair}"
  security_groups = ["${var.etcd_sg_ids}"]
  user_data       = "${data.ignition_config.tnc.*.rendered[count.index]}"

  metadata {
    role              = "etcd"
    tectonicClusterID = "${var.cluster_id}"
  }
}
