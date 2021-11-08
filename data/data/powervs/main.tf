provider "ibm" {
  alias = "vpc"
  ibmcloud_api_key      = var.powervs_api_key
  region                = var.powervs_vpc_region
}

provider "ibm" {
  alias = "powervs"
  ibmcloud_api_key      = var.powervs_api_key
  region                = var.powervs_region
}

resource "ibm_pi_key" "cluster_key" {
  provider = ibm.powervs
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

  cos_instance_location = var.powervs_cos_instance_location
  cos_bucket_location   = var.powervs_cos_bucket_location
  cos_storage_class     = var.powervs_cos_storage_class

  memory       = var.powervs_bootstrap_memory
  processors   = var.powervs_bootstrap_processors
  ignition     = var.ignition_bootstrap
  sys_type     = var.powervs_sys_type
  proc_type    = var.powervs_proc_type
  key_id       = ibm_pi_key.cluster_key.key_id
  image_name   = var.powervs_image_name
  network_name = var.powervs_network_name
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

  memory       = var.powervs_master_memory
  processors   = var.powervs_master_processors
  ignition     = var.ignition_master
  sys_type     = var.powervs_sys_type
  proc_type    = var.powervs_proc_type
  key_id       = ibm_pi_key.cluster_key.key_id
  image_name   = var.powervs_image_name
  network_name = var.powervs_network_name
}

data "ibm_is_subnet" "vpc_subnet" {
  provider = ibm.vpc
  name = var.powervs_vpc_subnet_name
}

data "ibm_pi_image" "boot_image" {
  provider = ibm.powervs
  pi_image_name    = var.powervs_image_name
  pi_cloud_instance_id = var.powervs_cloud_instance_id
}

data "ibm_pi_network" "pvs_net" {
  provider = ibm.powervs
  pi_network_name = var.powervs_network_name
  pi_cloud_instance_id = var.powervs_cloud_instance_id
}

module "loadbalancer" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./loadbalancer"

  cluster_id    = var.cluster_id
  vpc_name      = var.powervs_vpc_name
  vpc_subnet_id = data.ibm_is_subnet.vpc_subnet.id
  bootstrap_ip  = module.bootstrap.bootstrap_ip
  master_ips = module.master.master_ips
  resource_group    = var.powervs_resource_group
}


module "dns" {
  providers = {
    ibm = ibm.vpc
  }
  source = "./dns"

  cis_id                     = var.powervs_cis_crn
  base_domain                = var.base_domain
  cluster_domain             = var.cluster_domain
  load_balancer_hostname     = module.loadbalancer.powervs_lb_hostname
  load_balancer_int_hostname = module.loadbalancer.powervs_lb_int_hostname
}
