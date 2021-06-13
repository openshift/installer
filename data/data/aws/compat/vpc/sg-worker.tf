resource "aws_security_group" "worker" {
  vpc_id      = data.aws_vpc.cluster_vpc.id
  description = local.description

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
  description       = local.description

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "worker_ingress_icmp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol    = "icmp"
  cidr_blocks = var.cidr_blocks
  from_port   = -1
  to_port     = -1
}

resource "aws_security_group_rule" "worker_ingress_ssh" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol    = "tcp"
  cidr_blocks = var.cidr_blocks
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "worker_ingress_vxlan" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_vxlan_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "worker_ingress_geneve" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 6081
  to_port   = 6081
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_ike" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 500
  to_port   = 500
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_ike_nat_t" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 4500
  to_port   = 4500
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_esp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = 50
  from_port = 0
  to_port   = 0
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_geneve_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 6081
  to_port   = 6081
}

resource "aws_security_group_rule" "worker_ingress_ike_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 500
  to_port   = 500
}

resource "aws_security_group_rule" "worker_ingress_nat_t_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 4500
  to_port   = 4500
}

resource "aws_security_group_rule" "worker_ingress_esp_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = 50
  from_port = 0
  to_port   = 0
}

resource "aws_security_group_rule" "worker_ingress_internal" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_internal_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "worker_ingress_internal_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_internal_from_master_udp" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_kubelet_insecure_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
}

resource "aws_security_group_rule" "worker_ingress_services_tcp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_services_tcp_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
}

resource "aws_security_group_rule" "worker_ingress_services_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.worker.id
  description       = local.description

  protocol  = "udp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "worker_ingress_services_udp_from_master" {
  type                     = "ingress"
  security_group_id        = aws_security_group.worker.id
  source_security_group_id = aws_security_group.master.id
  description              = local.description

  protocol  = "udp"
  from_port = 30000
  to_port   = 32767
}
