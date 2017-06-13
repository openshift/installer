resource "openstack_compute_secgroup_v2" "master" {
  name        = "${var.cluster_name}_master"
  description = "security group for k8s masters"

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

resource "openstack_compute_secgroup_v2" "node" {
  name        = "${var.cluster_name}_worker"
  description = "security group for k8s nodes"

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
  name        = "${var.cluster_name}_master_self_hosted_etcd"
  description = "security group for self hosted etcd"

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

resource "openstack_compute_secgroup_v2" "etcd" {
  name        = "${var.cluster_name}_etcd_group"
  description = "security group for etcd"

  rule {
    from_port   = 2379
    to_port     = 2380
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "default" {
  name        = "${var.cluster_name}_default"
  description = "security group defaults: SSH and ping"

  rule {
    from_port   = 22
    to_port     = 22
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
