resource "openstack_networking_secgroup_rule_v2" "etcd" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 2379
  port_range_max    = 2380
  protocol          = "tcp"
  remote_ip_prefix  = "${var.cluster_cidr}"
  security_group_id = "${var.secgroup_id}"
}
