output "vnet_id" {
  value = "${length(var.tectonic_azure_external_vnet_id) > 0 ? var.tectonic_azure_external_vnet_id : azurerm_virtual_network.new_vnet.id}"
}

output "cluster_default_sg" {
  value = "${aws_security_group.cluster_default.id}"
}

# We have to do this join() & split() 'trick' because null_data_source and 
# the ternary operator can't output lists or maps
#
output "master_subnet_ids" {
  value = ["${split(",", var.tectonic_azure_external_vnet_id == "" ? join(",", azure_subnet.master_subnet.*.id) :  join(",", data.azure_subnet.external_master.*.id))}"]
}

output "worker_subnet_ids" {
  value = ["${split(",", var.tectonic_azure_external_vnet_id == "" ? join(",", azure_subnet.worker_subnet.*.id) :  join(",", data.azure_subnet.external_worker.*.id))}"]
}
