output "dns_server" {
  value = var.publish_strategy == "Internal" ? ibm_is_instance.dns_vm_vsi[0].primary_network_interface[0].primary_ipv4_address : "1.1.1.1"
}
