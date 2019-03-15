output "vnet_id" {
  value = "${local.vnet_id}"
}

output "cluster-pip" {
  value = "${azurerm_public_ip.cluster_public_ip.ip_address}"
}

output "public_subnet_id" {
  value = "${local.subnet_ids}"
}

output "lb_backend_pool_id" {
  value="${local.lb_backend_pool_id}"
}

output "bootstrap_ssh_natrule_id" {
  value = "${local.bootstrap_ssh_natrule_id}"
}

output "master_nsg_id" {
  value = "${azurerm_network_security_group.master.id}"
}
