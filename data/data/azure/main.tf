locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )

  master_subnet_cidr = cidrsubnet(var.machine_cidr, 3, 0) #master subnet is a smaller subnet within the vnet. i.e from /21 to /24
  node_subnet_cidr   = cidrsubnet(var.machine_cidr, 3, 1) #node subnet is a smaller subnet within the vnet. i.e from /21 to /24
}

provider "azurerm" {
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

module "bootstrap" {
  source                  = "./bootstrap"
  resource_group_name     = azurerm_resource_group.main.name
  region                  = var.azure_region
  vm_size                 = var.azure_bootstrap_vm_type
  vm_image                = var.azure_image_id
  identity                = azurerm_user_assigned_identity.main.id
  cluster_id              = var.cluster_id
  ignition                = var.ignition_bootstrap
  subnet_id               = module.vnet.public_subnet_id
  elb_backend_pool_id     = module.vnet.public_lb_backend_pool_id
  ilb_backend_pool_id     = module.vnet.internal_lb_backend_pool_id
  tags                    = local.tags
  boot_diag_blob_endpoint = azurerm_storage_account.bootdiag.primary_blob_endpoint
  ssh_nat_rule_id         = module.vnet.bootstrap_ssh_nat_rule_id

  # This is to create explicit dependency on private zone to exist before VMs are created in the vnet. https://github.com/MicrosoftDocs/azure-docs/issues/13728
  private_dns_zone_id = azurerm_dns_zone.private.id
}

module "vnet" {
  source              = "./vnet"
  resource_group_name = azurerm_resource_group.main.name
  vnet_cidr           = var.machine_cidr
  master_subnet_cidr  = local.master_subnet_cidr
  node_subnet_cidr    = local.node_subnet_cidr
  cluster_id          = var.cluster_id
  region              = var.azure_region
  dns_label           = var.cluster_id
  master_count        = var.master_count
}

module "master" {
  source                  = "./master"
  resource_group_name     = azurerm_resource_group.main.name
  cluster_id              = var.cluster_id
  region                  = var.azure_region
  vm_size                 = var.azure_master_vm_type
  vm_image                = var.azure_image_id
  identity                = azurerm_user_assigned_identity.main.id
  ignition                = var.ignition_master
  external_lb_id          = module.vnet.public_lb_id
  elb_backend_pool_id     = module.vnet.public_lb_backend_pool_id
  ilb_backend_pool_id     = module.vnet.internal_lb_backend_pool_id
  subnet_id               = module.vnet.public_subnet_id
  master_subnet_cidr      = local.master_subnet_cidr
  instance_count          = var.master_count
  boot_diag_blob_endpoint = azurerm_storage_account.bootdiag.primary_blob_endpoint
  os_volume_size          = var.azure_master_root_volume_size
  ssh_nat_rule_ids        = module.vnet.mmaster_ssh_nat_rule_ids

  # This is to create explicit dependency on private zone to exist before VMs are created in the vnet. https://github.com/MicrosoftDocs/azure-docs/issues/13728
  private_dns_zone_id = azurerm_dns_zone.private.id
}

module "dns" {
  source                          = "./dns"
  cluster_domain                  = var.cluster_domain
  base_domain                     = var.base_domain
  external_lb_fqdn                = module.vnet.public_lb_pip_fqdn
  internal_lb_ipaddress           = module.vnet.internal_lb_ip_address
  resource_group_name             = azurerm_resource_group.main.name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private_dns_zone_name           = azurerm_dns_zone.private.name
  etcd_count                      = var.master_count
  etcd_ip_addresses               = module.master.ip_addresses
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

resource "azurerm_storage_account" "bootdiag" {
  name                     = "bootdiagmasters${random_string.storage_suffix.result}"
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

# https://github.com/MicrosoftDocs/azure-docs/issues/13728
resource "azurerm_dns_zone" "private" {
  name                           = var.cluster_domain
  resource_group_name            = azurerm_resource_group.main.name
  zone_type                      = "Private"
  resolution_virtual_network_ids = [module.vnet.vnet_id]
}
