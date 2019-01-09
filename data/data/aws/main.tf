locals {
  private_zone_id = "${aws_route53_zone.int.zone_id}"
}

provider "aws" {
  region = "${var.aws_region}"
}

module "bootstrap" {
  source = "./bootstrap"

  ami                      = "${var.aws_ec2_ami_override}"
  cluster_name             = "${var.cluster_name}"
  iam_role                 = "${var.aws_controlplane_iam_role_name}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vpc.controlplane_subnet_ids[0]}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  vpc_security_group_ids   = "${list(module.vpc.controlplane_sg_id)}"

  tags = "${merge(map(
      "Name", "${var.cluster_name}-bootstrap",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.aws_extra_tags)}"
}

module "controlplane" {
  source = "./controlplane"

  base_domain              = "${var.base_domain}"
  cluster_id               = "${var.cluster_id}"
  cluster_name             = "${var.cluster_name}"
  ec2_type                 = "${var.aws_controlplane_ec2_type}"
  extra_tags               = "${var.aws_extra_tags}"
  instance_count           = "${var.controlplane_count}"
  controlplane_iam_role    = "${var.aws_controlplane_iam_role_name}"
  controlplane_sg_ids      = "${list(module.vpc.controlplane_sg_id)}"
  root_volume_iops         = "${var.aws_controlplane_root_volume_iops}"
  root_volume_size         = "${var.aws_controlplane_root_volume_size}"
  root_volume_type         = "${var.aws_controlplane_root_volume_type}"
  subnet_ids               = "${module.vpc.controlplane_subnet_ids}"
  target_group_arns        = "${module.vpc.aws_lb_target_group_arns}"
  target_group_arns_length = "${module.vpc.aws_lb_target_group_arns_length}"
  ec2_ami                  = "${var.aws_ec2_ami_override}"
  user_data_ign            = "${var.ignition_controlplane}"
}

module "iam" {
  source = "./iam"

  cluster_name     = "${var.cluster_name}"
  compute_iam_role = "${var.aws_compute_iam_role_name}"
}

module "dns" {
  source = "./route53"

  api_external_lb_dns_name = "${module.vpc.aws_lb_api_external_dns_name}"
  api_external_lb_zone_id  = "${module.vpc.aws_lb_api_external_zone_id}"
  api_internal_lb_dns_name = "${module.vpc.aws_lb_api_internal_dns_name}"
  api_internal_lb_zone_id  = "${module.vpc.aws_lb_api_internal_zone_id}"
  base_domain              = "${var.base_domain}"
  cluster_name             = "${var.cluster_name}"
  controlplane_count       = "${var.controlplane_count}"
  private_zone_id          = "${local.private_zone_id}"
  extra_tags               = "${var.aws_extra_tags}"
}

module "vpc" {
  source = "./vpc"

  base_domain  = "${var.base_domain}"
  cidr_block   = "${var.machine_cidr}"
  cluster_id   = "${var.cluster_id}"
  cluster_name = "${var.cluster_name}"
  region       = "${var.aws_region}"

  extra_tags = "${var.aws_extra_tags}"
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = "${var.controlplane_count}"
  type    = "A"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "${var.cluster_name}-etcd-${count.index}"
  records = ["${module.controlplane.ip_addresses[count.index]}"]
}

resource "aws_route53_record" "etcd_cluster" {
  type    = "SRV"
  ttl     = "60"
  zone_id = "${local.private_zone_id}"
  name    = "_etcd-server-ssl._tcp.${var.cluster_name}"
  records = ["${formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)}"]
}

resource "aws_route53_zone" "int" {
  name          = "${var.base_domain}"
  force_destroy = true

  vpc {
    vpc_id = "${module.vpc.vpc_id}"
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}_int",
      "KubernetesCluster", "${var.cluster_name}",
      "openshiftClusterID", "${var.cluster_id}"
    ), var.aws_extra_tags)}"
}
