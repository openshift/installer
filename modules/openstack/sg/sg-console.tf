resource "openstack_networking_secgroup_v2" "console" {
  name                 = "${var.cluster_name}_console_sg"
  description          = "tectonicClusterID: ${var.cluster_id}"
  delete_default_rules = true
}

resource "openstack_networking_secgroup_rule_v2" "console_egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.console.id}"
}

resource "openstack_networking_secgroup_rule_v2" "console_ingress_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.console.id}"
}

resource "openstack_networking_secgroup_rule_v2" "console_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.console.id}"
}
