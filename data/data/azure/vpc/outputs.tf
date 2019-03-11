output "vnet_id" {
  value = "${data.azure_vnet.cluster_vnet.id}"
}

output "cluster-pip" {
  value = "${azurerm_public_ip.cluster_public_ip.ip_address}"
}

output "public_subnet_ids" {
  value = "${local.public_subnet_ids}"
}

output "private_subnet_ids" {
  value = "${local.private_subnet_ids}"
}

output "master_nsg_id" {
  value = "${azurerm_network_security_group.master.id}"
}
