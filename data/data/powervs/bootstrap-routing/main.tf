locals {
  api_servers       = var.powervs_expose_bootstrap ? concat(var.control_plane_ips, [var.bootstrap_private_ip]) : var.control_plane_ips
  api_servers_count = length(local.api_servers)
}

provider "ibm" {
  ibmcloud_api_key = var.powervs_api_key
  region           = var.powervs_vpc_region
  zone             = var.powervs_vpc_zone
}

resource "ibm_is_lb_pool_member" "machine_config_member" {
  count          = local.api_servers_count
  lb             = var.lb_int_id
  pool           = var.machine_cfg_pool_id
  port           = 22623
  target_address = local.api_servers[count.index]
}

resource "ibm_is_lb_pool_member" "api_member_int" {
  count          = local.api_servers_count
  depends_on     = [ibm_is_lb_pool_member.machine_config_member]
  lb             = var.lb_int_id
  pool           = var.api_pool_int_id
  port           = 6443
  target_address = local.api_servers[count.index]
}

resource "ibm_is_lb_pool_member" "api_member" {
  count          = local.api_servers_count
  lb             = var.lb_ext_id
  pool           = var.api_pool_ext_id
  port           = 6443
  target_address = local.api_servers[count.index]
}
