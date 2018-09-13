provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

module "bootstrap" {
  source = "../../../modules/aws/bootstrap"

  ami                         = "${var.tectonic_aws_ec2_ami_override}"
  associate_public_ip_address = "${var.tectonic_aws_endpoints != "private"}"
  bucket                      = "${local.s3_bucket}"
  cluster_name                = "${var.tectonic_cluster_name}"
  elbs                        = "${local.aws_lbs}"
  iam_role                    = "${var.tectonic_aws_master_iam_role_name}"
  ignition                    = "${local.ignition_bootstrap}"
  subnet_id                   = "${local.subnet_ids[0]}"
  vpc_security_group_ids      = ["${concat(var.tectonic_aws_master_extra_sg_ids, list(local.sg_id))}"]

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-bootstrap",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}
