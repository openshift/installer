locals {
  ids               = data.ibm_pi_dhcps.dhcp_services.servers[*].dhcp_id
  names             = data.ibm_pi_dhcps.dhcp_services.servers[*].network_name
  dhcp_id_from_name = var.pvs_network_name == "" ? "" : matchkeys(local.ids, local.names, [var.pvs_network_name])[0]
}

data "ibm_pi_dhcps" "dhcp_services" {
  pi_cloud_instance_id = var.cloud_instance_id
}

resource "ibm_pi_dhcp" "new_dhcp_service" {
  count                  = var.pvs_network_name == "" ? 1 : 0
  pi_cloud_instance_id   = var.cloud_instance_id
  pi_cloud_connection_id = data.ibm_pi_cloud_connection.cloud_connection.id
  pi_cidr                = var.machine_cidr
  pi_dns_server          = var.dns_server
  # the pi_dhcp_name param will be prefixed by the DHCP ID when created, so keep it short here:
  pi_dhcp_name = var.cluster_id
}

resource "ibm_pi_cloud_connection" "new_cloud_connection" {
  count                              = var.cloud_conn_name == "" ? 1 : 0
  pi_cloud_instance_id               = var.cloud_instance_id
  pi_cloud_connection_name           = "cloud-con-${var.cluster_id}"
  pi_cloud_connection_speed          = 50
  pi_cloud_connection_global_routing = true
  pi_cloud_connection_vpc_enabled    = true
  pi_cloud_connection_vpc_crns       = [var.vpc_crn]
}

data "ibm_pi_cloud_connection" "cloud_connection" {
  pi_cloud_connection_name = var.cloud_conn_name == "" ? ibm_pi_cloud_connection.new_cloud_connection[0].pi_cloud_connection_name : var.cloud_conn_name
  pi_cloud_instance_id     = var.cloud_instance_id
}

data "ibm_pi_dhcp" "dhcp_service" {
  pi_cloud_instance_id = var.cloud_instance_id
  pi_dhcp_id           = var.pvs_network_name == "" ? ibm_pi_dhcp.new_dhcp_service[0].dhcp_id : local.dhcp_id_from_name
}
