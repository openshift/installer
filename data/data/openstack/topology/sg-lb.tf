resource "openstack_networking_secgroup_v2" "mcs" {
  name = "mcs"
  tags = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
}

resource "openstack_networking_secgroup_rule_v2" "mcs_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 49500
  port_range_max    = 49500
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.mcs.id}"
}

resource "openstack_networking_secgroup_v2" "api" {
  name = "api"
  tags = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
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

resource "openstack_networking_secgroup_v2" "console" {
  name = "console"
  tags = ["${format("tectonicClusterID=%s", var.cluster_id)}"]
}

resource "openstack_networking_secgroup_rule_v2" "console_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.console.id}"
}

resource "openstack_networking_secgroup_rule_v2" "console_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.console.id}"
}
