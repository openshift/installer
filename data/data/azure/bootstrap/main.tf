locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
  description = "Created By OpenShift Installer"
  # At this time min_tls_version is only supported in the Public Cloud and US Government Cloud.
  environments_with_min_tls_version = ["public", "usgovernment"]

}

provider "azurerm" {
  features {}
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
  environment     = var.azure_environment
}

provider "azureprivatedns" {
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
  environment     = var.azure_environment
}

module "bootstrap" {
  source                 = "./bootstrap"
  resource_group_name    = data.azurerm_resource_group.main.name
  region                 = var.azure_region
  vm_size                = var.azure_bootstrap_vm_type
  vm_image               = data.azurerm_image.cluster.id
  identity               = data.azurerm_user_assigned_identity.main.id
  cluster_id             = var.cluster_id
  ignition               = var.ignition_bootstrap
  subnet_id              = module.vnet.master_subnet_id
  elb_backend_pool_v4_id = module.vnet.public_lb_backend_pool_v4_id
  elb_backend_pool_v6_id = module.vnet.public_lb_backend_pool_v6_id
  ilb_backend_pool_v4_id = module.vnet.internal_lb_backend_pool_v4_id
  ilb_backend_pool_v6_id = module.vnet.internal_lb_backend_pool_v6_id
  tags                   = local.tags
  storage_account        = data.azurerm_storage_account.cluster
  nsg_name               = module.vnet.cluster_nsg_name
  private                = module.vnet.private
  outbound_udr           = var.azure_outbound_user_defined_routing

  use_ipv4 = var.use_ipv4
  use_ipv6 = var.use_ipv6
}

module "vnet" {
  source              = "./vnet"
  resource_group_name = data.azurerm_resource_group.main.name
  vnet_v4_cidrs       = var.machine_v4_cidrs
  vnet_v6_cidrs       = var.machine_v6_cidrs
  cluster_id          = var.cluster_id
  region              = var.azure_region
  dns_label           = var.cluster_id

  preexisting_network         = var.azure_preexisting_network
  network_resource_group_name = var.azure_network_resource_group_name
  virtual_network_name        = var.azure_virtual_network
  master_subnet               = var.azure_control_plane_subnet
  worker_subnet               = var.azure_compute_subnet
  private                     = var.azure_private
  outbound_udr                = var.azure_outbound_user_defined_routing

  use_ipv4 = var.use_ipv4
  use_ipv6 = var.use_ipv6
}

data "azurerm_resource_group" "main" {
  name = var.azure_resource_group_name == "" ? "${var.cluster_id}-rg" : var.azure_resource_group_name
}

data "azurerm_storage_account" "cluster" {
  name                = "cluster${var.azure_storage_suffix}"
  resource_group_name = data.azurerm_resource_group.main.name
}

data "azurerm_user_assigned_identity" "main" {
  name                = "${var.cluster_id}-identity"
  resource_group_name = data.azurerm_resource_group.main.name
}

data "azurerm_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = data.azurerm_resource_group.main.name
}
