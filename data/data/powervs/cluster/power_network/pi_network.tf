resource "ibm_pi_dhcp" "dhcp_service" {
  pi_cloud_instance_id   = var.cloud_instance_id
  pi_cloud_connection_id = ibm_pi_cloud_connection.cloud_connection.cloud_connection_id
}

resource "ibm_pi_cloud_connection" "cloud_connection" {
  pi_cloud_instance_id            = var.cloud_instance_id
  pi_cloud_connection_name        = "cloud-con-${var.cluster_id}"
  pi_cloud_connection_speed       = 50
  pi_cloud_connection_vpc_enabled = true
  pi_cloud_connection_vpc_crns    = [var.vpc_crn]
}
