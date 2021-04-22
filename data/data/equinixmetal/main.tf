provider "packet" {
  auth_token = var.metal_auth_token
}

module "bootstrap" {
  source     = "./bootstrap"
  project_id = var.metal_project_id
  facility   = var.metal_facility
  // metro                = var.metal_metro
  billing_cycle        = var.metal_billing_cycle
  plan                 = var.metal_machine_type
  operating_system     = var.metal_bootstrap_operating_system
  ssh_private_key_path = var.metal_ssh_private_key_path
  cluster_domain       = var.cluster_domain
  base_domain          = var.base_domain
  // cf_zone_id           = var.cf_zone_id
  //ocp_version         = var.metal_ocp_version
  //ocp_version_zstream = var.metal_ocp_version_zstream
  //depends              = [module.prepare_openshift.finished]

  ip_addresses = module.dns.boostrap_a

  ignition = var.ignition_bootstrap
}


module "master" {
  node_count = var.master_count
  source     = "./master"

  project_id = var.metal_project_id
  facility   = var.metal_facility
  // metro = var.metal_metro
  // billing_cycle        = var.metal_billing_cycle
  plan                 = var.metal_machine_type
  operating_system     = "custom_ipxe"
  ssh_private_key_path = var.metal_ssh_private_key_path
  cluster_domain       = var.cluster_domain
  base_domain          = var.base_domain
  // cf_zone_id           = var.cf_zone_id
  //ocp_version         = var.metal_ocp_version
  //ocp_version_zstream = var.metal_ocp_version_zstream
  //depends              = [module.prepare_openshift.finished]

  ip_addresses = module.dns.masters_a

  ignition = var.ignition_master

  bootstrap_ip = module.bootstrap.lb_ip
}

module "dns" {
  source = "./dns"

  cluster_name       = var.cluster_name
  cluster_basedomain = var.cluster_basedomain
}
/*
module "prepare_openshift" {

  source = "./modules/prereq"

  cluster_name         = var.cluster_name
  cluster_basedomain   = var.cluster_basedomain
  ocp_version          = var.ocp_version
  count_master         = var.count_master
  count_compute        = var.count_compute
  ssh_public_key_path  = var.ssh_public_key_path
  ssh_private_key_path = var.ssh_private_key_path
  bootstrap_ip           = module.bootstrap.lb_ip
  ocp_api_token        = var.ocp_cluster_manager_token
  depends              = [module.bootstrap.finished]
}

module "openshift_install" {
  source = "./modules/install"

  ssh_private_key_path = var.ssh_private_key_path
  operating_system     = var.bootstrap_operating_system
  bootstrap_ip           = module.bootstrap.lb_ip
  count_master         = var.count_master
  count_compute        = var.count_compute
  cluster_name         = var.cluster_name
  cluster_basedomain   = var.cluster_basedomain
  bootstrap_ip         = module.openshift_bootstrap.node_ip
  master_ips           = module.openshift_masters.node_ip
  worker_ips           = module.openshift_workers.node_ip
  depends              = [module.openshift_masters.node_ip, module.openshift_workers.node_ip]

  ocp_storage_nfs_enable    = var.ocp_storage_nfs_enable
  ocp_storage_ocs_enable    = var.ocp_storage_ocs_enable
  ocp_virtualization_enable = var.ocp_virtualization_enable
}
*/
