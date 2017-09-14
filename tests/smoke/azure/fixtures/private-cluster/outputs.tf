/***
    Outputs
***/

output "vpn-gw-endpoint" {
  value = "${azurerm_public_ip.vpn_gw_ip.fqdn}"
}

output "tectonic_azure_external_resource_group" {
  value = "${azurerm_resource_group.tectonic_cluster.id}"
}

output "tectonic_azure_external_vnet_id" {
  value = "${azurerm_virtual_network.tectonic_private_vnet.id}"
}
