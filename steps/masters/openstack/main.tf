provider "openstack" {
  cloud = "${var.tectonic_openstack_cloud}"
}

module "masters" {
  source = "../../../modules/openstack/masters"

  base_domain    = "${var.tectonic_base_domain}"
  cluster_id     = "${var.tectonic_cluster_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  flavor_name    = "${var.tectonic_openstack_master_flavor_name}"
  image_name     = "${var.tectonic_openstack_master_image_name}"
  instance_count = "${var.tectonic_bootstrap == "true" ? 1 : var.tectonic_master_count}"
  master_sg_ids  = "${concat(var.tectonic_openstack_master_extra_sg_ids, list(local.master_sg_id))}"
  key_pair       = "${var.tectonic_openstack_key_pair}"
  user_data_ign  = "${file("${path.cwd}/${var.tectonic_ignition_master}")}"
}
