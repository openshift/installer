locals {
  description = "Created By OpenShift Installer"
  tags = merge(
    {
      "OCP"                                     = "ISV Integration",
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.resource_tags,
  )
}

provider "alicloud" {
  access_key = var.access_key
  secret_key = var.secret_key
  region     = var.region_id
}

module "vpc" {
  source            = "./vpc"
  cluster_id        = var.cluster_id
  region_id         = var.region_id
  zone_ids          = var.zone_ids
  resource_group_id = var.resource_group_id
  vpc_cidr_block    = var.machine_v4_cidrs[0]
  tags              = local.tags
}

module "pvtz" {
  source            = "./privatezone"
  cluster_id        = var.cluster_id
  resource_group_id = var.resource_group_id
  vpc_id            = module.vpc.vpc_id
  cluster_domain    = var.cluster_domain
  base_domain       = var.base_domain
  slb_external_ip   = module.vpc.slb_external_ip
  slb_internal_ip   = module.vpc.slb_internal_ip
  tags              = local.tags
}

module "ram" {
  source     = "./ram"
  cluster_id = var.cluster_id
  tags       = local.tags
}

module "master" {
  source               = "./master"
  cluster_id           = var.cluster_id
  resource_group_id    = var.resource_group_id
  vpc_id               = module.vpc.vpc_id
  vswitch_ids          = module.vpc.vswitch_ids
  sg_id                = module.vpc.sg_master_id
  slb_id               = module.vpc.slb_external_id
  instance_type        = var.instance_type
  image_id             = var.image_id
  system_disk_size     = var.system_disk_size
  system_disk_category = var.system_disk_category
  user_data_ign        = var.ignition_master
  key_name             = var.key_name
  role_name            = module.ram.role_master_name
  tags                 = local.tags
}

module "bootstrap" {
  source               = "./bootstrap"
  cluster_id           = var.cluster_id
  resource_group_id    = var.resource_group_id
  ignition_file        = var.ignition_bootstrap_file
  ignition_bucket      = var.ignition_bucket
  ignition             = var.ignition_bootstrap
  vpc_id               = module.vpc.vpc_id
  vswitch_id           = module.vpc.vswitch_ids[0]
  slb_id               = module.vpc.slb_external_id
  instance_type        = var.instance_type
  image_id             = var.image_id
  system_disk_size     = var.system_disk_size
  system_disk_category = var.system_disk_category
  key_name             = var.key_name
  tags                 = local.tags
}
