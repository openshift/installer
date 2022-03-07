locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
  description = "Created By OpenShift Installer"
}

provider "azurestack" {
  arm_endpoint    = var.azure_arm_endpoint
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}


module "master" {
  source                 = "./master"
  resource_group_name    = var.resource_group_name
  cluster_id             = var.cluster_id
  region                 = var.azure_region
  vm_size                = var.azure_master_vm_type
  vm_image_uri           = var.vm_image
  ignition               = var.ignition_master
  elb_backend_pool_v4_id = var.elb_backend_pool_v4_id
  ilb_backend_pool_v4_id = var.ilb_backend_pool_v4_id
  subnet_id              = var.master_subnet_id
  instance_count         = var.master_count
  storage_account        = var.storage_account
  os_volume_type         = var.azure_master_root_volume_type
  os_volume_size         = var.azure_master_root_volume_size
  private                = var.azure_private
  availability_set_id    = var.availability_set_id
}

module "dns" {
  source                          = "./dns"
  cluster_domain                  = var.cluster_domain
  cluster_id                      = var.cluster_id
  base_domain                     = var.base_domain
  virtual_network_id              = var.virtual_network_id
  elb_fqdn_v4                     = var.elb_pip_v4_fqdn
  elb_pip_v4                      = var.elb_pip_v4
  ilb_ipaddress_v4                = var.ilb_ip_v4_address
  resource_group_name             = var.resource_group_name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private                         = var.azure_private
  tags                            = local.tags
}
