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

  image_name   = var.gcp_image_id
  machine_type = var.gcp_bootstrap_instance_type
  cluster_id   = var.cluster_id
  ignition     = var.ignition_bootstrap
  network      = module.network.network
  subnet       = module.network.master_subnet
  zone         = module.network.zones[0]

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
  zones          = module.network.zones

  labels = local.labels
}

module "network" {
  source = "./network"

  cluster_id         = var.cluster_id
  master_subnet_cidr = local.master_subnet_cidr
  worker_subnet_cidr = local.worker_subnet_cidr
  network_cidr       = var.machine_cidr

  bootstrap_instance = module.bootstrap.bootstrap_instance
  master_instances   = module.master.master_instances

  bootstrap_instance_group = module.bootstrap.bootstrap_instance_group
  master_instance_groups   = module.master.master_instance_groups
}
