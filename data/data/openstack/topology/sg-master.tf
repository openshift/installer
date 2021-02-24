resource "openstack_networking_secgroup_v2" "master" {
  count       = var.openstack_disable_sg ? 0 : 1
  name        = "${var.cluster_id}-master"
  tags        = ["openshiftClusterID=${var.cluster_id}"]
  description = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_mcs" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22623
  port_range_max    = 22623
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

# TODO(mandre) Explicitely enable egress

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "icmp"
  port_range_min = 0
  port_range_max = 0
  # FIXME(mandre) AWS only allows ICMP from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ssh" {
  count = var.openstack_disable_sg ? 0 : 1
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 22
  port_range_max = 22
  # FIXME(mandre) AWS only allows SSH from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_tcp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_mdns_udp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 5353
  port_range_max    = 5353
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https" {
  count = var.openstack_disable_sg ? 0 : 1
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 6443
  port_range_max = 6443
  # FIXME(mandre) AWS only allows API port from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ike" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 500
  port_range_max    = 500
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ike_nat_t" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4500
  port_range_max    = 4500
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_esp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "esp"
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_udp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_etcd" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_tcp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_udp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vrrp" {
  count = var.openstack_disable_sg ? 0 : 1
  direction = "ingress"
  ethertype = "IPv4"
  # Explicitly set the vrrp protocol number to prevent cases when the Neutron Plugin
  # is disabled and it cannot identify a number by name.
  protocol          = "112"
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master[0].id
  description       = local.description
}

