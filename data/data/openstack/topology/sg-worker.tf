resource "openstack_networking_secgroup_v2" "worker" {
  name = "worker"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_heapster" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 4194
  port_range_max    = 4194
  protocol          = "tcp"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_heapster_from_master" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 4194
  port_range_max    = 4194
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_flannel" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_flannel_from_master" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_node_exporter" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9100
  port_range_max    = 9100
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_node_exporter_from_master" {
  direction         = "ingress"
  protocol          = "tcp"
  ethertype         = "IPv4"
  port_range_min    = 9100
  port_range_max    = 9100
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_kubelet_insecure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_kubelet_insecure_from_master" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_kubelet_secure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10255
  port_range_max    = 10255
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_kubelet_secure_from_master" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10255
  port_range_max    = 10255
  remote_group_id   = "${openstack_networking_secgroup_v2.master.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_services" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_services_from_console" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_group_id   = "${openstack_networking_secgroup_v2.console.id}"
  security_group_id = "${openstack_networking_secgroup_v2.worker.id}"
}
