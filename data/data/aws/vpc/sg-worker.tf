resource "aws_security_group" "worker" {
  vpc_id = "${data.aws_vpc.cluster_vpc.id}"

  tags = "${merge(map(
      "Name", "${var.cluster_id}-worker-sg",
    ), var.tags)}"
}

resource "aws_security_group_rule" "worker_egress" {
  type              = "egress"
  security_group_id = "${aws_security_group.worker.id}"

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "worker_ingress_icmp" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol    = "icmp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 0
  to_port     = 0
}

resource "aws_security_group_rule" "worker_ingress_ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "worker_ingress_http" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 80
  to_port     = 80
}

resource "aws_security_group_rule" "worker_ingress_https" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 443
  to_port     = 443
}

resource "aws_security_group_rule" "worker_ingress_heapster" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 4194
  to_port   = 4194
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_heapster_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "tcp"
  from_port = 4194
  to_port   = 4194
}

resource "aws_security_group_rule" "worker_ingress_vxlan" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_vxlan_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "worker_ingress_internal" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_internal_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
}

resource "aws_security_group_rule" "worker_ingress_kubelet_secure" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 10255
  to_port   = 10255
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_kubelet_secure_from_master" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.master.id}"

  protocol  = "tcp"
  from_port = 10255
  to_port   = 10255
}

resource "aws_security_group_rule" "worker_ingress_services" {
  type              = "ingress"
  security_group_id = "${aws_security_group.worker.id}"

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_services_from_console" {
  type                     = "ingress"
  security_group_id        = "${aws_security_group.worker.id}"
  source_security_group_id = "${aws_security_group.console.id}"

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
}
