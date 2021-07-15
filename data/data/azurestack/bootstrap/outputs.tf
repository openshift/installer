output "bootstrap_ip" {
  value = var.azure_private ? azurestack_network_interface.bootstrap.private_ip_address : azurestack_public_ip.bootstrap_public_ip_v4[0].ip_address
}
