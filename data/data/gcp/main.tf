locals {
  labels = merge(
    {
      "kubernetes-io-cluster-${var.cluster_id}" = "owned"
    },
    var.gcp_extra_labels,
  )

  master_subnet_cidr = cidrsubnet(var.machine_v4_cidrs[0], 1, 0) #master subnet is a smaller subnet within the vnet. e.g., from /21 to /22
  worker_subnet_cidr = cidrsubnet(var.machine_v4_cidrs[0], 1, 1) #worker subnet is a smaller subnet within the vnet. e.g., from /21 to /22
  public_endpoints   = var.gcp_publish_strategy == "External" ? true : false

  gcp_image   = var.gcp_preexisting_image ? var.gcp_image : google_compute_image.cluster[0].self_link
  description = "Created By OpenShift Installer"
}

provider "google" {
  credentials = var.gcp_service_account
  project     = var.gcp_project_id
  region      = var.gcp_region
}

module "bootstrap" {
  source = "./bootstrap"
  count  = var.bootstrapping ? 1 : 0

  bootstrap_enabled = var.bootstrapping

  image            = local.gcp_image
  machine_type     = var.gcp_bootstrap_instance_type
  cluster_id       = var.cluster_id
  ignition         = var.ignition_bootstrap
  network          = module.network.network
  network_cidr     = var.machine_v4_cidrs[0]
  public_endpoints = local.public_endpoints
  subnet           = module.network.master_subnet
  zone             = var.gcp_master_availability_zones[0]
  region           = var.gcp_region

  root_volume_size         = var.gcp_master_root_volume_size
  root_volume_type         = var.gcp_master_root_volume_type
  root_volume_kms_key_link = var.gcp_root_volume_kms_key_link

  labels = local.labels
}

module "master" {
  source = "./master"

  image          = local.gcp_image
  instance_count = var.master_count
  machine_type   = var.gcp_master_instance_type
  cluster_id     = var.cluster_id
  ignition       = var.ignition_master
  subnet         = module.network.master_subnet
  zones          = distinct(var.gcp_master_availability_zones)

  root_volume_size         = var.gcp_master_root_volume_size
  root_volume_type         = var.gcp_master_root_volume_type
  root_volume_kms_key_link = var.gcp_root_volume_kms_key_link

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
  network_cidr       = var.machine_v4_cidrs[0]
  public_endpoints   = local.public_endpoints

  bootstrap_lb              = var.bootstrapping
  bootstrap_instances       = module.bootstrap.bootstrap_instances
  bootstrap_instance_groups = module.bootstrap.bootstrap_instance_groups

  master_instances       = module.master.master_instances
  master_instance_groups = module.master.master_instance_groups

  preexisting_network = var.gcp_preexisting_network
  cluster_network     = var.gcp_cluster_network
  master_subnet       = var.gcp_control_plane_subnet
  worker_subnet       = var.gcp_compute_subnet
}

module "dns" {
  source = "./dns"

  cluster_id           = var.cluster_id
  public_dns_zone_name = var.gcp_public_dns_zone_name
  network              = module.network.network
  cluster_domain       = var.cluster_domain
  api_external_lb_ip   = module.network.cluster_public_ip
  api_internal_lb_ip   = module.network.cluster_ip
  public_endpoints     = local.public_endpoints
}

resource "google_compute_image" "cluster" {
  count       = var.gcp_preexisting_image ? 0 : 1
  description = local.description

  name = "${var.cluster_id}-rhcos-image"

  # See https://github.com/openshift/installer/issues/2546
  guest_os_features {
    type = "SECURE_BOOT"
  }
  guest_os_features {
    type = "UEFI_COMPATIBLE"
  }

  raw_disk {
    source = var.gcp_image_uri
  }

  licenses = var.gcp_image_licenses
}
