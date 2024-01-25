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

module "vm" {
  providers = {
    ibm = ibm.powervs
  }
  source = "./vm"

  resource_group        = var.powervs_resource_group
  cluster_id            = var.cluster_id
  ssh_key_name          = var.cluster_key_name
  cos_bucket_location   = var.powervs_vpc_region
  cos_instance_location = var.powervs_cos_instance_location
  cos_storage_class     = var.powervs_cos_storage_class
  ignition              = var.ignition_bootstrap
  memory                = var.powervs_bootstrap_memory
  processors            = var.powervs_bootstrap_processors
  proc_type             = var.powervs_proc_type
  image_id              = var.boot_image_id
  sys_type              = var.powervs_sys_type
  cloud_instance_id     = module.iaas.si_guid
  dhcp_network_id       = var.dhcp_network_id
  dhcp_id               = var.dhcp_id
  proxy_server_ip       = var.proxy_server_ip
  enable_snat           = var.powervs_enable_snat
}

module "lb" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./lb"

  bootstrap_ip        = module.vm.bootstrap_ip
  lb_int_id           = var.lb_int_id
  lb_ext_id           = var.lb_ext_id
  machine_cfg_pool_id = var.machine_cfg_pool_id
  api_pool_int_id     = var.api_pool_int_id
  api_pool_ext_id     = var.api_pool_ext_id
}

module "iaas" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./iaas"

  #
  # define and pass variables to:
  # data/data/powervs/bootstrap/iaas/variables.tf
  #
  cluster_id            = var.cluster_id
  resource_group        = var.powervs_resource_group
  service_instance_name = var.powervs_service_instance_name
}
