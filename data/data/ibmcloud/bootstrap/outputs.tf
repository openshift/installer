output "bootstrap_ip" {
  value = local.public_endpoints ? ibm_is_floating_ip.bootstrap_floatingip[0].address : ibm_is_instance.bootstrap_node.primary_network_interface[0].primary_ipv4_address
}
