#######################################
# Master module outputs
#######################################

output "name_list" {
  value = ibm_is_instance.master_node.*.name
}

output "primary_ipv4_address_list" {
  value = ibm_is_instance.master_node.*.primary_network_interface.0.primary_ipv4_address
}