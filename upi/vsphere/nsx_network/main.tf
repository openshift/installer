resource "nsxt_logical_switch" "${var.logical_switch}" {
  admin_state       = "UP"
  description       = "LS1 provisioned by Terraform"
  display_name      = "LS1"
  transport_zone_id = "${var.transport_zone_id}"
  replication_mode  = "MTEP"
}

resource "nsxt_logical_port" "logical_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${nsxt_logical_switch.${var.logical_switch}.id}"
}

rresource "nsxt_logical_router_downlink_port" "downlink_port" {
  description                   = "DP1 provisioned by Terraform"
  display_name                  = "DP1"
  logical_router_id             = "${var.t1_router}"
  linked_logical_switch_port_id = "${nsxt_logical_port.logical_port.id}"
  ip_address                    = "${var.logical_switch_ip_address}"
}

data "vsphere_network" "${var.logical_switch}" {
    name = "${nsxt_logical_switch.${var.logical_switch}.display_name}"
    datacenter_id = "${var.datacenter_id}"
    depends_on = ["nsxt_logical_switch.${var.logical_switch}"]
}

resource "nsxt_dhcp_server_profile" "dhcp_profile" {
  description                 = "dhcp_profile provisioned by Terraform"
  display_name                = "dhcp_profile"
  edge_cluster_id             = "${data.nsxt_edge_cluster.${var.nsx_edge_cluster}.id}"
  edge_cluster_member_indexes = [0, 1]
}

resource "nsxt_logical_dhcp_server" "logical_dhcp_server" {
  display_name    = "logical_dhcp_server"
  dhcp_profile_id = "${nsxt_dhcp_server_profile.dhcp_profile.id}"
  dhcp_server_ip  = "${var.dhcp_server_ip}/24"
  gateway_ip      = "${var.gateway_ip}"
}

resource "nsxt_logical_dhcp_port" "dhcp_port" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${nsxt_logical_switch.${var.logical_switch}.id}"
  dhcp_server_id    = "${nsxt_logical_dhcp_server.logical_dhcp_server.id}"
}

resource "nsxt_ip_block" "ip_block" {
  description  = "ip_block provisioned by Terraform"
  display_name = "ip_block"
  cidr         = "${var.ip_block_cidr}/24"
}

resource "nsxt_ip_block_subnet" "ip_block_subnet" {
  description = "ip_block_subnet"
  block_id    = "${nsxt_ip_block.ip_block.id}"
  size        = 16
}

resource "nsxt_ip_pool" "ip_pool" {
  description = "ip_pool provisioned by Terraform"
  display_name = "ip_pool"


  subnet {
    allocation_ranges = ["${var.ip_pool_start}", "${var.ip_pool_end}"]
    cidr              = "${var.ip_pool_cidr}/24"
    gateway_ip        = "${var.gateway_ip}"
    dns_suffix        = "${var.base_domain}"
    dns_nameservers   = ["${var.dns_nameservers}"]
  }
}