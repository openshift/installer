locals {
  compute_subnet_cidr = cidrsubnet(var.machine_cidr, 3, 1) #compute subnet is a smaller subnet within the vnet. i.e from /21 to /24
  control_subnet_cidr = cidrsubnet(var.machine_cidr, 3, 0) #control plane subnet is a smaller subnet within the vnet. i.e from /21 to /24
}

provider "google" {
  region  = var.google_region
  project = var.google_project
}

module "network" {
  source = "./network"

  cluster_id          = var.cluster_id
  compute_subnet_cidr = local.compute_subnet_cidr
  control_subnet_cidr = local.control_subnet_cidr
  network_cidr        = var.machine_cidr
}
