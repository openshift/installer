resource "aws_security_group" "etcd" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}_etcd_sg",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"
}

resource "aws_security_group_rule" "etcd_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.etcd.id}"

  from_port   = 0
  cidr_blocks = ["0.0.0.0/0"]
  to_port     = 0
  protocol    = "-1"
}

resource "aws_security_group_rule" "etcd_ingress_icmp" {
  type              = "ingress"
  security_group_id = "${aws_security_group.etcd.id}"

  protocol    = "icmp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 0
  to_port     = 0
}

resource "aws_security_group_rule" "etcd_ingress_ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "tcp"
  from_port = 22
  to_port   = 22
  self      = true
}

resource "aws_security_group_rule" "etcd_ingress_etcd" {
  type              = "ingress"
  security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "tcp"
  from_port = 2379
  to_port   = 2379
  self      = true
}

resource "aws_security_group_rule" "etcd_ingress_peer" {
  type              = "ingress"
  security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "tcp"
  from_port = 2380
  to_port   = 2380
  self      = true
}

resource "aws_security_group_rule" "etcd_ingress_flannel" {
  type              = "ingress"
  security_group_id = "${aws_security_group.etcd.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "etcd_ingress_flannel_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.etcd.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "etcd_ingress_flannel_from_worker" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.etcd.id}"
  source_security_group_id = "${aws_security_group.worker.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "etcd_ingress_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.etcd.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "tcp"
  from_port = 0
  to_port   = 65535
}

resource "aws_security_group_rule" "etcd_ingress_from_worker" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.etcd.id}"
  source_security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 0
  to_port   = 65535
}
