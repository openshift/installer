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
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  port_range_min    = 0
  port_range_max    = 0
  remote_ip_prefix  = "0.0.0.0/0"
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
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_mdns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 5353
  port_range_max    = 5353
  remote_ip_prefix  = "${var.cidr_block}"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6443
  port_range_max    = 6443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_group_id   = openstack_networking_secgroup_v2.worker.id
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
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
  port_range_min    = 10250
  port_range_max    = 10250
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure_from_worker" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
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

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  security_group_id = openstack_networking_secgroup_v2.master.id
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vrrp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = 112
  security_group_id = openstack_networking_secgroup_v2.master.id
}

