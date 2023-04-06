locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
  description = "Created By OpenShift Installer"

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
  tags     = var.azure_extra_tags
}

data "azurerm_resource_group" "main" {
  name = var.azure_resource_group_name == "" ? "${var.cluster_id}-rg" : var.azure_resource_group_name

  depends_on = [azurerm_resource_group.main]
}

data "azurerm_resource_group" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  name = var.azure_network_resource_group_name
}

resource "azurerm_user_assigned_identity" "main" {
  resource_group_name = data.azurerm_resource_group.main.name
  location            = data.azurerm_resource_group.main.location
  name                = "${var.cluster_id}-identity"
  tags                = var.azure_extra_tags
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


# Creates Shared Image Gallery
# https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/shared_image_gallery
resource "azurerm_shared_image_gallery" "sig" {
  name                = "gallery_${replace(var.cluster_id, "-", "_")}"
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region
  tags                = var.azure_extra_tags
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

  tags = var.azure_extra_tags
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

  tags = var.azure_extra_tags
}
