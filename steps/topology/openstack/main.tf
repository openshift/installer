provider "openstack" {
  cloud = "${var.tectonic_openstack_cloud}"
}

module "sg" {
  source = "../../../modules/openstack/sg"

  cluster_id   = "${var.tectonic_cluster_id}"
  cluster_name = "${var.tectonic_cluster_name}"
}
