output "control_plane_ips" {
  value = ibm_is_instance.master_node[*].primary_network_interface[0].primary_ipv4_address
}
