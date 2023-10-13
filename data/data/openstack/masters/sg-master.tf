resource "openstack_networking_secgroup_v2" "master" {
  name        = "${var.cluster_id}-master"
  tags        = ["openshiftClusterID=${var.cluster_id}"]
  description = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_mcs" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22623
  port_range_max    = 22623
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_mcs_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 22623
  port_range_max    = 22623
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

# TODO(mandre) Explicitely enable egress

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "icmp"
  port_range_min = 0
  port_range_max = 0
  # FIXME(mandre) AWS only allows ICMP from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp_v6" {
  count          = length(var.machine_v6_cidrs)
  direction      = "ingress"
  ethertype      = "IPv6"
  protocol       = "ipv6-icmp"
  port_range_min = 0
  port_range_max = 0
  # FIXME(mandre) AWS only allows ICMP from cidr_block
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ssh" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ssh_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_tcp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_tcp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_api" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 6443
  port_range_max = 6443
  # FIXME(mandre) AWS only allows API port from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_api_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ike" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 500
  port_range_max    = 500
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ike_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 500
  port_range_max    = 500
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ike_nat_t" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4500
  port_range_max    = 4500
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_esp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "esp"
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_esp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "esp"
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_udp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_udp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_etcd" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_etcd_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_tcp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_tcp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_udp" {
  count             = length(var.machine_v4_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_udp_v6" {
  count             = length(var.machine_v6_cidrs)
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vrrp" {
  count     = length(var.machine_v4_cidrs)
  direction = "ingress"
  ethertype = "IPv4"
  # Explicitly set the vrrp protocol number to prevent cases when the Neutron Plugin
  # is disabled and it cannot identify a number by name.
  protocol          = "112"
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vrrp_v6" {
  count     = length(var.machine_v6_cidrs)
  direction = "ingress"
  ethertype = "IPv6"
  # Explicitly set the vrrp protocol number to prevent cases when the Neutron Plugin
  # is disabled and it cannot identify a number by name.
  protocol          = "112"
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_http" {
  count             = var.masters_schedulable ? 1 : 0
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https" {
  count             = var.masters_schedulable ? 1 : 0
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_http_v6" {
  count             = (var.masters_schedulable && length(var.machine_v6_cidrs) > 0) ? 1 : 0
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https_v6" {
  count             = (var.masters_schedulable && length(var.machine_v6_cidrs) > 0) ? 1 : 0
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "::/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_router" {
  count             = var.masters_schedulable ? length(var.machine_v4_cidrs) : 0
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 1936
  port_range_max    = 1936
  remote_ip_prefix  = element(var.machine_v4_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_router_v6" {
  count             = (var.masters_schedulable && length(var.machine_v6_cidrs) > 0) ? 1 : 0
  direction         = "ingress"
  ethertype         = "IPv6"
  protocol          = "tcp"
  port_range_min    = 1936
  port_range_max    = 1936
  remote_ip_prefix  = element(var.machine_v6_cidrs, count.index)
  security_group_id = openstack_networking_secgroup_v2.master.id
  description       = local.description
}