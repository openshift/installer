output "public_lb_backend_pool_v4_id" {
  value = local.need_public_ipv4 ? azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id : null
}

output "public_lb_backend_pool_v6_id" {
  value = local.need_public_ipv6 ? azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id : null
}

output "internal_lb_backend_pool_v4_id" {
  value = var.use_ipv4 ? azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v4[0].id : null
}

output "internal_lb_backend_pool_v6_id" {
  value = var.use_ipv6 ? azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[0].id : null
}

output "public_lb_id" {
  value = var.private ? null : azurerm_lb.public.id
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

output "cluster_nsg_name" {
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

output "private" {
  value = var.private
}
