output "dhcp_id" {
  value = ibm_pi_dhcp.dhcp_service.dhcp_id
}

output "dhcp_network_id" {
  value = ibm_pi_dhcp.dhcp_service.network
}
