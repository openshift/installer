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
  resource_group_name             = data.azurerm_resource_group.main.name
  base_domain_resource_group_name = var.azure_base_domain_resource_group_name
  private                         = module.vnet.private

  use_ipv4 = var.use_ipv4
  use_ipv6 = var.use_ipv6
}

resource "azurerm_resource_group" "main" {
  count = var.azure_resource_group_name == "" ? 1 : 0

  name     = "${var.cluster_id}-rg"
  location = var.azure_region
  tags     = local.tags
}

data "azurerm_resource_group" "main" {
  name = var.azure_resource_group_name == "" ? "${var.cluster_id}-rg" : var.azure_resource_group_name

  depends_on = [azurerm_resource_group.main]
}

data "azurerm_resource_group" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  name = var.azure_network_resource_group_name
}

resource "azurerm_storage_account" "cluster" {
  name                     = "cluster${var.azure_storage_suffix}"
  resource_group_name      = data.azurerm_resource_group.main.name
  location                 = var.azure_region
  account_tier             = "Standard"
  account_replication_type = "LRS"
  min_tls_version          = contains(local.environments_with_min_tls_version, var.azure_environment) ? "TLS1_2" : null
}

resource "azurerm_user_assigned_identity" "main" {
  resource_group_name = data.azurerm_resource_group.main.name
  location            = data.azurerm_resource_group.main.location

  name = "${var.cluster_id}-identity"
}

resource "azurerm_role_assignment" "main" {
  scope                = data.azurerm_resource_group.main.id
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
  name                   = "rhcos${var.azure_storage_suffix}.vhd"
  storage_account_name   = azurerm_storage_account.cluster.name
  storage_container_name = azurerm_storage_container.vhd.name
  type                   = "Page"
  source_uri             = var.azure_image_url
  metadata               = map("source_uri", var.azure_image_url)
}

resource "azurerm_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = azurerm_storage_blob.rhcos_image.url
  }
}
