data "nsxt_transport_zone" "overlay_tz" {
  display_name = "${var.overlay_tz}"
}

data "nsxt_edge_cluster" "${data.nsxt_edge_cluster}" {
  display_name = "${var.nsx_edge_cluster}"
}

resource "nsxt_logical_switch" "${var.logical_switch}" {
  admin_state       = "UP"
  description       = "LS created by Terraform"
  display_name      = "${var.logical_switch}"
  transport_zone_id = "${data.nsxt_transport_zone.overlay_tz.id}"
  replication_mode  = "MTEP"
}

resource "nsxt_logical_tier1_router" "tier1_router" {
  description                 = "Tier1 router provisioned by Terraform"
  display_name                = "TfTier1"
  failover_mode               = "PREEMPTIVE"
  high_availability_mode      = "ACTIVE_STANDBY"
  edge_cluster_id             = "${data.nsxt_edge_cluster.${var.nsx_edge_cluster}.id}"
  enable_router_advertisement = true
  advertise_connected_routes  = true
  advertise_static_routes     = false
  advertise_nat_routes        = true
}

resource "nsxt_logical_router_link_port_on_tier0" "link_port_tier0" {
  description       = "TIER0_PORT1 provisioned by Terraform"
  display_name      = "TIER0_PORT1"
  logical_router_id = "${data.nsxt_logical_tier0_router.${var.t0_router}.id}"
}

resource "nsxt_logical_router_link_port_on_tier1" "link_port_tier1" {
  description                   = "TIER1_PORT1 provisioned by Terraform"
  display_name                  = "TIER1_PORT1"
  logical_router_id             = "${var.t1_router}"
  linked_logical_router_port_id = "${nsxt_logical_router_link_port_on_tier0.link_port_tier0.id}"
}

resource "nsxt_logical_port" "logical_port1" {
  admin_state       = "UP"
  description       = "LP1 provisioned by Terraform"
  display_name      = "LP1"
  logical_switch_id = "${nsxt_logical_switch.${var.logical_switch}.id}"
}

resource "nsxt_logical_router_downlink_port" "downlink_port" {
  description                   = "DP1 provisioned by Terraform"
  display_name                  = "DP1"
  logical_router_id             = "${var.t1_router}"
  linked_logical_switch_port_id = "${nsxt_logical_port.logical_port1.id}"
  ip_address                    = "192.168.245.1/24"
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