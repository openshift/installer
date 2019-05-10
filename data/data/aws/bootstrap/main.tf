resource "aws_s3_bucket" "ignition" {
  acl = "private"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    var.tags,
  )

  lifecycle {
    ignore_changes = all
  }
}

resource "aws_s3_bucket_object" "ignition" {
  bucket  = aws_s3_bucket.ignition.id
  key     = "bootstrap.ign"
  content = var.ignition
  acl     = "private"

  server_side_encryption = "AES256"

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap"
    },
    var.tags,
  )

  lifecycle {
    ignore_changes = all
  }
}

data "ignition_config" "redirect" {
  replace {
    source = "s3://${aws_s3_bucket.ignition.id}/bootstrap.ign"
  }
}

resource "aws_iam_instance_profile" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-profile"

  role = aws_iam_role.bootstrap.name
}

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

  tags = merge(
    {
      "Name" = "${var.cluster_id}-bootstrap-role"
    },
    var.tags,
  )
}

resource "aws_iam_role_policy" "bootstrap" {
  name = "${var.cluster_id}-bootstrap-policy"
  role = aws_iam_role.bootstrap.id

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
  ami = var.ami

  iam_instance_profile        = aws_iam_instance_profile.bootstrap.name
  instance_type               = var.instance_type
  subnet_id                   = var.subnet_id
  user_data                   = data.ignition_config.redirect.rendered
  vpc_security_group_ids      = flatten([var.vpc_security_group_ids, aws_security_group.bootstrap.id])
  associate_public_ip_address = true

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new OS release comes out.
    ignore_changes = [ami]
  }

  tags = merge(
    {
    "Name" = "${var.cluster_id}-bootstrap"
    },
    var.tags,
  )

  root_block_device {
    volume_type = var.volume_type
    volume_size = var.volume_size
    iops        = var.volume_type == "io1" ? var.volume_iops : 0
  }

  volume_tags = merge(
    {
    "Name" = "${var.cluster_id}-bootstrap-vol"
    },
    var.tags,
  )
}

resource "aws_lb_target_group_attachment" "bootstrap" {
  count = var.target_group_arns_length

  target_group_arn = var.target_group_arns[count.index]
  target_id        = aws_instance.bootstrap.private_ip
}

resource "aws_security_group" "bootstrap" {
  vpc_id = var.vpc_id

  timeouts {
    create = "20m"
  }

  tags = merge(
    {
    "Name" = "${var.cluster_id}-bootstrap-sg"
    },
    var.tags,
  )
}

resource "aws_security_group_rule" "ssh" {
  type              = "ingress"
  security_group_id = aws_security_group.bootstrap.id

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 22
  to_port     = 22
}

resource "aws_security_group_rule" "bootstrap_journald_gateway" {
  type              = "ingress"
  security_group_id = aws_security_group.bootstrap.id

  protocol    = "tcp"
  cidr_blocks = ["0.0.0.0/0"]
  from_port   = 19531
  to_port     = 19531
}

