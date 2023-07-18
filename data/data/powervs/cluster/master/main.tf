provider "ibm" {
  alias            = "vpc"
  ibmcloud_api_key = var.api_key
  region           = var.vpc_region
  zone             = var.vpc_zone
}

provider "ibm" {
  alias            = "powervs"
  ibmcloud_api_key = var.api_key
  region           = var.powervs_region
  zone             = var.powervs_zone
}

module "vm" {
  providers = {
    ibm = ibm.powervs
  }
  source = "./vm"

  instance_count    = var.instance_count
  memory            = var.memory
  processors        = var.processors
  cluster_id        = var.cluster_id
  proc_type         = var.proc_type
  image_id          = var.image_id
  sys_type          = var.sys_type
  cloud_instance_id = var.cloud_instance_id
  resource_group    = var.resource_group
  ignition          = var.ignition
  ssh_key_name      = var.ssh_key_name
  dhcp_id           = var.dhcp_id
  dhcp_network_id   = var.dhcp_network_id
}

module "lb" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./lb"

  instance_count      = var.instance_count
  master_ips          = module.vm.master_ips
  lb_int_id           = var.lb_int_id
  lb_ext_id           = var.lb_ext_id
  api_pool_ext_id     = var.api_pool_ext_id
  api_pool_int_id     = var.api_pool_int_id
  machine_cfg_pool_id = var.machine_cfg_pool_id
}
