locals {
  labels = var.gcp_extra_labels

  master_subnet_cidr = cidrsubnet(var.machine_cidr, 3, 0) #master subnet is a smaller subnet within the vnet. i.e from /21 to /24
  worker_subnet_cidr = cidrsubnet(var.machine_cidr, 3, 1) #worker subnet is a smaller subnet within the vnet. i.e from /21 to /24
}

provider "google" {
  credentials = var.gcp_service_account
  project     = var.gcp_project_id
  region      = var.gcp_region
}

module "bootstrap" {
  source = "./bootstrap"

  bootstrap_enabled = var.gcp_bootstrap_enabled

  image_name   = var.gcp_image_id
  machine_type = var.gcp_bootstrap_instance_type
  cluster_id   = var.cluster_id
  ignition     = var.ignition_bootstrap
  network      = module.network.network
  subnet       = module.network.master_subnet
  zone         = var.gcp_master_availability_zones[0]

  labels = local.labels
}

module "master" {
  source = "./master"

  image_name     = var.gcp_image_id
  instance_count = var.master_count
  machine_type   = var.gcp_master_instance_type
  cluster_id     = var.cluster_id
  ignition       = var.ignition_master
  network        = module.network.network
  subnet         = module.network.master_subnet
  zones          = distinct(var.gcp_master_availability_zones)

  labels = local.labels
}

module "iam" {
  source = "./iam"

  cluster_id = var.cluster_id
}

module "network" {
  source = "./network"

  cluster_id         = var.cluster_id
  master_subnet_cidr = local.master_subnet_cidr
  worker_subnet_cidr = local.worker_subnet_cidr
  network_cidr       = var.machine_cidr

  bootstrap_lb              = var.gcp_bootstrap_enabled && var.gcp_bootstrap_lb
  bootstrap_instances       = module.bootstrap.bootstrap_instances
  bootstrap_instance_groups = module.bootstrap.bootstrap_instance_groups

  master_instances       = module.master.master_instances
  master_instance_groups = module.master.master_instance_groups
}

module "dns" {
  source = "./dns"

  cluster_id           = var.cluster_id
  public_dns_zone_name = var.gcp_public_dns_zone_name
  network              = module.network.network
  etcd_ip_addresses    = flatten(module.master.ip_addresses)
  etcd_count           = var.master_count
  cluster_domain       = var.cluster_domain
  api_external_lb_ip   = module.network.cluster_public_ip
  api_internal_lb_ip   = module.network.cluster_private_ip
}
