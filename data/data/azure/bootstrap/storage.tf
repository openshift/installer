locals {
  tags = merge(
    {
      "kubernetes.io_cluster.${var.cluster_id}" = "owned"
    },
    var.azure_extra_tags,
  )

  # At this time min_tls_version is only supported in the Public Cloud and US Government Cloud.
  environments_with_min_tls_version = ["public", "usgovernment"]
}

resource "azurerm_storage_account" "cluster" {
  name                            = "cluster${var.random_storage_account_suffix}"
  resource_group_name             = var.resource_group_name
  location                        = var.azure_region
  account_tier                    = "Standard"
  account_replication_type        = "LRS"
  min_tls_version                 = contains(local.environments_with_min_tls_version, var.azure_environment) ? "TLS1_2" : null
  allow_nested_items_to_be_public = false
  tags                            = var.azure_extra_tags
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

resource "azurerm_shared_image" "bootstrap_gen2" {
  name                = "${var.cluster_id}-bootstrap-gen2"
  gallery_name        = var.image_version_gallery_name
  resource_group_name = var.resource_group_name
  location            = var.azure_region
  os_type             = "Linux"
  hyper_v_generation  = "V2"
  architecture        = var.azure_vm_architecture

  identifier {
    publisher = "RedHat-gen2"
    offer     = "rhcos-gen2"
    sku       = "bootstrap"
  }

  tags = var.azure_extra_tags
}

resource "azurerm_shared_image_version" "bootstrap_image_version" {
  name                = var.azure_image_release
  gallery_name        = azurerm_shared_image.bootstrap_gen2.gallery_name
  image_name          = azurerm_shared_image.bootstrap_gen2.name
  resource_group_name = var.resource_group_name
  location            = var.azure_region

  blob_uri           = azurerm_storage_blob.rhcos_image.url
  storage_account_id = azurerm_storage_account.cluster.id

  target_region {
    name                   = var.azure_region
    regional_replica_count = 1
  }

  tags = var.azure_extra_tags
}
