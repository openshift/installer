provider "openstack" {
  auth_url            = "${var.tectonic_openstack_credentials_auth_url}"
  cert                = "${var.tectonic_openstack_credentials_cert}"
  cloud               = "${var.tectonic_openstack_credentials_cloud}"
  domain_id           = "${var.tectonic_openstack_credentials_domain_id}"
  domain_name         = "${var.tectonic_openstack_credentials_domain_name}"
  endpoint_type       = "${var.tectonic_openstack_credentials_endpoint_type}"
  insecure            = "${var.tectonic_openstack_credentials_insecure}"
  key                 = "${var.tectonic_openstack_credentials_key}"
  password            = "${var.tectonic_openstack_credentials_password}"
  project_domain_id   = "${var.tectonic_openstack_credentials_project_domain_id}"
  project_domain_name = "${var.tectonic_openstack_credentials_project_domain_name}"
  region              = "${var.tectonic_openstack_region}"
  region              = "${var.tectonic_openstack_credentials_region}"
  swauth              = "${var.tectonic_openstack_credentials_swauth}"
  tenant_id           = "${var.tectonic_openstack_credentials_tenant_id}"
  tenant_name         = "${var.tectonic_openstack_credentials_tenant_name}"
  token               = "${var.tectonic_openstack_credentials_token}"
  use_octavia         = "${var.tectonic_openstack_credentials_use_octavia}"
  user_domain_id      = "${var.tectonic_openstack_credentials_user_domain_id}"
  user_domain_name    = "${var.tectonic_openstack_credentials_user_domain_name}"
  user_id             = "${var.tectonic_openstack_credentials_user_id}"
  user_name           = "${var.tectonic_openstack_credentials_user_name}"
  version             = ">=1.6.0"
}

module "bootstrap" {
  source = "./bootstrap"

  swift_container   = "${openstack_objectstorage_container_v1.tectonic.name}"
  cluster_name      = "${var.tectonic_cluster_name}"
  cluster_id        = "${var.tectonic_cluster_id}"
  image_name        = "${var.tectonic_openstack_base_image}"
  flavor_name       = "${var.tectonic_openstack_master_flavor_name}"
  ignition          = "${var.ignition_bootstrap}"
  bootstrap_port_id = "${module.topology.bootstrap_port_id}"
}

module "masters" {
  source = "./masters"

  base_image     = "${var.tectonic_openstack_base_image}"
  cluster_id     = "${var.tectonic_cluster_id}"
  cluster_name   = "${var.tectonic_cluster_name}"
  flavor_name    = "${var.tectonic_openstack_master_flavor_name}"
  instance_count = "${var.tectonic_master_count}"
  master_sg_ids  = "${concat(var.tectonic_openstack_master_extra_sg_ids, list(module.topology.master_sg_id))}"
  subnet_ids     = "${module.topology.master_subnet_ids}"
  user_data_igns = ["${var.ignition_masters}"]
}

# TODO(shadower) add a dns module here

module "topology" {
  source = "./topology"

  cidr_block                 = "${var.tectonic_openstack_network_cidr_block}"
  cluster_id                 = "${var.tectonic_cluster_id}"
  cluster_name               = "${var.tectonic_cluster_name}"
  external_master_subnet_ids = "${compact(var.tectonic_openstack_external_master_subnet_ids)}"
  external_network           = "${var.tectonic_openstack_external_network}"
  masters_count              = "${var.tectonic_master_count}"
}

resource "openstack_objectstorage_container_v1" "tectonic" {
  name = "${lower(var.tectonic_cluster_name)}.${var.tectonic_base_domain}"

  metadata = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-master",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_openstack_extra_tags)}"
}
