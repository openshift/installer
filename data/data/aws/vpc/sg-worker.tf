resource "aws_security_group" "worker" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  timeouts {
    create = "20m"
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-worker-sg"
    },
    var.tags,
  )
}

resource "aws_security_group_rule" "worker_egress" {
  type              = "egress"
  security_group_id = aws_security_group.worker.id

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "worker_ingress_icmp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol    = "icmp"
  cidr_blocks = [data.aws_vpc.cluster_vpc.cidr_block]
  from_port   = -1
  to_port     = -1
}

resource "aws_security_group_rule" "worker_ingress_ssh" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol    = "tcp"
  cidr_blocks = [data.aws_vpc.cluster_vpc.cidr_block]
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "worker_ingress_vxlan" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_vxlan_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "worker_ingress_internal" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_internal_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "worker_ingress_internal_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_internal_from_master_udp" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
}

resource "aws_security_group_rule" "worker_ingress_services_tcp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_services_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_bgp" {
  count                    = "${var.allow_bgp}"
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 179
  to_port   = 179
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_bgp_from_master" {
  count                    = "${var.allow_bgp}"
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 179
  to_port   = 179
}

resource "aws_security_group_rule" "worker_ingress_ipip" {
  count             = "${var.allow_ipip}"
  type              = "ingress"
  security_group_id = aws_security_group.worker.id

  protocol  = "4"
  from_port = -1
  to_port   = -1
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_ipip_from_master" {
  count             = "${var.allow_ipip}"
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id

  protocol  = "4"
  from_port = -1
  to_port   = -1
}
