locals {
  description = "Created By OpenShift Installer"
  is_external = var.ali_publish_strategy == "External" ? true : false
  tags = merge(
    {
      "GISV"                                      = "ocp",
      "sigs.k8s.io/cloud-provider-alibaba/origin" = "ocp",
      "kubernetes.io/cluster/${var.cluster_id}"   = "owned",
      "ack.aliyun.com"                            = var.cluster_id
    },
    var.ali_extra_tags,
  )
}

provider "alicloud" {
  access_key = var.ali_access_key
  secret_key = var.ali_secret_key
  region     = var.ali_region_id
}

module "resource_group" {
  source            = "./resourcegroup"
  cluster_id        = var.cluster_id
  resource_group_id = var.ali_resource_group_id
}

module "vpc" {
  source      = "./vpc"
  vpc_id      = var.ali_vpc_id
  vswitch_ids = var.ali_vswitch_ids
  cluster_id  = var.cluster_id
  region_id   = var.ali_region_id
  zone_ids = distinct(
    concat(
      var.ali_master_availability_zone_ids,
      var.ali_worker_availability_zone_ids,
    ),
  )
  nat_gateway_zone_id = var.ali_nat_gateway_zone_id
  resource_group_id   = module.resource_group.resource_group_id
  vpc_cidr_block      = var.machine_v4_cidrs[0]
  tags                = local.tags
  publish_strategy    = var.ali_publish_strategy
}

module "dns" {
  source            = "./dns"
  cluster_id        = var.cluster_id
  private_zone_id   = var.ali_private_zone_id
  resource_group_id = module.resource_group.resource_group_id
  vpc_id            = module.vpc.vpc_id
  cluster_domain    = var.cluster_domain
  base_domain       = var.base_domain
  slb_external_ip   = module.vpc.slb_external_ip
  slb_internal_ip   = module.vpc.slb_internal_ip
  tags              = local.tags
  publish_strategy  = var.ali_publish_strategy
}

module "ram" {
  source     = "./ram"
  cluster_id = var.cluster_id
  tags       = local.tags
}

module "master" {
  source               = "./master"
  cluster_id           = var.cluster_id
  resource_group_id    = module.resource_group.resource_group_id
  vpc_id               = module.vpc.vpc_id
  zone_ids             = var.ali_master_availability_zone_ids
  az_to_vswitch_id     = module.vpc.az_to_vswitch_id
  sg_id                = module.vpc.sg_master_id
  slb_ids              = module.vpc.slb_ids
  slb_group_length     = module.vpc.slb_group_length
  instance_type        = var.ali_master_instance_type
  instance_count       = var.master_count
  image_id             = var.ali_image_id
  system_disk_size     = var.ali_system_disk_size
  system_disk_category = var.ali_system_disk_category
  user_data_ign        = var.ignition_master
  role_name            = module.ram.role_master_name
  tags                 = local.tags
  publish_strategy     = var.ali_publish_strategy
}
