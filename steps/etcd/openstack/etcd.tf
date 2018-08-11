provider "openstack" {
  cloud = "${var.tectonic_openstack_cloud}"
}

module "defaults" {
  source = "../../../modules/openstack/target-defaults"

  etcd_count = "${var.tectonic_etcd_count}"
}

data "template_file" "etcd_hostname_list" {
  count    = "${module.defaults.etcd_count}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "etcd" {
  source = "../../../modules/openstack/etcd"

  base_domain     = "${var.tectonic_base_domain}"
  cluster_id      = "${var.tectonic_cluster_id}"
  cluster_name    = "${var.tectonic_cluster_name}"
  container_image = "${var.tectonic_container_images["etcd"]}"
  etcd_sg_ids     = "${concat(var.tectonic_openstack_etcd_extra_sg_ids, list(local.etcd_sg_id))}"
  flavor_name     = "${var.tectonic_openstack_etcd_flavor_name}"
  image_name      = "${var.tectonic_openstack_etcd_image_name}"
  instance_count  = "${length(data.template_file.etcd_hostname_list.*.id)}"
  key_pair        = "${var.tectonic_openstack_key_pair}"
}
