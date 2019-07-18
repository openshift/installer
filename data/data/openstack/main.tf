provider "openstack" {
  auth_url            = var.openstack_credentials_auth_url
  cert                = var.openstack_credentials_cert
  cloud               = var.openstack_credentials_cloud
  domain_id           = var.openstack_credentials_domain_id
  domain_name         = var.openstack_credentials_domain_name
  endpoint_type       = var.openstack_credentials_endpoint_type
  insecure            = var.openstack_credentials_insecure
  key                 = var.openstack_credentials_key
  password            = var.openstack_credentials_password
  project_domain_id   = var.openstack_credentials_project_domain_id
  project_domain_name = var.openstack_credentials_project_domain_name
  region              = var.openstack_credentials_region
  swauth              = var.openstack_credentials_swauth
  tenant_id           = var.openstack_credentials_tenant_id
  tenant_name         = var.openstack_credentials_tenant_name
  token               = var.openstack_credentials_token
  use_octavia         = var.openstack_credentials_use_octavia
  user_domain_id      = var.openstack_credentials_user_domain_id
  user_domain_name    = var.openstack_credentials_user_domain_name
  user_id             = var.openstack_credentials_user_id
  user_name           = var.openstack_credentials_user_name
}

module "service" {
  source = "./service"

  swift_container   = openstack_objectstorage_container_v1.container.name
  cluster_id        = var.cluster_id
  cluster_domain    = var.cluster_domain
  image_name        = var.openstack_base_image
  flavor_name       = var.openstack_master_flavor_name
  ignition          = var.ignition_bootstrap
  lb_floating_ip    = var.openstack_lb_floating_ip
  service_port_id   = module.topology.service_port_id
  service_port_ip   = module.topology.service_port_ip
  master_ips        = module.topology.master_ips
  master_port_names = module.topology.master_port_names
  bootstrap_ip      = module.topology.bootstrap_port_ip
}

module "bootstrap" {
  source = "./bootstrap"

  swift_container     = openstack_objectstorage_container_v1.container.name
  cluster_id          = var.cluster_id
  cluster_domain      = var.cluster_domain
  image_name          = var.openstack_base_image
  flavor_name         = var.openstack_master_flavor_name
  ignition            = var.ignition_bootstrap
  bootstrap_port_id   = module.topology.bootstrap_port_id
  service_vm_fixed_ip = module.topology.service_vm_fixed_ip
}

module "masters" {
  source = "./masters"

  base_image          = var.openstack_base_image
  bootstrap_ip        = module.topology.bootstrap_port_ip
  cluster_id          = var.cluster_id
  cluster_domain      = var.cluster_domain
  flavor_name         = var.openstack_master_flavor_name
  instance_count      = var.master_count
  lb_floating_ip      = var.openstack_lb_floating_ip
  master_ips          = module.topology.master_ips
  master_port_ids     = module.topology.master_port_ids
  master_port_names   = module.topology.master_port_names
  user_data_ign       = var.ignition_master
  service_vm_fixed_ip = module.topology.service_vm_fixed_ip
  api_int_ip          = var.openstack_api_int_ip
  node_dns_ip         = var.openstack_node_dns_ip
  master_sg_ids = concat(
    var.openstack_master_extra_sg_ids,
    [module.topology.master_sg_id],
  )
}

# TODO(shadower) add a dns module here

module "topology" {
  source = "./topology"

  cidr_block          = var.machine_cidr
  cluster_id          = var.cluster_id
  external_network    = var.openstack_external_network
  external_network_id = var.openstack_external_network_id
  masters_count       = var.master_count
  lb_floating_ip      = var.openstack_lb_floating_ip
  api_int_ip          = var.openstack_api_int_ip
  node_dns_ip         = var.openstack_node_dns_ip
  trunk_support       = var.openstack_trunk_support
  octavia_support     = var.openstack_octavia_support
}

resource "openstack_objectstorage_container_v1" "container" {
  name = var.cluster_id

  # "kubernetes.io/cluster/${var.cluster_id}" = "owned"
  metadata = merge(
    {
      "Name"               = "${var.cluster_id}-ignition-master"
      "openshiftClusterID" = var.cluster_id
    },
    var.openstack_extra_tags,
  )
}

