resource "nsxt_logical_port" "logical_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${var.logical_switch}"
}

resource "nsxt_logical_router_downlink_port" "downlink_port" {
  description                   = "DP1 provisioned by Terraform"
  display_name                  = "DP1"
  logical_router_id             = "${var.t1_router}"
  linked_logical_switch_port_id = "${nsxt_logical_port.logical_port.id}"
  ip_address                    = "${var.logical_switch_ip_address}"
}


resource "nsxt_dhcp_server_profile" "dhcp_profile" {
  description                 = "dhcp_profile provisioned by Terraform"
  display_name                = "dhcp_profile"
  edge_cluster_id             = "${var.nsx_edge_cluster}"
}

resource "nsxt_logical_dhcp_server" "logical_dhcp_server" {
  display_name    = "logical_dhcp_server"
  dhcp_profile_id = "${nsxt_dhcp_server_profile.dhcp_profile.id}"
  dhcp_server_ip  = "${var.dhcp_server_ip}/24"
  gateway_ip      = "${var.gateway_ip}"
  domain_name = "${var.base_domain}"
  dns_name_servers = ["${var.dns_nameservers}"]
}

resource "nsxt_logical_dhcp_port" "dhcp_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${var.logical_switch}"
  dhcp_server_id    = "${nsxt_logical_dhcp_server.logical_dhcp_server.id}"
}

resource "nsxt_dhcp_server_ip_pool" "dhcp_ip_pool" {
  display_name = "ip pool"
  description = "ip pool"
  logical_dhcp_server_id = "${nsxt_logical_dhcp_server.logical_dhcp_server.id}"
  gateway_ip = "${var.gateway_ip}"
  lease_time = 180
  error_threshold = 98
  warning_threshold = 70

  ip_range {
    start = "${var.ip_pool_start}"
    end = "${var.ip_pool_end}"
  }
}


resource "nsxt_ip_block" "ip_block" {
  description  = "ip_block provisioned by Terraform"
  display_name = "ip_block"
  cidr         = "${var.ip_block_cidr}"
}

resource "nsxt_ip_block_subnet" "ip_block_subnet" {
  description = "ip_block_subnet"
  block_id    = "${nsxt_ip_block.ip_block.id}"
  size        = 16
}

