resource "aws_security_group" "controlplane" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}_controlplane_sg",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_security_group_rule" "controlplane_mcs" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 49500
  to_port     = 49500
}

resource "aws_security_group_rule" "controlplane_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.controlplane.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "controlplane_ingress_icmp" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol    = "icmp"
  cidr_blocks = ["${data.aws_vpc.cluster_vpc.cidr_block}"]
  from_port   = 0
  to_port     = 0
}

resource "aws_security_group_rule" "controlplane_ingress_ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "controlplane_ingress_http" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol    = "tcp"
  cidr_blocks = ["${data.aws_vpc.cluster_vpc.cidr_block}"]
  from_port   = 80
  to_port     = 80
}

resource "aws_security_group_rule" "controlplane_ingress_https" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol    = "tcp"
  cidr_blocks = ["${data.aws_vpc.cluster_vpc.cidr_block}"]
  from_port   = 6443
  to_port     = 6445
}

resource "aws_security_group_rule" "controlplane_ingress_heapster" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 4194
  to_port   = 4194
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_heapster_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 4194
  to_port   = 4194
}

resource "aws_security_group_rule" "controlplane_ingress_flannel" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_flannel_from_etcd" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "controlplane_ingress_flannel_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "controlplane_ingress_internal" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9990
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_internal_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9990
}

resource "aws_security_group_rule" "controlplane_ingress_kube_scheduler" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 10251
  to_port   = 10251
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_kube_scheduler_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 10251
  to_port   = 10251
}

resource "aws_security_group_rule" "controlplane_ingress_kube_controller_manager" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 10252
  to_port   = 10252
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_kube_controller_manager_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 10252
  to_port   = 10252
}

resource "aws_security_group_rule" "controlplane_ingress_kubelet_insecure" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_kubelet_insecure_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
}

resource "aws_security_group_rule" "controlplane_ingress_kubelet_secure" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 10255
  to_port   = 10255
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_kubelet_secure_from_compute" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.compute.id}"

  protocol  = "tcp"
  from_port = 10255
  to_port   = 10255
}

resource "aws_security_group_rule" "controlplane_ingress_etcd" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 2379
  to_port   = 2380
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_bootstrap_etcd" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 12379
  to_port   = 12380
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_services" {
  type              = "ingress"
  security_group_id = "${aws_security_group.controlplane.id}"

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "controlplane_ingress_services_from_console" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.console.id}"

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
}

resource "aws_security_group_rule" "controlplane_ingress_from_etcd" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.controlplane.id}"
  source_security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "tcp"
  from_port = 0
  to_port   = 65535
}
