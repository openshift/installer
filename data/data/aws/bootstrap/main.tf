# AWS permission: s3:CreateBucket
# AWS permission: s3:DeleteBucket
resource "aws_s3_bucket" "ignition" {
  acl = "private"

  # AWS permission: s3:GetBucketTagging
  # AWS permission: s3:PutBucketTagging
  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap",
  ), var.tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

# AWS permission: s3:GetObject
# AWS permission: s3:DeleteObject
resource "aws_s3_bucket_object" "ignition" {
  bucket  = "${aws_s3_bucket.ignition.id}"
  key     = "bootstrap.ign"
  content = "${var.ignition}"
  acl     = "private"

  server_side_encryption = "AES256"

  # AWS permission: s3:GetObjectTagging
  # AWS permission: s3:PutObjectTagging
  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap",
  ), var.tags)}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

data "ignition_config" "redirect" {
  replace {
    source = "s3://${aws_s3_bucket.ignition.id}/bootstrap.ign"
  }
}

# AWS permission: iam:CreateInstanceProfile
# AWS permission: iam:DeleteInstanceProfile
resource "aws_iam_instance_profile" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-profile"

  role = "${aws_iam_role.bootstrap.name}"
}

# AWS permission: iam:CreateRole
# AWS permission: iam:DeleteRole
resource "aws_iam_role" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-role"
  path = "/"

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

  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap-role",
  ), var.tags)}"
}

# AWS permission: iam:PutRolePolicy
# AWS permission: iam:DeleteRolePolicy
resource "aws_iam_role_policy" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-policy"
  role = "${aws_iam_role.bootstrap.id}"

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

# AWS permission: ec2:RunInstances
# AWS permission: ec2:TerminateInstances
resource "aws_instance" "bootstrap" {
  ami = "${var.ami}"

  iam_instance_profile        = "${aws_iam_instance_profile.bootstrap.name}"
  instance_type               = "${var.instance_type}"
  subnet_id                   = "${var.subnet_id}"
  user_data                   = "${data.ignition_config.redirect.rendered}"
  vpc_security_group_ids      = ["${var.vpc_security_group_ids}", "${aws_security_group.bootstrap.id}"]
  associate_public_ip_address = true

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new OS release comes out.
    ignore_changes = ["ami"]
  }

  # AWS permission: ec2:CreateTags
  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap",
  ), var.tags)}"

  # AWS permission: ec2:CreateVolume
  # AWS permission: ec2:DeleteVolume
  root_block_device {
    volume_type = "${var.volume_type}"
    volume_size = "${var.volume_size}"
    iops        = "${var.volume_type == "io1" ? var.volume_iops : 0}"
  }

  # AWS permission: ec2:CreateTags
  volume_tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap-vol",
  ), var.tags)}"
}

# AWS permission: elasticloadbalancing:RegisterTargets
# AWS permission: elasticloadbalancing:DeregisterTargets
resource "aws_lb_target_group_attachment" "bootstrap" {
  count = "${var.target_group_arns_length}"

  target_group_arn = "${var.target_group_arns[count.index]}"
  target_id        = "${aws_instance.bootstrap.private_ip}"
}

# AWS permission: ec2:CreateSecurityGroup
# AWS permission: ec2:DeleteSecurityGroup
resource "aws_security_group" "bootstrap" {
  vpc_id = "${var.vpc_id}"

  timeouts {
    create = "20m"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-bootstrap-sg",
  ), var.tags)}"
}

# AWS permission: ec2:UpdateSecurityGroupRuleDescriptionsIngress
resource "aws_security_group_rule" "ssh" {
  type              = "ingress"
  security_group_id = "${aws_security_group.bootstrap.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 22
  to_port     = 22
}

# AWS permission: ec2:UpdateSecurityGroupRuleDescriptionsIngress
resource "aws_security_group_rule" "bootstrap_journald_gateway" {
  type              = "ingress"
  security_group_id = "${aws_security_group.bootstrap.id}"

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 19531
  to_port     = 19531
}
