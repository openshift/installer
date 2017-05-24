resource "openstack_compute_secgroup_v2" "master" {
  name        = "${var.cluster_name}_${var.hostname_infix}"
  description = "security group for k8s masters"
  count       = "${var.hostname_infix == "master" ? 1 : 0}"

  // icmp
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  // SSH
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // k8s API
  rule {
    from_port   = 443
    to_port     = 443
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // cAdvisor
  rule {
    from_port   = 4194
    to_port     = 4194
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // flannel overlay vxlan
  rule {
    from_port   = 8472
    to_port     = 8472
    ip_protocol = "udp"
    cidr        = "0.0.0.0/0"
  }

  // kubelet API
  rule {
    from_port   = 10250
    to_port     = 10250
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "node" {
  name        = "${var.cluster_name}_${var.hostname_infix}"
  description = "security group for k8s nodes"
  count       = "${var.hostname_infix == "worker" ? 1 : 0}"

  // ICMP
  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  // SSH
  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // k8s API
  rule {
    from_port   = 443
    to_port     = 443
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // cAdvisor
  rule {
    from_port   = 4194
    to_port     = 4194
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // flannel
  rule {
    from_port   = 4789
    to_port     = 4789
    ip_protocol = "udp"
    cidr        = "0.0.0.0/0"
  }

  // flannel overlay vxlan
  rule {
    from_port   = 8472
    to_port     = 8472
    ip_protocol = "udp"
    cidr        = "0.0.0.0/0"
  }

  // kubelet API
  rule {
    from_port   = 10250
    to_port     = 10250
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // External services
  rule {
    from_port   = 30000
    to_port     = 32767
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "self_hosted_etcd" {
  name        = "${var.cluster_name}_${var.hostname_infix}_self_hosted_etcd"
  description = "security group for self hosted etcd"

  count = "${var.tectonic_experimental ? 1 : 0}"

  // etcd
  rule {
    from_port   = 2379
    to_port     = 2380
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  // bootstrap etcd
  rule {
    from_port   = 12379
    to_port     = 12380
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}
