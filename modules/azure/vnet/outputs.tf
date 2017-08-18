output "vnet_id" {
  value = "${var.external_vnet_id == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : var.external_vnet_id }"
}

output "master_subnet" {
  value = "${var.external_vnet_id == "" ?  join(" ", azurerm_subnet.master_subnet.*.id) : var.external_master_subnet_id }"
}

output "worker_subnet" {
  value = "${var.external_vnet_id == "" ?  join(" ", azurerm_subnet.worker_subnet.*.id) : var.external_worker_subnet_id }"
}

output "worker_subnet_name" {
  value = "${var.external_vnet_id == "" ?  join(" ", azurerm_subnet.worker_subnet.*.name) : replace(var.external_vnet_id, "${var.const_id_to_group_name_regex}", "$2") }"
}

# TODO: Allow user to provide their own network
output "etcd_cidr" {
  value = "${azurerm_subnet.master_subnet.address_prefix}"
}

output "master_cidr" {
  value = "${azurerm_subnet.master_subnet.address_prefix}"
}

output "worker_cidr" {
  value = "${azurerm_subnet.worker_subnet.address_prefix}"
}

output "etcd_nsg_name" {
  value = "${var.external_nsg_etcd_id == "" ? join(" ", azurerm_network_security_group.etcd.*.name) : replace(var.external_nsg_etcd_id, "${var.const_id_to_group_name_regex}", "$2")}"
}

# TODO: Allow user to provide their own network
output "worker_nsg_name" {
  value = "${var.external_nsg_worker_id == "" ? join(" ", azurerm_network_security_group.worker.*.name) : var.external_nsg_worker_id }"
}

output "etcd_network_interface_ids" {
  value = ["${azurerm_network_interface.etcd_nic.*.id}"]
}

output "etcd_endpoints" {
  value = "${azurerm_network_interface.etcd_nic.*.private_ip_address}"
}

output "master_network_interface_ids" {
  value = ["${azurerm_network_interface.tectonic_master.*.id}"]
}

output "worker_network_interface_ids" {
  value = ["${azurerm_network_interface.tectonic_worker.*.id}"]
}

output "master_private_ip_addresses" {
  value = ["${azurerm_network_interface.tectonic_master.*.private_ip_address}"]
}

output "worker_private_ip_addresses" {
  value = ["${azurerm_network_interface.tectonic_worker.*.private_ip_address}"]
}

output "api_ip_addresses" {
  value = ["${azurerm_public_ip.api_ip.ip_address}"]
}

output "console_ip_addresses" {
  value = ["${azurerm_public_ip.console_ip.ip_address}"]
}

output "ingress_fqdn" {
  value = "${var.base_domain == "" ? azurerm_public_ip.console_ip.fqdn : "${var.cluster_name}.${var.base_domain}"}"
}

output "api_fqdn" {
  value = "${var.base_domain == "" ? azurerm_public_ip.api_ip.fqdn : "${var.cluster_name}-api.${var.base_domain}"}"
}
