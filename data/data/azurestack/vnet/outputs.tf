output "public_lb_backend_pool_v4_id" {
  value = ! var.private ? azurestack_lb_backend_address_pool.public_lb_pool_v4[0].id : null
}

output "internal_lb_backend_pool_v4_id" {
  value = azurestack_lb_backend_address_pool.internal_lb_controlplane_pool_v4.id
}

output "public_lb_id" {
  value = ! var.private ? azurestack_lb.public[0].id : null
}

// TODO: This would be used by the CNAME record
output "public_lb_pip_v4_fqdn" {
  // TODO: Do we really need to get the fqdn from a data source instead of the resource?
  value = ! var.private ? azurestack_public_ip.cluster_public_ip_v4[0].fqdn : null
}

output "public_lb_pip_v4" {
  value = ! var.private ? azurestack_public_ip.cluster_public_ip_v4[0].ip_address : null
}

output "internal_lb_ip_v4_address" {
  value = azurestack_lb.internal.private_ip_addresses[0]
}

output "cluster_nsg_name" {
  value = azurestack_network_security_group.cluster.name
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
