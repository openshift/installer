resource "aws_s3_bucket_object" "ignition" {
  bucket  = "${var.bucket}"
  key     = "bootstrap.ign"
  content = "${var.ignition}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${var.tags}"

  lifecycle {
    ignore_changes = ["*"]
  }
}

data "ignition_config" "redirect" {
  replace {
    source = "s3://${var.bucket}/bootstrap.ign"
  }
}

resource "aws_iam_instance_profile" "bootstrap" {
  name = "${var.cluster_name}-bootstrap-profile"

  role = "${var.iam_role == "" ?
    join("|", aws_iam_role.bootstrap.*.name) :
    join("|", data.aws_iam_role.bootstrap.*.name)
  }"
}

data "aws_iam_role" "bootstrap" {
  count = "${var.iam_role == "" ? 0 : 1}"
  name  = "${var.iam_role}"
}

resource "aws_iam_role" "bootstrap" {
  count = "${var.iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}-bootstrap-role"
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

resource "aws_iam_role_policy" "bootstrap" {
  count = "${var.iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}-bootstrap-policy"
  role  = "${aws_iam_role.bootstrap.id}"

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
  ami = "${var.ami}"

  iam_instance_profile        = "${aws_iam_instance_profile.bootstrap.name}"
  instance_type               = "${var.instance_type}"
  subnet_id                   = "${var.subnet_id}"
  user_data                   = "${data.ignition_config.redirect.rendered}"
  vpc_security_group_ids      = ["${var.vpc_security_group_ids}"]
  associate_public_ip_address = "${var.associate_public_ip_address}"

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new OS release comes out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
    "kubernetes.io/cluster/${var.cluster_name}", "owned",
  ), var.tags)}"

  root_block_device {
    volume_type = "${var.volume_type}"
    volume_size = "${var.volume_size}"
    iops        = "${var.volume_type == "io1" ? var.volume_iops : 0}"
  }

  volume_tags = "${var.tags}"
}

resource "aws_elb_attachment" "bootstrap" {
  count    = "${var.elbs_length}"
  elb      = "${var.elbs[count.index]}"
  instance = "${aws_instance.bootstrap.id}"
}
