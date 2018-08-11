resource "openstack_networking_secgroup_v2" "etcd" {
  name                 = "${var.cluster_name}_etcd_sg"
  description          = "tectonicClusterID: ${var.cluster_id}"
  delete_default_rules = true
}

resource "openstack_networking_secgroup_rule_v2" "etcd_egress" {
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.etcd.id}"
}

resource "openstack_networking_secgroup_rule_v2" "etcd_ingress_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.etcd.id}"
}

resource "openstack_networking_secgroup_rule_v2" "etcd_ingress_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.etcd.id}"
}

resource "openstack_networking_secgroup_rule_v2" "etcd_ingress_etcd" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  remote_group_id   = "${openstack_networking_secgroup_v2.etcd.id}"
  security_group_id = "${openstack_networking_secgroup_v2.etcd.id}"
}

resource "openstack_networking_secgroup_rule_v2" "etcd_ingress_etcd_from_master" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2379
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.etcd.id}"
}
