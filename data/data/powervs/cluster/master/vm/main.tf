# Create the master instances
resource "ibm_pi_instance" "master" {
  count                = var.instance_count
  pi_memory            = var.memory
  pi_processors        = var.processors
  pi_instance_name     = "${var.cluster_id}-master-${count.index}"
  pi_proc_type         = var.proc_type
  pi_image_id          = var.image_id
  pi_sys_type          = var.sys_type
  pi_cloud_instance_id = var.cloud_instance_id
  pi_network {
    network_id = var.dhcp_network_id
  }
  pi_user_data     = base64encode(var.ignition)
  pi_key_pair_name = var.ssh_key_name
  pi_health_status = "WARNING"
}

resource "time_sleep" "wait_for_master_macs" {
  create_duration = "3m"

  depends_on = [ibm_pi_instance.master]
}

data "ibm_pi_dhcp" "dhcp_service_refresh" {
  depends_on           = [time_sleep.wait_for_master_macs]
  pi_cloud_instance_id = var.cloud_instance_id
  pi_dhcp_id           = var.dhcp_id
}

locals {
  macs       = flatten(ibm_pi_instance.master[*].pi_network[0].mac_address)
  master_ips = [for lease in data.ibm_pi_dhcp.dhcp_service_refresh.leases : lease.instance_ip if contains(local.macs, lease.instance_mac)]
}
