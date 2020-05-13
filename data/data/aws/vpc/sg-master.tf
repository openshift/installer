resource "aws_security_group" "master" {
  vpc_id = data.aws_vpc.cluster_vpc.id

  timeouts {
    create = "20m"
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-master-sg"
    },
    var.tags,
  )
}

resource "aws_security_group_rule" "master_mcs" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol    = "tcp"
  cidr_blocks = var.cidr_blocks
  from_port   = 22623
  to_port     = 22623
}

resource "aws_security_group_rule" "master_egress" {
  type              = "egress"
  security_group_id = aws_security_group.master.id

  from_port   = 0
  to_port     = 0
  protocol    = "-1"
  cidr_blocks = ["0.0.0.0/0"]
}

resource "aws_security_group_rule" "master_ingress_icmp" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol    = "icmp"
  cidr_blocks = var.cidr_blocks
  from_port   = -1
  to_port     = -1
}

resource "aws_security_group_rule" "master_ingress_ssh" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol    = "tcp"
  cidr_blocks = var.cidr_blocks
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "master_ingress_https" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol    = "tcp"
  cidr_blocks = var.cidr_blocks
  from_port   = 6443
  to_port     = 6443
}

resource "aws_security_group_rule" "master_ingress_vxlan" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
  self      = true
}

resource "aws_security_group_rule" "master_ingress_vxlan_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 4789
  to_port   = 4789
}

resource "aws_security_group_rule" "master_ingress_geneve" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 6081
  to_port   = 6081
  self      = true
}

resource "aws_security_group_rule" "master_ingress_geneve_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 6081
  to_port   = 6081
}

resource "aws_security_group_rule" "master_ingress_ovndb" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 6641
  to_port   = 6642
  self      = true
}

resource "aws_security_group_rule" "master_ingress_ovndb_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 6641
  to_port   = 6642
}

resource "aws_security_group_rule" "master_ingress_internal" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "master_ingress_internal_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "master_ingress_internal_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
  self      = true
}

resource "aws_security_group_rule" "master_ingress_internal_from_worker_udp" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 9000
  to_port   = 9999
}

resource "aws_security_group_rule" "master_ingress_kube_scheduler" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 10259
  to_port   = 10259
  self      = true
}

resource "aws_security_group_rule" "master_ingress_kube_scheduler_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 10259
  to_port   = 10259
}

resource "aws_security_group_rule" "master_ingress_kube_controller_manager" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 10257
  to_port   = 10257
  self      = true
}

resource "aws_security_group_rule" "master_ingress_kube_controller_manager_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 10257
  to_port   = 10257
}

resource "aws_security_group_rule" "master_ingress_kubelet_secure" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
  self      = true
}

resource "aws_security_group_rule" "master_ingress_kubelet_secure_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 10250
  to_port   = 10250
}

resource "aws_security_group_rule" "master_ingress_etcd" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 2379
  to_port   = 2380
  self      = true
}

resource "aws_security_group_rule" "master_ingress_services_tcp" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "master_ingress_services_tcp_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "tcp"
  from_port = 30000
  to_port   = 32767
}

resource "aws_security_group_rule" "master_ingress_services_udp" {
  type              = "ingress"
  security_group_id = aws_security_group.master.id

  protocol  = "udp"
  from_port = 30000
  to_port   = 32767
  self      = true
}

resource "aws_security_group_rule" "master_ingress_services_udp_from_worker" {
  type                     = "ingress"
  security_group_id        = aws_security_group.master.id
  source_security_group_id = aws_security_group.worker.id

  protocol  = "udp"
  from_port = 30000
  to_port   = 32767
}
