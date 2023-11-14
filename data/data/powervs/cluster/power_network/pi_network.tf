locals {
  ids   = data.ibm_pi_dhcps.dhcp_services.servers[*].dhcp_id
  names = data.ibm_pi_dhcps.dhcp_services.servers[*].network_name
}

data "ibm_pi_dhcps" "dhcp_services" {
  pi_cloud_instance_id = var.cloud_instance_id
}

resource "ibm_pi_dhcp" "new_dhcp_service" {
  count                = 1
  pi_cloud_instance_id = var.cloud_instance_id
  pi_cidr              = var.machine_cidr
  pi_dns_server        = var.dns_server
  pi_dhcp_snat_enabled = var.enable_snat
  # the pi_dhcp_name param will be prefixed by the DHCP ID when created, so keep it short here:
  pi_dhcp_name = var.cluster_id
}

data "ibm_pi_dhcp" "dhcp_service" {
  pi_cloud_instance_id = var.cloud_instance_id
  pi_dhcp_id           = ibm_pi_dhcp.new_dhcp_service[0].dhcp_id
}
