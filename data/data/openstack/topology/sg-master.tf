resource "openstack_networking_secgroup_v2" "master" {
  name = "${var.cluster_id}-master"
  tags = ["openshiftClusterID=${var.cluster_id}"]
}

resource "openstack_networking_secgroup_rule_v2" "master_mcs" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22623
  port_range_max    = 22623
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "${var.cidr_block}"
  security_group_id = "${openstack_networking_secgroup_v2.master.id}"
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "${var.cidr_block}"
  security_group_id = "${openstack_networking_secgroup_v2.master.id}"
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6445
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_heapster" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 4194
  port_range_max    = 4194
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_heapster_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 4194
  port_range_max    = 4194
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_flannel" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_flannel_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_from_worker_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_insecure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_insecure_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10255
  port_range_max    = 10255
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10255
  port_range_max    = 10255
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_etcd" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_bootstrap_etcd" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 12379
  port_range_max    = 12380
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  security_group_id = openstack_networking_secgroup_v2.master.id
}

