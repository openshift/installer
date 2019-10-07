output "cluster-pip" {
  value = azurerm_public_ip.cluster_public_ip.ip_address
}

output "network_id" {
  value = data.azurerm_virtual_network.cluster_vnet.id
}

output "public_subnet_id" {
  value = azurerm_subnet.master_subnet.id
}

output "public_lb_backend_pool_id" {
  value = azurerm_lb_backend_address_pool.master_public_lb_pool.id
}

output "internal_lb_backend_pool_id" {
  value = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool.id
}

output "public_lb_id" {
  value = azurerm_lb.public.id
}

output "public_lb_pip_fqdn" {
  value = data.azurerm_public_ip.cluster_public_ip.fqdn
}

output "internal_lb_ip_address" {
  value = azurerm_lb.internal.private_ip_address
}

output "master_nsg_name" {
  value = azurerm_network_security_group.master.name
}
