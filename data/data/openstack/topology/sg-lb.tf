resource "openstack_networking_secgroup_v2" "api" {
  name = "api"
  tags = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_secgroup_rule_v2" "api_mcs" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 49500
  port_range_max    = 49500
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}

resource "openstack_networking_secgroup_rule_v2" "api_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}

resource "openstack_networking_secgroup_rule_v2" "api_ingress_ssh_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}

resource "openstack_networking_secgroup_rule_v2" "api_ingress_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.api.id}"
}
