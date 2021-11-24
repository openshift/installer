resource "alicloud_security_group" "sg_master" {
  name              = "${local.prefix}-sg-master"
  description       = local.description
  resource_group_id = var.resource_group_id
  vpc_id            = alicloud_vpc.vpc.id
  tags = merge(
    {
      "Name" = "${local.prefix}-sg-master"
    },
    var.tags,
  )
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_mcs" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22623/22623"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_icmp" {
  type              = "ingress"
  ip_protocol       = "icmp"
  policy            = "accept"
  port_range        = "-1/-1"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ssh" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22/22"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_https" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "6443/6443"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_vxlan" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "4789/4789"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_geneve" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "6081/6081"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ike" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "500/500"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ike_nat_t" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "4500/4500"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ovndb" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "6641/6642"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_vxlan_from_worker" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "4789/4789"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_geneve_from_worker" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "6081/6081"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ike_from_worker" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "500/500"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ike_nat_t_from_worker" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "4500/4500"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_ovndb_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "6641/6642"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_internal" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "9000/9999"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_internal_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "9000/9999"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_internal_udp" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "9000/9999"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_internal_from_worker_udp" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "9000/9999"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_kube_scheduler" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "10259/10259"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_kube_scheduler_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10259/10259"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}


resource "alicloud_security_group_rule" "sg_rule_master_kube_controller_manager" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "10257/10257"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}


resource "alicloud_security_group_rule" "sg_rule_master_ingress_kube_controller_manager_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10257/10257"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_kubelet_insecure" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10250/10250"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_kubelet_insecure_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "10250/10250"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_etcd" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "2379/2380"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_services_tcp" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "30000/32767"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_services_tcp_from_worker" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "accept"
  port_range               = "30000/32767"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}

resource "alicloud_security_group_rule" "sg_rule_master_services_udp" {
  type              = "ingress"
  ip_protocol       = "udp"
  policy            = "accept"
  port_range        = "30000/32767"
  security_group_id = alicloud_security_group.sg_master.id
  cidr_ip           = var.vpc_cidr_block
}

resource "alicloud_security_group_rule" "sg_rule_master_ingress_services_udb_from_worker" {
  type                     = "ingress"
  ip_protocol              = "udp"
  policy                   = "accept"
  port_range               = "30000/32767"
  security_group_id        = alicloud_security_group.sg_master.id
  source_security_group_id = alicloud_security_group.sg_worker.id
}
