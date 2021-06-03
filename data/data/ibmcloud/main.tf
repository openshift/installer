locals {
  description       = "Created By OpenShift Installer"
  resource_group_id = var.ibmcloud_resource_group_name == "" ? ibm_resource_group.group.0.id : data.ibm_resource_group.group.0.id
  tags = concat(
    [ "kubernetes.io_cluster_${var.cluster_id}:owned" ],
    var.ibmcloud_extra_tags
  )
}

############################################
# IBM Cloud provider
############################################

provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
  region           = var.ibmcloud_region
}

############################################
# Resource group
############################################

resource "ibm_resource_group" "group" {
  count = var.ibmcloud_resource_group_name == "" ? 1 : 0
  name  = var.cluster_id
}

data "ibm_resource_group" "group" {
  count = var.ibmcloud_resource_group_name == "" ? 0 : 1
  name  = var.ibmcloud_resource_group_name
}

############################################
# Shared COS Instance
############################################
resource "ibm_resource_instance" "cos" {
  name              = "${var.cluster_id}-cos"
  service           = "cloud-object-storage"
  plan              = "standard"
  location          = "global"
  resource_group_id = local.resource_group_id
  tags              = local.tags
}

############################################
# Import VPC Custom Image
############################################

module "image" {
  source = "./image"

  name                     = "${var.cluster_id}-rhcos"
  image_filepath           = var.ibmcloud_image_filepath
  cluster_id               = var.cluster_id
  region                   = var.ibmcloud_region
  resource_group_id        = local.resource_group_id
  tags                     = local.tags
  cos_resource_instance_id = ibm_resource_instance.cos.id
}

############################################
# Bootstrap module
############################################

module "bootstrap" {
  source = "./bootstrap"
  
  cluster_id               = var.cluster_id
  cos_resource_instance_id = ibm_resource_instance.cos.id
  cos_bucket_region        = var.ibmcloud_region
  ignition_file            = var.ignition_bootstrap_file
  resource_group_id        = local.resource_group_id
  security_group_id        = module.vpc.control_plane_security_group_id
  subnet_id                = module.vpc.control_plane_subnet_id_list[0]
  tags                     = local.tags
  vpc_id                   = module.vpc.vpc_id
  vsi_image_id             = module.image.vsi_image_id
  vsi_profile              = var.ibmcloud_bootstrap_instance_type
  zone                     = module.vpc.control_plane_subnet_zone_list[0]

  lb_kubernetes_api_public_id       = module.vpc.lb_kubernetes_api_public_id
  lb_kubernetes_api_private_id      = module.vpc.lb_kubernetes_api_private_id
  lb_pool_kubernetes_api_public_id  = module.vpc.lb_pool_kubernetes_api_public_id
  lb_pool_kubernetes_api_private_id = module.vpc.lb_pool_kubernetes_api_private_id
  lb_pool_machine_config_id         = module.vpc.lb_pool_machine_config_id
}

############################################
# Master module
############################################

module "master" {
  source     = "./master"
  
  cluster_id        = var.cluster_id
  instance_count    = var.master_count
  ignition          = var.ignition_master
  resource_group_id = local.resource_group_id
  security_group_id = module.vpc.control_plane_security_group_id
  subnet_id_list    = module.vpc.control_plane_subnet_id_list
  tags              = local.tags
  vpc_id            = module.vpc.vpc_id
  vsi_image_id      = module.image.vsi_image_id
  vsi_profile       = var.ibmcloud_master_instance_type
  zone_list         = module.vpc.control_plane_subnet_zone_list

  lb_kubernetes_api_public_id       = module.vpc.lb_kubernetes_api_public_id
  lb_kubernetes_api_private_id      = module.vpc.lb_kubernetes_api_private_id
  lb_pool_kubernetes_api_public_id  = module.vpc.lb_pool_kubernetes_api_public_id
  lb_pool_kubernetes_api_private_id = module.vpc.lb_pool_kubernetes_api_private_id
  lb_pool_machine_config_id         = module.vpc.lb_pool_machine_config_id
}

############################################
# CIS module
############################################

module "cis" {
  source     = "./cis"
  
  cis_id         = var.ibmcloud_cis_crn
  base_domain    = var.base_domain
  cluster_domain = var.cluster_domain

  bootstrap_name         = module.bootstrap.name
  bootstrap_ipv4_address = module.bootstrap.primary_ipv4_address

  master_count             = var.master_count
  master_name_list         = module.master.name_list
  master_ipv4_address_list = module.master.primary_ipv4_address_list

  lb_kubernetes_api_public_hostname  = module.vpc.lb_kubernetes_api_public_hostname
  lb_kubernetes_api_private_hostname = module.vpc.lb_kubernetes_api_private_hostname
}

############################################
# VPC module
############################################

module "vpc" {
  source = "./vpc"
  
  cluster_id        = var.cluster_id
  resource_group_id = local.resource_group_id
  region            = var.ibmcloud_region
  tags              = local.tags
  zone_list         = distinct(var.ibmcloud_master_availability_zones)
}
