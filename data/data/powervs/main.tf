provider "ibm" {
  ibmcloud_api_key      = var.ibmcloud_api_key
  region                = var.ibmcloud_region
  zone                  = var.ibmcloud_zone
  iaas_classic_username = "apikey"
  iaas_classic_api_key  = var.ibmcloud_api_key
}

module "bootstrap" {
  source            = "./bootstrap"
  cloud_instance_id = var.powervs_cloud_instance_id
  cluster_id        = var.powervs_cluster_id
  resource_group    = var.powervs_resource_group

  cos_instance_location = var.powervs_cos_instance_location
  cos_bucket_location   = var.powervs_cos_bucket_location
  cos_storage_class     = var.powervs_cos_storage_class

  memory     = var.powervs_bootstrap_memory
  processors = var.powervs_bootstrap_processors
  ignition   = var.powervs_bootstrap_ignition
  sys_type   = var.powervs_sys_type
  proc_type  = var.powervs_proc_type

  # TODO(mjturek): image and network IDs are not derived during terraform
  #                for other providers. Need to investigate and follow how
  #                other providers do this. cnorman's branch has some work
  #                towards this already.
  image_name   = var.powervs_image_name
  network_name = var.powervs_network_name
}

data "ibm_is_subnet" "vpc_subnet" {
  name = var.powervs_vpc_subnet_name
}

module "loadbalancer" {
  source = "./loadbalancer"

  cluster_id    = var.powervs_cluster_id
  vpc_name      = var.powervs_vpc_name
  vpc_subnet_id = data.ibm_is_subnet.vpc_subnet.id
  bootstrap_ip  = module.bootstrap.bootstrap_ip

  # TODO add resources for master/controller
  master_ips = []
}


module "dns" {
  source = "./dns"

  base_domain                = var.powervs_base_domain
  cluster_id                 = var.powervs_cluster_id
  cluster_domain             = var.powervs_cluster_domain
  load_balancer_hostname     = module.loadbalancer.powervs_lb_hostname
  load_balancer_int_hostname = module.loadbalancer.powervs_lb_int_hostname
}
