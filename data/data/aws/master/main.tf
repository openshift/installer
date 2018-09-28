locals {
  arn = "aws"
}

resource "aws_iam_instance_profile" "master" {
  name = "${var.cluster_name}-master-profile"

  role = "${var.master_iam_role == "" ?
    join("|", aws_iam_role.master_role.*.name) :
    join("|", data.aws_iam_role.master_role.*.name)
  }"
}

data "aws_iam_role" "master_role" {
  count = "${var.master_iam_role == "" ? 0 : 1}"
  name  = "${var.master_iam_role}"
}

resource "aws_iam_role" "master_role" {
  count = "${var.master_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}-master-role"
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

resource "aws_iam_role_policy" "master_policy" {
  count = "${var.master_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}_master_policy"
  role  = "${aws_iam_role.master_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "ec2:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action": "iam:PassRole",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action" : [
        "s3:GetObject"
      ],
      "Resource": "arn:${local.arn}:s3:::*",
      "Effect": "Allow"
    },
    {
      "Action": "elasticloadbalancing:*",
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_instance" "master" {
  count = "${var.instance_count}"
  ami   = "${var.ec2_ami}"

  iam_instance_profile = "${aws_iam_instance_profile.master.name}"
  instance_type        = "${var.ec2_type}"
  subnet_id            = "${element(var.subnet_ids, count.index)}"
  user_data            = "${var.user_data_igns[count.index]}"

  vpc_security_group_ids      = ["${var.master_sg_ids}"]
  associate_public_ip_address = "${var.public_endpoints}"

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new CoreOS Release comes
    # out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}-master-${count.index}",
      "kubernetes.io/cluster/${var.cluster_name}", "owned",
      "tectonicClusterID", "${var.cluster_id}"
    ), var.extra_tags)}"

  root_block_device {
    volume_type = "${var.root_volume_type}"
    volume_size = "${var.root_volume_size}"
    iops        = "${var.root_volume_type == "io1" ? var.root_volume_iops : 0}"
  }

  volume_tags = "${merge(map(
    "Name", "${var.cluster_name}-master-${count.index}-vol",
    "kubernetes.io/cluster/${var.cluster_name}", "owned",
    "tectonicClusterID", "${var.cluster_id}"
  ), var.extra_tags)}"
}

resource "aws_elb_attachment" "masters_internal" {
  count    = "${var.private_endpoints ? var.instance_count : 0}"
  elb      = "${var.elb_api_internal_id}"
  instance = "${aws_instance.master.*.id[count.index]}"
}

resource "aws_elb_attachment" "masters_external" {
  count    = "${var.public_endpoints ? var.instance_count : 0}"
  elb      = "${var.elb_api_external_id}"
  instance = "${aws_instance.master.*.id[count.index]}"
}

resource "aws_elb_attachment" "masters_console" {
  count    = "${var.instance_count}"
  elb      = "${var.elb_console_id}"
  instance = "${aws_instance.master.*.id[count.index]}"
}
