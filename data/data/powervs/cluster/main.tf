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

  cluster_id           = var.cluster_id
  publish_strategy     = var.powervs_publish_strategy
  resource_group       = var.powervs_resource_group
  vpc_zone             = var.powervs_vpc_zone
  vpc_subnet_name      = var.powervs_vpc_subnet_name
  vpc_name             = var.powervs_vpc_name
  vpc_gateway_name     = var.powervs_vpc_gateway_name
  vpc_gateway_attached = var.powervs_vpc_gateway_attached
  wait_for_vpc         = var.powervs_wait_for_vpc
}

module "pi_network" {
  providers = {
    ibm = ibm.powervs
  }
  source = "./power_network"

  cluster_id        = var.cluster_id
  cloud_instance_id = module.iaas.si_guid
  resource_group    = var.powervs_resource_group
  machine_cidr      = var.machine_v4_cidrs[0]
  vpc_crn           = module.vpc.vpc_crn
  dns_server        = module.dns.dns_server
  enable_snat       = var.powervs_enable_snat
}

resource "ibm_pi_key" "cluster_key" {
  provider             = ibm.powervs
  pi_key_name          = "${var.cluster_id}-key"
  pi_ssh_key           = var.powervs_ssh_key
  pi_cloud_instance_id = module.iaas.si_guid
}

module "master" {
  providers = {
    ibm = ibm.powervs
  }
  source = "./master"

  cloud_instance_id   = module.iaas.si_guid
  cluster_id          = var.cluster_id
  resource_group      = var.powervs_resource_group
  instance_count      = var.master_count
  api_key             = var.powervs_api_key
  powervs_region      = var.powervs_region
  powervs_zone        = var.powervs_zone
  vpc_region          = var.powervs_vpc_region
  vpc_zone            = module.vpc.vpc_zone
  memory              = var.powervs_master_memory
  processors          = var.powervs_master_processors
  ignition            = var.ignition_master
  sys_type            = var.powervs_sys_type
  proc_type           = var.powervs_proc_type
  ssh_key_name        = ibm_pi_key.cluster_key.name
  image_id            = ibm_pi_image.boot_image.image_id
  dhcp_network_id     = module.pi_network.dhcp_network_id
  dhcp_id             = module.pi_network.dhcp_id
  lb_ext_id           = module.loadbalancer.lb_ext_id
  lb_int_id           = module.loadbalancer.lb_int_id
  machine_cfg_pool_id = module.loadbalancer.machine_cfg_pool_id
  api_pool_int_id     = module.loadbalancer.api_pool_int_id
  api_pool_ext_id     = module.loadbalancer.api_pool_ext_id
}

resource "ibm_pi_image" "boot_image" {
  provider                  = ibm.powervs
  pi_image_name             = "rhcos-${var.cluster_id}"
  pi_cloud_instance_id      = module.iaas.si_guid
  pi_image_bucket_name      = var.powervs_image_bucket_name
  pi_image_bucket_access    = "public"
  pi_image_bucket_region    = var.powervs_cos_region
  pi_image_bucket_file_name = var.powervs_image_bucket_file_name
  pi_image_storage_type     = var.powervs_image_storage_type
}

data "ibm_pi_dhcp" "dhcp_service" {
  provider             = ibm.powervs
  depends_on           = [module.master]
  pi_cloud_instance_id = module.iaas.si_guid
  pi_dhcp_id           = module.pi_network.dhcp_id
}

module "loadbalancer" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./loadbalancer"

  cluster_id     = var.cluster_id
  enable_snat    = var.powervs_enable_snat
  master_count   = var.master_count
  resource_group = var.powervs_resource_group
  vpc_id         = module.vpc.vpc_id
  vpc_subnet_id  = module.vpc.vpc_subnet_id
}

locals {
  dns_service_id = var.powervs_publish_strategy == "Internal" ? var.powervs_dns_guid : var.powervs_cis_crn
}

module "dns" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./dns"

  service_id                 = local.dns_service_id
  base_domain                = var.base_domain
  cluster_domain             = var.cluster_domain
  load_balancer_hostname     = module.loadbalancer.lb_hostname
  load_balancer_int_hostname = module.loadbalancer.lb_int_hostname
  cluster_id                 = var.cluster_id
  vpc_crn                    = module.vpc.vpc_crn
  vpc_id                     = module.vpc.vpc_id
  vpc_subnet_id              = module.vpc.vpc_subnet_id
  vpc_zone                   = module.vpc.vpc_zone
  vpc_region                 = var.powervs_vpc_region
  vpc_permitted              = var.powervs_vpc_permitted
  ssh_key                    = var.powervs_ssh_key
  publish_strategy           = var.powervs_publish_strategy
  enable_snat                = var.powervs_enable_snat
  # dns_vm_image_name        = @FUTURE
}

module "transit_gateway" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./transit_gateway"

  cluster_id               = var.cluster_id
  resource_group           = var.powervs_resource_group
  service_instance_crn     = module.iaas.si_crn
  attached_transit_gateway = var.powervs_attached_transit_gateway
  tg_connection_vpc_id     = var.powervs_tg_connection_vpc_id
  vpc_crn                  = module.vpc.vpc_crn
  vpc_region               = var.powervs_vpc_region
}

module "iaas" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./iaas"

  #
  # define and pass variables to:
  # data/data/powervs/cluster/iaas/variables.tf
  #
  cluster_id            = var.cluster_id
  resource_group        = var.powervs_resource_group
  powervs_zone          = var.powervs_zone
  service_instance_name = var.powervs_service_instance_name
  wait_for_workspace    = var.powervs_wait_for_workspace
}
