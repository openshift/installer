output "cluster-pip" {
  value = azurerm_public_ip.cluster_public_ip.ip_address
}

output "public_subnet_id" {
  value = local.subnet_ids
}

output "public_lb_backend_pool_id" {
  value = azurerm_lb_backend_address_pool.master_public_lb_pool.id
}

output "internal_lb_backend_pool_id" {
  value = local.internal_lb_controlplane_pool_id
}

output "public_lb_id" {
  value = local.public_lb_id
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

output "bootstrap_ssh_nat_rule_id" {
  value = azurerm_lb_nat_rule.bootstrap_ssh.id
}

output "mmaster_ssh_nat_rule_ids" {
  value = azurerm_lb_nat_rule.master_ssh.*.id
}

