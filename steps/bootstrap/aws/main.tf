locals {
  private_endpoints = "${var.tectonic_aws_endpoints != "public"}"
  public_endpoints  = "${var.tectonic_aws_endpoints != "private"}"
}

provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

resource "aws_s3_bucket_object" "ignition_bootstrap" {
  bucket  = "${local.s3_bucket}"
  key     = "bootstrap.ign"
  content = "${local.ignition_bootstrap}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-bootstrap",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

data "ignition_config" "bootstrap_redirect" {
  replace {
    source = "s3://${local.s3_bucket}/bootstrap.ign"
  }
}

module "ami" {
  source = "../../../modules/aws/ami"

  region          = "${var.tectonic_aws_region}"
  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

resource "aws_iam_instance_profile" "bootstrap" {
  name = "${var.tectonic_cluster_name}-bootstrap-profile"

  role = "${var.tectonic_aws_master_iam_role_name == "" ? join("|", aws_iam_role.bootstrap_role.*.name) : join("|", data.aws_iam_role.bootstrap_role.*.name)}"
}

data "aws_iam_role" "bootstrap_role" {
  count = "${var.tectonic_aws_master_iam_role_name == "" ? 0 : 1}"
  name  = "${var.tectonic_aws_master_iam_role_name}"
}

resource "aws_iam_role" "bootstrap_role" {
  count = "${var.tectonic_aws_master_iam_role_name == "" ? 1 : 0}"
  name  = "${var.tectonic_cluster_name}-bootstrap-role"
  path  = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.amazonaws.com"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF
}

resource "aws_iam_role_policy" "bootstrap_policy" {
  count = "${var.tectonic_aws_master_iam_role_name == "" ? 1 : 0}"
  name  = "${var.tectonic_cluster_name}-bootstrap-policy"
  role  = "${aws_iam_role.bootstrap_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:Describe*",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:AttachVolume",
      "Resource": "*"
    },
    {
      "Effect": "Allow",
      "Action": "ec2:DetachVolume",
      "Resource": "*"
    },
    {
      "Action" : [
        "s3:GetObject"
      ],
      "Resource": "arn:aws:s3:::*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_instance" "bootstrap" {
  ami = "${coalesce(var.tectonic_aws_ec2_ami_override, module.ami.id)}"

  iam_instance_profile        = "${aws_iam_instance_profile.bootstrap.name}"
  instance_type               = "${var.tectonic_aws_master_ec2_type}"
  subnet_id                   = "${local.subnet_ids[0]}"
  user_data                   = "${data.ignition_config.bootstrap_redirect.rendered}"
  vpc_security_group_ids      = ["${concat(var.tectonic_aws_master_extra_sg_ids, list(local.sg_id))}"]
  associate_public_ip_address = "${local.public_endpoints}"

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new OS release comes out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-bootstrap",
      "kubernetes.io/cluster/${var.tectonic_cluster_name}", "owned",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"

  root_block_device {
    volume_type = "${var.tectonic_aws_master_root_volume_type}"
    volume_size = "${var.tectonic_aws_master_root_volume_size}"
    iops        = "${var.tectonic_aws_master_root_volume_type == "io1" ? var.tectonic_aws_master_root_volume_iops : 0}"
  }

  volume_tags = "${merge(map(
    "Name", "${var.tectonic_cluster_name}-bootstrap",
    "kubernetes.io/cluster/${var.tectonic_cluster_name}", "owned",
    "tectonicClusterID", "${var.tectonic_cluster_id}"
  ), var.tectonic_aws_extra_tags)}"
}

resource "aws_elb_attachment" "bootstrap" {
  count    = "${length(local.aws_lbs)}"
  elb      = "${local.aws_lbs[count.index]}"
  instance = "${aws_instance.bootstrap.0.id}"
}
