output "vnet_id" {
  value = "${var.external_vnet_name == "" ? join("",azurerm_virtual_network.tectonic_vnet.*.name) : var.external_vnet_name }"
}

# We have to do this join() & split() 'trick' because null_data_source and 
# the ternary operator can't output lists or maps
#
output "master_subnet" {
  value = "${azurerm_subnet.master_subnet.id}"
}

output "worker_subnet" {
  value = "${azurerm_subnet.worker_subnet.id}"
}
