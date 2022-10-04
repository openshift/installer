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

resource "azurestack_resource_group" "main" {
  count = var.azure_resource_group_name == "" ? 1 : 0

  name     = "${var.cluster_id}-rg"
  location = var.azure_region
  tags     = local.tags
}

data "azurestack_resource_group" "main" {
  name = var.azure_resource_group_name == "" ? "${var.cluster_id}-rg" : var.azure_resource_group_name

  depends_on = [azurestack_resource_group.main]
}

data "azurestack_resource_group" "network" {
  count = var.azure_preexisting_network ? 1 : 0

  name = var.azure_network_resource_group_name
}

resource "azurestack_storage_account" "cluster" {
  name                     = "cluster${var.random_storage_account_suffix}"
  resource_group_name      = data.azurestack_resource_group.main.name
  location                 = var.azure_region
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

# copy over the vhd to cluster resource group and create an image using that
resource "azurestack_storage_container" "vhd" {
  name                 = "vhd"
  resource_group_name  = data.azurestack_resource_group.main.name
  storage_account_name = azurestack_storage_account.cluster.name
}

resource "azurestack_storage_blob" "rhcos_image" {
  name                   = "rhcos${var.random_storage_account_suffix}.vhd"
  resource_group_name    = data.azurestack_resource_group.main.name
  storage_account_name   = azurestack_storage_account.cluster.name
  storage_container_name = azurestack_storage_container.vhd.name
  type                   = "page"
  source_uri             = var.azure_image_url
}

resource "azurestack_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = azurestack_storage_blob.rhcos_image.url
  }
}

resource "azurestack_availability_set" "cluster_availability_set" {
  name                = "${var.cluster_id}-cluster"
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region
  managed             = true
}
