resource "openstack_networking_secgroup_v2" "worker" {
  name = "${var.cluster_id}-worker"
  tags = ["openshiftClusterID=${var.cluster_id}"]
}

# TODO(mandre) Explicitely enable egress

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_icmp" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "icmp"
  port_range_min = 0
  port_range_max = 0
  # FIXME(mandre) AWS only allows ICMP from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_mdns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 5353
  port_range_max    = 5353
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_http" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_https" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_router" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 1936
  port_range_max    = 1936
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_vxlan" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_geneve" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_ike" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 500
  port_range_max    = 500
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_ike_nat_t" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4500
  port_range_max    = 4500
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_esp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "esp"
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_internal" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_internal_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_kubelet_insecure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_services_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_services_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}

resource "openstack_networking_secgroup_rule_v2" "worker_ingress_vrrp" {
  direction = "ingress"
  ethertype = "IPv4"
  # Explicitly set the vrrp protocol number to prevent cases when the Neutron Plugin
  # is disabled and it cannot identify a number by name.
  protocol          = "112"
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.worker.id
}
