locals {
  tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    var.aws_extra_tags,
  )
  description = "Created By OpenShift Installer"

  public_endpoints = var.aws_publish_strategy == "External" ? true : false
  volume_type      = "gp2"
  volume_size      = 30
  volume_iops      = local.volume_type == "io1" ? 100 : 0

  // s3 object supports only 10 tags. The first 8 tags from
  // the list of aws_extra_tags are used for s3 object
  // slice function uses new_tag_len as excluding index
  new_tag_len    = 8
  sliced_tag_map = length(var.aws_extra_tags) <= local.new_tag_len ? var.aws_extra_tags : { for k in slice(keys(var.aws_extra_tags), 0, local.new_tag_len) : k => var.aws_extra_tags[k] }
  s3_object_tags = merge(
    {
      "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    },
    local.sliced_tag_map,
  )
}

provider "aws" {
  region = var.aws_region

  skip_region_validation = true

  endpoints {
    ec2     = lookup(var.custom_endpoints, "ec2", null)
    elb     = lookup(var.custom_endpoints, "elasticloadbalancing", null)
    iam     = lookup(var.custom_endpoints, "iam", null)
    route53 = lookup(var.custom_endpoints, "route53", null)
    s3      = lookup(var.custom_endpoints, "s3", null)
    sts     = lookup(var.custom_endpoints, "sts", null)
  }
}

data "aws_partition" "current" {}

data "aws_ebs_default_kms_key" "current" {}

resource "aws_s3_bucket" "ignition" {
  count         = var.aws_preserve_bootstrap_ignition ? 0 : 1
  bucket        = var.aws_ignition_bucket
  force_destroy = true

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    local.tags,
  )

  lifecycle {
    ignore_changes = all
  }
}

resource "aws_s3_object" "ignition" {
  count  = var.aws_preserve_bootstrap_ignition ? 0 : 1
  bucket = aws_s3_bucket.ignition[0].id
  key    = "bootstrap.ign"
  source = var.ignition_bootstrap_file

  server_side_encryption = "AES256"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    local.s3_object_tags,
  )

  lifecycle {
    ignore_changes = all
  }
}

resource "aws_iam_instance_profile" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-profile"

  role = var.aws_master_iam_role_name != "" ? var.aws_master_iam_role_name : aws_iam_role.bootstrap[0].name

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-profile"
    },
    local.tags,
  )
}

resource "aws_iam_role" "bootstrap" {
  count = var.aws_master_iam_role_name == "" ? 1 : 0

  name = "${var.cluster_id}-bootstrap-role"
  path = "/"

  assume_role_policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Principal": {
                "Service": "ec2.${data.aws_partition.current.dns_suffix}"
            },
            "Effect": "Allow",
            "Sid": ""
        }
    ]
}
EOF

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-role"
    },
    local.tags,
  )
}

resource "aws_iam_role_policy" "bootstrap" {
  count = var.aws_master_iam_role_name == "" ? 1 : 0
  name = "${var.cluster_id}-bootstrap-policy"
  role = aws_iam_role.bootstrap[0].id

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
    }
  ]
}
EOF
}

resource "aws_instance" "bootstrap" {
  ami = var.ami_id

  iam_instance_profile        = aws_iam_instance_profile.bootstrap.name
  instance_type               = var.aws_bootstrap_instance_type
  subnet_id                   = var.aws_publish_strategy == "External" ? var.public_subnet_ids[0] : var.private_subnet_ids[0]
  user_data                   = var.aws_bootstrap_stub_ignition
  vpc_security_group_ids      = [var.master_sg_id, aws_security_group.bootstrap.id]
  associate_public_ip_address = local.public_endpoints && var.aws_public_ipv4_pool == ""

  dynamic "instance_market_options" {
    for_each = var.aws_master_use_spot_instance ? [1] : []
    content {
      market_type = "spot"
      spot_options {
        spot_instance_type = "one-time"
      }
    }
  }

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new OS release comes out.
    ignore_changes = [ami]
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    local.tags,
  )

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = var.aws_bootstrap_instance_metadata_authentication
  }

  root_block_device {
    volume_type = local.volume_type
    volume_size = local.volume_size
    iops        = local.volume_iops
    encrypted   = true
    kms_key_id  = var.aws_master_root_volume_kms_key_id == "" ? data.aws_ebs_default_kms_key.current.key_arn : var.aws_master_root_volume_kms_key_id
  }

  volume_tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-vol"
    },
    local.tags,
  )

  depends_on = [
    aws_s3_object.ignition,
    # https://bugzilla.redhat.com/show_bug.cgi?id=1859153
    aws_iam_instance_profile.bootstrap,
  ]
}

resource "aws_lb_target_group_attachment" "bootstrap" {
  // Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot use a dynamic list for count
  // and therefore are force to implicitly assume that the list is of lb_target_group_arns_length - 1, in case there is no api_external
  count = local.public_endpoints ? var.lb_target_group_arns_length : var.lb_target_group_arns_length - 1

  target_group_arn = var.lb_target_group_arns[count.index]
  target_id        = aws_instance.bootstrap.private_ip
}

resource "aws_security_group" "bootstrap" {
  vpc_id      = var.vpc_id
  description = local.description

  timeouts {
    create = "20m"
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-sg"
    },
    local.tags,
  )
}

resource "aws_security_group_rule" "ssh" {
  type              = "ingress"
  security_group_id = aws_security_group.bootstrap.id
  description       = local.description

  protocol    = "tcp"
  cidr_blocks = local.public_endpoints ? ["0.0.0.0/0"] : var.machine_v4_cidrs
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "bootstrap_journald_gateway" {
  type              = "ingress"
  security_group_id = aws_security_group.bootstrap.id
  description       = local.description

  protocol    = "tcp"
  cidr_blocks = local.public_endpoints ? ["0.0.0.0/0"] : var.machine_v4_cidrs
  from_port   = 19531
  to_port     = 19531
}

resource "aws_eip" "bootstrap" {
  count            = var.aws_public_ipv4_pool == "" ? 0 : 1
  domain           = "vpc"
  instance         = aws_instance.bootstrap.id
  public_ipv4_pool = var.aws_public_ipv4_pool

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-eip"
    },
    local.tags,
  )

  depends_on = [aws_instance.bootstrap]
}
