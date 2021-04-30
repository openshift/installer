locals {
  description = "Created By OpenShift Installer"
}

############################################
# IBM Cloud provider
############################################

provider "ibm" {
  region = var.ibmcloud_region
}

############################################
# Datasources
############################################

data "ibm_is_image" "vsi_image" {
  name = var.ibmcloud_vsi_image
}

############################################
# Bootstrap module
############################################

module "bootstrap" {
  source = "./bootstrap"
  
  cluster_id        = var.cluster_id
  cos_bucket_region = var.ibmcloud_region
  ignition_file     = var.ignition_bootstrap_file
  resource_group_id = local.resource_group_id
  security_group_id = module.vpc.control_plane_security_group_id
  subnet_id         = module.vpc.control_plane_subnet_id_list[0]
  vpc_id            = module.vpc.vpc_id
  vsi_image_id      = data.ibm_is_image.vsi_image.id
  vsi_profile       = var.ibmcloud_vsi_profile
  zone              = module.vpc.control_plane_subnet_zone_list[0]

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
  depends_on = [ module.bootstrap ]
  
  cluster_id        = var.cluster_id
  instance_count    = var.master_count
  ignition          = var.ignition_master
  resource_group_id = local.resource_group_id
  security_group_id = module.vpc.control_plane_security_group_id
  subnet_id_list    = module.vpc.control_plane_subnet_id_list
  vpc_id            = module.vpc.vpc_id
  vsi_image_id      = data.ibm_is_image.vsi_image.id
  vsi_profile       = var.ibmcloud_vsi_profile
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
  
  cis_id     = var.ibmcloud_cis_id
  cluster_id = var.cluster_id
  domain     = var.base_domain

  bootstrap_name         = module.bootstrap.name
  bootstrap_ipv4_address = module.bootstrap.primary_ipv4_address

  master_count             = var.master_count
  master_name_list         = module.master.name_list
  master_ipv4_address_list = module.master.primary_ipv4_address_list

  lb_kubernetes_api_public_hostname       = module.vpc.lb_kubernetes_api_public_hostname
  lb_kubernetes_api_private_hostname      = module.vpc.lb_kubernetes_api_private_hostname
  lb_application_ingress_public_hostname  = module.vpc.lb_application_ingress_public_hostname
  lb_application_ingress_private_hostname = module.vpc.lb_application_ingress_private_hostname
}

############################################
# VPC module
############################################

module "vpc" {
  source = "./vpc"
  
  cluster_id        = var.cluster_id
  resource_group_id = local.resource_group_id
  region            = var.ibmcloud_region
  zone_list         = distinct(var.ibmcloud_master_availability_zones)
}
