locals {
  master_subnet_cidr = cidrsubnet(var.machine_v4_cidrs[0], 1, 0) #master subnet is a smaller subnet within the vnet. e.g., from /21 to /22
  worker_subnet_cidr = cidrsubnet(var.machine_v4_cidrs[0], 1, 1) #worker subnet is a smaller subnet within the vnet. e.g., from /21 to /22
  public_endpoints   = var.gcp_publish_strategy == "External" ? true : false

  gcp_image   = var.gcp_image
  description = "Created By OpenShift Installer"
}

provider "google" {
  credentials = var.gcp_service_account
  project     = var.gcp_project_id
  region      = var.gcp_region
}

module "master" {
  source = "./master"

  image           = local.gcp_image
  instance_count  = var.master_count
  machine_type    = var.gcp_master_instance_type
  project_id      = var.gcp_project_id
  cluster_id      = var.cluster_id
  service_account = var.gcp_instance_service_account
  ignition        = var.ignition_master
  subnet          = module.network.master_subnet
  zones           = distinct(var.gcp_master_availability_zones)
  secure_boot     = var.gcp_master_secure_boot

  root_volume_size = var.gcp_master_root_volume_size
  root_volume_type = var.gcp_master_root_volume_type

  root_volume_kms_key_link = var.gcp_root_volume_kms_key_link

  confidential_compute = var.gcp_master_confidential_compute
  on_host_maintenance  = var.gcp_master_on_host_maintenance
  gcp_extra_labels     = var.gcp_extra_labels
  gcp_extra_tags       = var.gcp_extra_tags

  tags = var.gcp_control_plane_tags
}

module "iam" {
  source = "./iam"

  project_id = var.gcp_project_id
  cluster_id = var.cluster_id

  service_account = var.gcp_instance_service_account
}

module "network" {
  source = "./network"

  cluster_id         = var.cluster_id
  master_subnet_cidr = local.master_subnet_cidr
  worker_subnet_cidr = local.worker_subnet_cidr
  network_cidr       = var.machine_v4_cidrs[0]
  public_endpoints   = local.public_endpoints

  preexisting_network = var.gcp_preexisting_network
  cluster_network     = var.gcp_cluster_network
  master_subnet       = var.gcp_control_plane_subnet
  worker_subnet       = var.gcp_compute_subnet
  network_project_id  = var.gcp_network_project_id

  create_firewall_rules = var.gcp_create_firewall_rules
}

module "dns" {
  count = var.gcp_user_provisioned_dns ? 0 : 1

  source = "./dns"

  cluster_id         = var.cluster_id
  public_zone_name   = var.gcp_public_zone_name
  private_zone_name  = var.gcp_private_zone_name
  network            = module.network.network
  cluster_domain     = var.cluster_domain
  api_external_lb_ip = module.network.cluster_public_ip
  api_internal_lb_ip = module.network.cluster_ip
  public_endpoints   = local.public_endpoints
  project_id         = var.gcp_project_id
  gcp_extra_labels   = var.gcp_extra_labels
}
