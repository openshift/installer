locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
}

provider "azurerm" {
  features {}
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

provider "azureprivatedns" {
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

module "bootstrap" {
  source                 = "./bootstrap"
  resource_group_name    = azurerm_resource_group.main.name
  region                 = var.azure_region
  vm_size                = var.azure_bootstrap_vm_type
  vm_image               = azurerm_image.cluster.id
  identity               = azurerm_user_assigned_identity.main.id
  cluster_id             = var.cluster_id
  ignition               = var.ignition_bootstrap
  subnet_id              = module.vnet.master_subnet_id
  elb_backend_pool_v4_id = module.vnet.public_lb_backend_pool_v4_id
  elb_backend_pool_v6_id = module.vnet.public_lb_backend_pool_v6_id
  ilb_backend_pool_v4_id = module.vnet.internal_lb_backend_pool_v4_id
  ilb_backend_pool_v6_id = module.vnet.internal_lb_backend_pool_v6_id
  tags                   = local.tags
  storage_account        = azurerm_storage_account.cluster
  nsg_name               = module.vnet.cluster_nsg_name
  private                = module.vnet.private

  use_ipv4                  = var.use_ipv4 || var.azure_emulate_single_stack_ipv6
  use_ipv6                  = var.use_ipv6
  emulate_single_stack_ipv6 = var.azure_emulate_single_stack_ipv6
}

module "vnet" {
  source              = "./vnet"
  resource_group_name = azurerm_resource_group.main.name
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

  use_ipv4                  = var.use_ipv4 || var.azure_emulate_single_stack_ipv6
  use_ipv6                  = var.use_ipv6
  emulate_single_stack_ipv6 = var.azure_emulate_single_stack_ipv6
}

module "master" {
  source                 = "./master"
  resource_group_name    = azurerm_resource_group.main.name
  cluster_id             = var.cluster_id
  region                 = var.azure_region
  availability_zones     = var.azure_master_availability_zones
  vm_size                = var.azure_master_vm_type
  vm_image               = azurerm_image.cluster.id
  identity               = azurerm_user_assigned_identity.main.id
  ignition               = var.ignition_master
  elb_backend_pool_v4_id = module.vnet.public_lb_backend_pool_v4_id
  elb_backend_pool_v6_id = module.vnet.public_lb_backend_pool_v6_id
  ilb_backend_pool_v4_id = module.vnet.internal_lb_backend_pool_v4_id
  ilb_backend_pool_v6_id = module.vnet.internal_lb_backend_pool_v6_id
  subnet_id              = module.vnet.master_subnet_id
  instance_count         = var.master_count
  storage_account        = azurerm_storage_account.cluster
  os_volume_type         = var.azure_master_root_volume_type
  os_volume_size         = var.azure_master_root_volume_size
  private                = module.vnet.private

  use_ipv4                  = var.use_ipv4 || var.azure_emulate_single_stack_ipv6
  use_ipv6                  = var.use_ipv6
  emulate_single_stack_ipv6 = var.azure_emulate_single_stack_ipv6
}

module "dns" {
  source                          = "./dns"
  cluster_domain                  = var.cluster_domain
  cluster_id                      = var.cluster_id
  base_domain                     = var.base_domain
  virtual_network_id              = module.vnet.virtual_network_id
  external_lb_fqdn_v4             = module.vnet.public_lb_pip_v4_fqdn
  external_lb_fqdn_v6             = module.vnet.public_lb_pip_v6_fqdn
  internal_lb_ipaddress_v4        = module.vnet.internal_lb_ip_v4_address
  internal_lb_ipaddress_v6        = module.vnet.internal_lb_ip_v6_address
  resource_group_name             = azurerm_resource_group.main.name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private                         = module.vnet.private

  use_ipv4                  = var.use_ipv4 || var.azure_emulate_single_stack_ipv6
  use_ipv6                  = var.use_ipv6
  emulate_single_stack_ipv6 = var.azure_emulate_single_stack_ipv6
}

resource "random_string" "storage_suffix" {
  length  = 5
  upper   = false
  special = false
}

resource "azurerm_resource_group" "main" {
  name     = "${var.cluster_id}-rg"
  location = var.azure_region
  tags     = local.tags
}

data "azurerm_resource_group" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  name = var.azure_network_resource_group_name
}

resource "azurerm_storage_account" "cluster" {
  name                     = "cluster${random_string.storage_suffix.result}"
  resource_group_name      = azurerm_resource_group.main.name
  location                 = var.azure_region
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_user_assigned_identity" "main" {
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location

  name = "${var.cluster_id}-identity"
}

resource "azurerm_role_assignment" "main" {
  scope                = azurerm_resource_group.main.id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}

resource "azurerm_role_assignment" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  scope                = data.azurerm_resource_group.network[0].id
  role_definition_name = "Contributor"
  principal_id         = azurerm_user_assigned_identity.main.principal_id
}

# copy over the vhd to cluster resource group and create an image using that
resource "azurerm_storage_container" "vhd" {
  name                 = "vhd"
  storage_account_name = azurerm_storage_account.cluster.name
}

resource "azurerm_storage_blob" "rhcos_image" {
  name                   = "rhcos${random_string.storage_suffix.result}.vhd"
  storage_account_name   = azurerm_storage_account.cluster.name
  storage_container_name = azurerm_storage_container.vhd.name
  type                   = "Block"
  source_uri             = var.azure_image_url
  metadata               = map("source_uri", var.azure_image_url)
}

resource "azurerm_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = azurerm_resource_group.main.name
  location            = var.azure_region

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = azurerm_storage_blob.rhcos_image.url
  }
}
