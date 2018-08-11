provider "openstack" {
  cloud = "${var.tectonic_openstack_cloud}"
}

module "workers" {
  source = "../../../modules/openstack/workers"

  base_domain    = "${var.tectonic_base_domain}"
  cluster_id     = "${var.tectonic_cluster_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  flavor_name    = "${var.tectonic_openstack_worker_flavor_name}"
  image_name     = "${var.tectonic_openstack_worker_image_name}"
  instance_count = "${var.tectonic_worker_count}"
  key_pair       = "${var.tectonic_openstack_key_pair}"
  worker_sg_ids  = "${concat(var.tectonic_openstack_worker_extra_sg_ids, list(local.worker_sg_id))}"
  user_data_ign  = "${file("${path.cwd}/${var.tectonic_ignition_worker}")}"
}
