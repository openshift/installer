output "vnet_id" {
  value = "${length(var.tectonic_azure_external_vnet_id) > 0 ? var.tectonic_azure_external_vnet_id : azurerm_virtual_network.tectonic_vnet.id}"
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
