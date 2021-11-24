locals {
  resource_group_id = var.ibmcloud_resource_group_name == "" ? ibm_resource_group.group.0.id : data.ibm_resource_group.group.0.id
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

  name                      = "${var.cluster_id}-rhcos"
  image_filepath            = var.ibmcloud_image_filepath
  cluster_id                = var.cluster_id
  region                    = var.ibmcloud_region
  resource_group_id         = local.resource_group_id
  tags                      = local.tags
  cos_resource_instance_crn = ibm_resource_instance.cos.crn
}

############################################
# CIS module
############################################

module "cis" {
  source = "./cis"

  cis_id         = var.ibmcloud_cis_crn
  base_domain    = var.base_domain
  cluster_domain = var.cluster_domain

  lb_kubernetes_api_public_hostname  = module.vpc.lb_kubernetes_api_public_hostname
  lb_kubernetes_api_private_hostname = module.vpc.lb_kubernetes_api_private_hostname
}

############################################
# Dedicated Host module
############################################

module "dhost" {
  source = "./dhost"

  cluster_id             = var.cluster_id
  dedicated_hosts_master = var.ibmcloud_master_dedicated_hosts
  dedicated_hosts_worker = var.ibmcloud_worker_dedicated_hosts
  resource_group_id      = local.resource_group_id
  zones_master           = distinct(var.ibmcloud_master_availability_zones)
  zones_worker           = distinct(var.ibmcloud_worker_availability_zones)
}

############################################
# VPC module
############################################

module "vpc" {
  source = "./vpc"

  cluster_id        = var.cluster_id
  public_endpoints  = local.public_endpoints
  resource_group_id = local.resource_group_id
  tags              = local.tags
  zones_master      = distinct(var.ibmcloud_master_availability_zones)
  zones_worker      = distinct(var.ibmcloud_worker_availability_zones)
}
