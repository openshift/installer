output "cluster-pip" {
  value = var.private ? null : azurerm_public_ip.cluster_public_ip.ip_address
}

output "public_lb_backend_pool_id" {
  value = azurerm_lb_backend_address_pool.master_public_lb_pool.id
}

output "internal_lb_backend_pool_id" {
  value = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool.id
}

output "public_lb_id" {
  value = var.private ? null : azurerm_lb.public.id
}

output "public_lb_pip_fqdn" {
  value = var.private ? null : data.azurerm_public_ip.cluster_public_ip.fqdn
}

output "internal_lb_ip_address" {
  value = azurerm_lb.internal.private_ip_address
}

output "master_nsg_name" {
  value = azurerm_network_security_group.master.name
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
