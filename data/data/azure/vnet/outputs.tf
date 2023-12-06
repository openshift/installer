output "elb_backend_pool_v4_id" {
  value = local.need_public_ipv4 ? azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id : null
}

output "elb_backend_pool_v6_id" {
  value = local.need_public_ipv6 ? azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id : null
}

output "ilb_backend_pool_v4_id" {
  value = var.use_ipv4 ? azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v4[0].id : null
}

output "ilb_backend_pool_v6_id" {
  value = var.use_ipv6 ? azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[0].id : null
}

output "public_lb_pip_v4_fqdn" {
  value = local.need_public_ipv4 ? data.azurerm_public_ip.cluster_public_ip_v4[0].fqdn : null
}

output "public_lb_pip_v6_fqdn" {
  value = local.need_public_ipv6 ? data.azurerm_public_ip.cluster_public_ip_v6[0].fqdn : null
}

output "internal_lb_ip_v4_address" {
  value = var.use_ipv4 ? azurerm_lb.internal.private_ip_addresses[0] : null
}

output "internal_lb_ip_v6_address" {
  value = var.use_ipv6 ? azurerm_lb.internal.private_ip_addresses[1] : null
}

output "nsg_name" {
  value = azurerm_network_security_group.cluster.name
}

output "virtual_network_id" {
  value = local.virtual_network_id
}

output "master_subnet_id" {
  value = local.master_subnet_id
}

output "worker_subnet_id" {
  value = local.worker_subnet_id
}

output "resource_group_name" {
  value = data.azurerm_resource_group.main.name
}

output "identity" {
  value = azurerm_user_assigned_identity.main.id
}

output "key_vault_key_id" {
  value = var.azure_keyvault_name != "" ? data.azurerm_key_vault.keyvault[0].id : null
}

output "user_assigned_identity_id" {
  value = var.azure_keyvault_name != "" ? data.azurerm_user_assigned_identity.keyvault_identity[0].id : null
}

output "subnet_id" {
  value = local.master_subnet_id
}

output "image_version_gallery_name" {
  value = azurerm_shared_image.cluster.gallery_name
}

output "image_version_gen2_gallery_name" {
  value = azurerm_shared_image.clustergen2.gallery_name
}

output "image_version_name" {
  value = azurerm_shared_image.cluster.name
}

output "image_version_gen2_name" {
  value = azurerm_shared_image.clustergen2.name
}
