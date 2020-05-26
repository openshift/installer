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

module "bootstrap" {
  source = "./bootstrap"

  cluster_id              = var.cluster_id
  extra_tags              = var.openstack_extra_tags
  base_image_id           = data.openstack_images_image_v2.base_image.id
  flavor_name             = var.openstack_master_flavor_name
  ignition                = var.ignition_bootstrap
  api_int_ip              = var.openstack_api_int_ip
  node_dns_ip             = var.openstack_node_dns_ip
  external_network        = var.openstack_external_network
  cluster_domain          = var.cluster_domain
  nodes_subnet_id         = module.topology.nodes_subnet_id
  private_network_id      = module.topology.private_network_id
  master_sg_id            = module.topology.master_sg_id
  bootstrap_shim_ignition = var.openstack_bootstrap_shim_ignition
  master_port_ids         = module.topology.master_port_ids
  root_volume_size        = var.openstack_master_root_volume_size
  root_volume_type        = var.openstack_master_root_volume_type
}

module "masters" {
  source = "./masters"

  base_image_id   = data.openstack_images_image_v2.base_image.id
  cluster_id      = var.cluster_id
  flavor_name     = var.openstack_master_flavor_name
  instance_count  = var.master_count
  master_port_ids = module.topology.master_port_ids
  user_data_ign   = var.ignition_master
  master_sg_ids = concat(
    var.openstack_master_extra_sg_ids,
    [module.topology.master_sg_id],
  )
  root_volume_size       = var.openstack_master_root_volume_size
  root_volume_type       = var.openstack_master_root_volume_type
  server_group_name      = var.openstack_master_server_group_name
  additional_network_ids = var.openstack_additional_network_ids
}

module "topology" {
  source = "./topology"

  cidr_block          = var.machine_v4_cidrs[0]
  cluster_id          = var.cluster_id
  cluster_domain      = var.cluster_domain
  external_network    = var.openstack_external_network
  external_network_id = var.openstack_external_network_id
  masters_count       = var.master_count
  lb_floating_ip      = var.openstack_lb_floating_ip
  api_int_ip          = var.openstack_api_int_ip
  node_dns_ip         = var.openstack_node_dns_ip
  ingress_ip          = var.openstack_ingress_ip
  external_dns        = var.openstack_external_dns
  trunk_support       = var.openstack_trunk_support
  octavia_support     = var.openstack_octavia_support
  machines_subnet_id  = var.openstack_machines_subnet_id
  machines_network_id = var.openstack_machines_network_id
}

data "openstack_images_image_v2" "base_image" {
  name = var.openstack_base_image_name
}
