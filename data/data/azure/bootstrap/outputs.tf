output "bootstrap_ip" {
  value = var.azure_private ? azurerm_network_interface.bootstrap.private_ip_address : azurerm_public_ip.bootstrap_public_ip_v4[0].ip_address
}

output "storage_account_id" {
  value = azurerm_storage_account.cluster.id
}

output "storage_rhcos_image_url" {
  value = azurerm_storage_blob.rhcos_image.url
}