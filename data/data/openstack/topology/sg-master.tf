resource "openstack_networking_secgroup_v2" "master" {
  name = "${var.cluster_id}-master"
  tags = ["openshiftClusterID=${var.cluster_id}"]
}

// We can't create all security group rules at once because it may lead to
// conflicts in Neutron. Therefore we have to create rules sequentially by
// setting explicit dependencies between them.
// For more information: https://github.com/hashicorp/terraform/issues/7519

// FIXME(mfedosin): ideally we need to resolve this in the OpenStack Terraform
// provider.
// Remove the dependencies when https://github.com/terraform-providers/terraform-provider-openstack/issues/952
// is fixed.

resource "openstack_networking_secgroup_rule_v2" "master_mcs" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22623
  port_range_max    = 22623
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id
}
#
# TODO(mandre) Explicitely enable egress

resource "openstack_networking_secgroup_rule_v2" "master_ingress_icmp" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "icmp"
  port_range_min = 0
  port_range_max = 0
  # FIXME(mandre) AWS only allows ICMP from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_mcs]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ssh" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 22
  port_range_max = 22
  # FIXME(mandre) AWS only allows SSH from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_icmp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_ssh]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_dns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_dns_tcp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_mdns_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 5353
  port_range_max    = 5353
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_dns_udp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_https" {
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 6443
  port_range_max = 6443
  # FIXME(mandre) AWS only allows API port from cidr_block
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_mdns_udp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vxlan" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 4789
  port_range_max    = 4789
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_https]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_geneve" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 6081
  port_range_max    = 6081
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_vxlan]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_ovndb" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 6641
  port_range_max    = 6642
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_geneve]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_ovndb]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_internal_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 9000
  port_range_max    = 9999
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_internal]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_scheduler" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10259
  port_range_max    = 10259
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_internal_udp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kube_controller_manager" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10257
  port_range_max    = 10257
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_kube_scheduler]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_kubelet_secure" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 10250
  port_range_max    = 10250
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_kube_controller_manager]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_etcd" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 2379
  port_range_max    = 2380
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_kubelet_secure]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_tcp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_etcd]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_services_udp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 30000
  port_range_max    = 32767
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_services_tcp]
}

resource "openstack_networking_secgroup_rule_v2" "master_ingress_vrrp" {
  direction = "ingress"
  ethertype = "IPv4"
  # Explicitly set the vrrp protocol number to prevent cases when the Neutron Plugin
  # is disabled and it cannot identify a number by name.
  protocol          = "112"
  remote_ip_prefix  = var.cidr_block
  security_group_id = openstack_networking_secgroup_v2.master.id

  depends_on = [openstack_networking_secgroup_rule_v2.master_ingress_services_udp]
}

