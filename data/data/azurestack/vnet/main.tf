locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )
  description = "Created By OpenShift Installer"
  upload = var.azure_image_path != ""
  rhcos_blob_uri = local.upload ? azurestackprivate_vhd_blob.rhcos_image[0].url : azurestack_storage_blob.rhcos_image[0].url
}

provider "azurestack" {
  arm_endpoint    = var.azure_arm_endpoint
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

provider "azurestackprivate" {
  arm_endpoint    = var.azure_arm_endpoint
  subscription_id = var.azure_subscription_id
  client_id       = var.azure_client_id
  client_secret   = var.azure_client_secret
  tenant_id       = var.azure_tenant_id
}

resource "random_string" "storage_suffix" {
  length  = 5
  upper   = false
  special = false
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
  name                     = "cluster${random_string.storage_suffix.result}"
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
  count = local.upload ? 0 : 1

  name                   = "rhcos${random_string.storage_suffix.result}.vhd"
  resource_group_name    = data.azurestack_resource_group.main.name
  storage_account_name   = azurestack_storage_account.cluster.name
  storage_container_name = azurestack_storage_container.vhd.name
  type                   = "page"
  source_uri             = var.azure_image_url
}

resource "azurestackprivate_vhd_blob" "rhcos_image" {
  count = local.upload ? 1 : 0

  name                   = "rhcos${random_string.storage_suffix.result}.vhd"
  resource_group_name    = data.azurestack_resource_group.main.name
  storage_account_name   = azurestack_storage_account.cluster.name
  storage_container_name = azurestack_storage_container.vhd.name
  source                 = var.azure_image_path
}

resource "azurestack_image" "cluster" {
  name                = var.cluster_id
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region

  os_disk {
    os_type  = "Linux"
    os_state = "Generalized"
    blob_uri = local.rhcos_blob_uri
  }
}

resource "azurestack_availability_set" "master_availability_set" {
  name                = "${var.cluster_id}-master"
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region
  managed             = true
}
