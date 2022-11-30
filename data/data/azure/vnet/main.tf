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
  subscription_id             = var.azure_subscription_id
  client_id                   = var.azure_client_id
  client_secret               = var.azure_client_secret
  client_certificate_password = var.azure_certificate_password
  client_certificate_path     = var.azure_certificate_path
  tenant_id                   = var.azure_tenant_id
  environment                 = var.azure_environment
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
  name                            = "cluster${var.random_storage_account_suffix}"
  resource_group_name             = data.azurerm_resource_group.main.name
  location                        = var.azure_region
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  min_tls_version                 = contains(local.environments_with_min_tls_version, var.azure_environment) ? "TLS1_2" : null
  allow_nested_items_to_be_public = false
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
  name                   = "rhcos${var.random_storage_account_suffix}.vhd"
  storage_account_name   = azurerm_storage_account.cluster.name
  storage_container_name = azurerm_storage_container.vhd.name
  type                   = "Page"
  source_uri             = var.azure_image_url
  metadata               = tomap({ source_uri = var.azure_image_url })
}

# Creates Shared Image Gallery
# https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/shared_image_gallery
resource "azurerm_shared_image_gallery" "sig" {
  name                = "gallery_${replace(var.cluster_id, "-", "_")}"
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region
}

# Creates image definition
# https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/shared_image
resource "azurerm_shared_image" "cluster" {
  name                = var.cluster_id
  gallery_name        = azurerm_shared_image_gallery.sig.name
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region
  os_type             = "Linux"
  architecture        = var.azure_vm_architecture

  identifier {
    publisher = "RedHat"
    offer     = "rhcos"
    sku       = "basic"
  }
}

resource "azurerm_shared_image" "clustergen2" {
  name                = "${var.cluster_id}-gen2"
  gallery_name        = azurerm_shared_image_gallery.sig.name
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region
  os_type             = "Linux"
  hyper_v_generation  = "V2"
  architecture        = var.azure_vm_architecture

  identifier {
    publisher = "RedHat-gen2"
    offer     = "rhcos-gen2"
    sku       = "gen2"
  }
}

resource "azurerm_shared_image_version" "cluster_image_version" {
  name                = var.azure_image_release
  gallery_name        = azurerm_shared_image.cluster.gallery_name
  image_name          = azurerm_shared_image.cluster.name
  resource_group_name = azurerm_shared_image.cluster.resource_group_name
  location            = azurerm_shared_image.cluster.location

  blob_uri           = azurerm_storage_blob.rhcos_image.url
  storage_account_id = azurerm_storage_account.cluster.id

  target_region {
    name                   = azurerm_shared_image.cluster.location
    regional_replica_count = 1
  }
}

resource "azurerm_shared_image_version" "clustergen2_image_version" {
  name                = var.azure_image_release
  gallery_name        = azurerm_shared_image.clustergen2.gallery_name
  image_name          = azurerm_shared_image.clustergen2.name
  resource_group_name = azurerm_shared_image.clustergen2.resource_group_name
  location            = azurerm_shared_image.clustergen2.location

  blob_uri           = azurerm_storage_blob.rhcos_image.url
  storage_account_id = azurerm_storage_account.cluster.id

  target_region {
    name                   = azurerm_shared_image.clustergen2.location
    regional_replica_count = 1
  }
}

