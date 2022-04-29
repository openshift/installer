provider "ibm" {
  alias            = "vpc"
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_vpc_region
  zone             = var.powervs_vpc_zone
}

provider "ibm" {
  alias            = "powervs"
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_region
  zone             = var.powervs_zone
}

module "vpc" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./vpc"

  cluster_id     = var.cluster_id
  resource_group = var.powervs_resource_group
  vpc_zone       = var.powervs_vpc_zone
}

module "pi_network" {
  providers = {
    ibm = ibm.powervs
  }
  source = "./power_network"

  cluster_id        = var.cluster_id
  cloud_instance_id = var.powervs_cloud_instance_id
  resource_group    = var.powervs_resource_group
  vpc_crn           = module.vpc.vpc_crn
}

resource "ibm_pi_key" "cluster_key" {
  provider             = ibm.powervs
  pi_key_name          = "${var.cluster_id}-key"
  pi_ssh_key           = var.powervs_ssh_key
  pi_cloud_instance_id = var.powervs_cloud_instance_id
}

module "bootstrap" {
  providers = {
    ibm = ibm.powervs
  }
  source            = "./bootstrap"
  cloud_instance_id = var.powervs_cloud_instance_id
  cluster_id        = var.cluster_id
  resource_group    = var.powervs_resource_group

  api_key               = var.powervs_api_key
  powervs_region        = var.powervs_region
  powervs_zone          = var.powervs_zone
  vpc_region            = var.powervs_vpc_region
  vpc_zone              = var.powervs_vpc_zone
  cos_instance_location = var.powervs_cos_instance_location
  cos_bucket_location   = var.powervs_cos_bucket_location
  cos_storage_class     = var.powervs_cos_storage_class
  memory                = var.powervs_bootstrap_memory
  processors            = var.powervs_bootstrap_processors
  ignition              = var.ignition_bootstrap
  sys_type              = var.powervs_sys_type
  proc_type             = var.powervs_proc_type
  ssh_key_id            = ibm_pi_key.cluster_key.key_id
  image_id              = ibm_pi_image.boot_image.image_id
  dhcp_network_id       = module.pi_network.dhcp_network_id
  dhcp_id               = module.pi_network.dhcp_id
  vpc_id                = module.vpc.vpc_id
  lb_ext_id             = module.loadbalancer.lb_ext_id
  lb_int_id             = module.loadbalancer.lb_int_id
  machine_cfg_pool_id   = module.loadbalancer.machine_cfg_pool_id
  api_pool_int_id       = module.loadbalancer.api_pool_int_id
  api_pool_ext_id       = module.loadbalancer.api_pool_ext_id
}

module "master" {
  providers = {
    ibm = ibm.powervs
  }
  source            = "./master"
  cloud_instance_id = var.powervs_cloud_instance_id
  cluster_id        = var.cluster_id
  resource_group    = var.powervs_resource_group
  instance_count    = var.master_count

  api_key                     = var.powervs_api_key
  powervs_region              = var.powervs_region
  powervs_zone                = var.powervs_zone
  vpc_region                  = var.powervs_vpc_region
  vpc_zone                    = var.powervs_vpc_zone
  memory                      = var.powervs_master_memory
  processors                  = var.powervs_master_processors
  ignition                    = var.ignition_master
  sys_type                    = var.powervs_sys_type
  proc_type                   = var.powervs_proc_type
  ssh_key_id                  = ibm_pi_key.cluster_key.key_id
  image_id                    = ibm_pi_image.boot_image.image_id
  dhcp_network_id             = module.pi_network.dhcp_network_id
  dhcp_id                     = module.pi_network.dhcp_id
  lb_ext_id                   = module.loadbalancer.lb_ext_id
  lb_int_id                   = module.loadbalancer.lb_int_id
  machine_cfg_pool_id         = module.loadbalancer.machine_cfg_pool_id
  api_pool_int_id             = module.loadbalancer.api_pool_int_id
  api_pool_ext_id             = module.loadbalancer.api_pool_ext_id
  bootstrap_api_member_int_id = module.bootstrap.api_member_int_id
  bootstrap_api_member_ext_id = module.bootstrap.api_member_ext_id
}

resource "ibm_pi_image" "boot_image" {
  provider                  = ibm.powervs
  pi_image_name             = "rhcos-${var.cluster_id}"
  pi_cloud_instance_id      = var.powervs_cloud_instance_id
  pi_image_bucket_name      = "rhcos-powervs-images-${var.powervs_vpc_region}"
  pi_image_bucket_access    = "public"
  pi_image_bucket_region    = var.powervs_vpc_region
  pi_image_bucket_file_name = var.powervs_image_bucket_file_name
  pi_image_storage_type     = var.powervs_image_storage_type
}

data "ibm_pi_dhcp" "dhcp_service" {
  provider             = ibm.powervs
  depends_on           = [module.bootstrap, module.master]
  pi_cloud_instance_id = var.powervs_cloud_instance_id
  pi_dhcp_id           = module.pi_network.dhcp_id
}

module "loadbalancer" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./loadbalancer"

  cluster_id     = var.cluster_id
  master_count   = var.master_count
  resource_group = var.powervs_resource_group
  vpc_id         = module.vpc.vpc_id
  vpc_subnet_id  = module.vpc.vpc_subnet_id
}

module "dns" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./dns"

  cis_id                     = var.powervs_cis_crn
  base_domain                = var.base_domain
  cluster_domain             = var.cluster_domain
  load_balancer_hostname     = module.loadbalancer.lb_hostname
  load_balancer_int_hostname = module.loadbalancer.lb_int_hostname
}
