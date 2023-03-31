output "resource_group_name" {
  value = data.azurerm_resource_group.main.name
}

output "storage_account_name" {
  value = azurerm_storage_account.cluster.name
}

output "rhcos_image_url" {
  value = azurerm_storage_blob.rhcos_image.url
}