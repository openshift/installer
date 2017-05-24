resource "openstack_compute_secgroup_v2" "etcd" {
  name        = "${var.cluster_name}_etcd_group"
  description = "security group for etcd: SSH and etcd client / cluster"
  count       = "${var.tectonic_experimental ? 0 : 1}"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 2379
    to_port     = 2380
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}
