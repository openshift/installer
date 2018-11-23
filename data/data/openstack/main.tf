provider "openstack" {
  auth_url            = "${var.openstack_credentials_auth_url}"
  cert                = "${var.openstack_credentials_cert}"
  cloud               = "${var.openstack_credentials_cloud}"
  domain_id           = "${var.openstack_credentials_domain_id}"
  domain_name         = "${var.openstack_credentials_domain_name}"
  endpoint_type       = "${var.openstack_credentials_endpoint_type}"
  insecure            = "${var.openstack_credentials_insecure}"
  key                 = "${var.openstack_credentials_key}"
  password            = "${var.openstack_credentials_password}"
  project_domain_id   = "${var.openstack_credentials_project_domain_id}"
  project_domain_name = "${var.openstack_credentials_project_domain_name}"
  region              = "${var.openstack_region}"
  region              = "${var.openstack_credentials_region}"
  swauth              = "${var.openstack_credentials_swauth}"
  tenant_id           = "${var.openstack_credentials_tenant_id}"
  tenant_name         = "${var.openstack_credentials_tenant_name}"
  token               = "${var.openstack_credentials_token}"
  use_octavia         = "${var.openstack_credentials_use_octavia}"
  user_domain_id      = "${var.openstack_credentials_user_domain_id}"
  user_domain_name    = "${var.openstack_credentials_user_domain_name}"
  user_id             = "${var.openstack_credentials_user_id}"
  user_name           = "${var.openstack_credentials_user_name}"
  version             = ">=1.6.0"
}

module "lb" {
  source = "./lb"

  swift_container   = "${openstack_objectstorage_container_v1.container.name}"
  cluster_name      = "${var.cluster_name}"
  cluster_id        = "${var.cluster_id}"
  cluster_domain    = "${var.base_domain}"
  image_name        = "${var.openstack_base_image}"
  flavor_name       = "${var.openstack_master_flavor_name}"
  ignition          = "${var.ignition_bootstrap}"
  lb_port_id        = "${module.topology.lb_port_id}"
  master_ips        = "${module.topology.master_ips}"
  master_port_names = "${module.topology.master_port_names}"
}

module "bootstrap" {
  source = "./bootstrap"

  swift_container   = "${openstack_objectstorage_container_v1.container.name}"
  cluster_name      = "${var.cluster_name}"
  cluster_id        = "${var.cluster_id}"
  image_name        = "${var.openstack_base_image}"
  flavor_name       = "${var.openstack_master_flavor_name}"
  ignition          = "${var.ignition_bootstrap}"
  bootstrap_port_id = "${module.topology.bootstrap_port_id}"
  service_vm_fixed_ip = "${module.topology.service_vm_fixed_ip}"
}

module "masters" {
  source = "./masters"

  base_image     = "${var.openstack_base_image}"
  cluster_id     = "${var.cluster_id}"
  cluster_name   = "${var.cluster_name}"
  flavor_name    = "${var.openstack_master_flavor_name}"
  instance_count = "${var.master_count}"
  master_sg_ids  = "${concat(var.openstack_master_extra_sg_ids, list(module.topology.master_sg_id))}"
  subnet_ids     = "${module.topology.master_subnet_ids}"
  user_data_ign  = "${var.ignition_master}"
  service_vm_fixed_ip = "${module.topology.service_vm_fixed_ip}"
}

# TODO(shadower) add a dns module here

module "topology" {
  source = "./topology"

  cidr_block                 = "${var.openstack_network_cidr_block}"
  cluster_id                 = "${var.cluster_id}"
  cluster_name               = "${var.cluster_name}"
  external_master_subnet_ids = "${compact(var.openstack_external_master_subnet_ids)}"
  external_network           = "${var.openstack_external_network}"
  masters_count              = "${var.master_count}"
}

resource "openstack_objectstorage_container_v1" "container" {
  name = "${lower(var.cluster_name)}.${var.base_domain}"

  metadata = "${merge(map(
      "Name", "${var.cluster_name}-ignition-master",
      "KubernetesCluster", "${var.cluster_name}",
      "tectonicClusterID", "${var.cluster_id}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.openstack_extra_tags)}"
}
