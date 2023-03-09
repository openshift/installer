output "bootstrap_ip" {
  value = var.azure_private ? azurerm_network_interface.bootstrap.private_ip_address : azurerm_public_ip.bootstrap_public_ip_v4[0].ip_address
}

output "vm_image" {
  value = var.azure_hypervgeneration_version == "V2" ? azurerm_shared_image_version.clustergen2_image_version.id : azurerm_shared_image_version.cluster_image_version.id
}

output "storage_account_name" {
  value = azurerm_storage_account.cluster.name
}
