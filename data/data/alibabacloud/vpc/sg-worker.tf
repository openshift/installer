resource "alicloud_security_group" "sg_worker" {
  name        = "${local.prefix}-sg_worker"
  description = local.description
  vpc_id      = alicloud_vpc.vpc.id
  tags = merge(
    {
      "Name" = "${local.prefix}-sg-worker"
    },
    var.tags,
  )
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_icmp" {
  type              = "ingress"
  ip_protocol       = "icmp"
  policy            = "accept"
  port_range        = "-1/-1"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_ssh" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22/22"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_vxlan" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "4789/4789"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_geneve" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "6081/6081"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_ike" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "500/500"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_ike_nat_t" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "4500/4500"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_vxlan_from_master" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "4789/4789"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_geneve_from_master" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "6081/6081"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_ike_from_master" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "500/500"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_ike_nat_t_from_master" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "4500/4500"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_internal" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "9000/9999"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_internal_from_master" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "9000/9999"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_internal_udp" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "9000/9999"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_internal_from_master_udp" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "9000/9999"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_kubelet_insecure" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10250/10250"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_kubelet_insecure_from_master" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10250/10250"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_services_tcp" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "30000/32767"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_services_tcp_from_master" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "30000/32767"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}

resource "alicloud_security_group_rule" "sg_rule_worker_services_udp" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "30000/32767"
  security_group_id = alicloud_security_group.sg_worker.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_worker_ingress_services_udb_from_master" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "30000/32767"
  security_group_id        = alicloud_security_group.sg_worker.id
  source_security_group_id = alicloud_security_group.sg_master.id
}
