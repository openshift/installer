resource "openstack_networking_secgroup_v2" "api" {
  name                 = "${var.cluster_name}_api_sg"
  description          = "tectonicClusterID: ${var.cluster_id}"
  delete_default_rules = true
}

resource "openstack_networking_secgroup_rule_v2" "api_egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}

resource "openstack_networking_secgroup_rule_v2" "api_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}
