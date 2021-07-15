output "elb_backend_pool_v4_id" {
  value = ! var.azure_private ? azurestack_lb_backend_address_pool.public_lb_pool_v4[0].id : null
}

output "ilb_backend_pool_v4_id" {
  value = azurestack_lb_backend_address_pool.internal_lb_controlplane_pool_v4.id
}

output "elb_pip_v4_fqdn" {
  // TODO: Do we really need to get the fqdn from a data source instead of the resource?
  value = ! var.azure_private ? azurestack_public_ip.cluster_public_ip_v4[0].fqdn : null
}

output "elb_pip_v4" {
  value = ! var.azure_private ? azurestack_public_ip.cluster_public_ip_v4[0].ip_address : null
}

output "ilb_ip_v4_address" {
  value = azurestack_lb.internal.private_ip_addresses[0]
}

output "nsg_name" {
  value = azurestack_network_security_group.cluster.name
}

output "virtual_network_id" {
  value = local.virtual_network_id
}

output "master_subnet_id" {
  value = local.master_subnet_id
}

output "resource_group_name" {
  value = data.azurestack_resource_group.main.name
}

output "vm_image" {
  value = azurestack_image.cluster.id
}

output "storage_account" {
  value = azurestack_storage_account.cluster
}

output "availability_set_id" {
  value = azurestack_availability_set.master_availability_set.id
}
