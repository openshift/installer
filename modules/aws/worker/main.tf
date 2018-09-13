locals {
  arn = "aws"
}

resource "aws_iam_instance_profile" "worker" {
  name = "${var.cluster_name}-worker-profile"

  role = "${var.worker_iam_role == "" ?
    join("|", aws_iam_role.worker_role.*.name) :
    join("|", data.aws_iam_role.worker_role.*.name)
  }"
}

data "aws_iam_role" "worker_role" {
  count = "${var.worker_iam_role == "" ? 0 : 1}"
  name  = "${var.worker_iam_role}"
}

resource "aws_iam_role" "worker_role" {
  count = "${var.worker_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}-worker-role"
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

resource "aws_iam_role_policy" "worker_policy" {
  count = "${var.worker_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}_worker_policy"
  role  = "${aws_iam_role.worker_role.id}"

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
      "Action": "elasticloadbalancing:*",
      "Resource": "*",
      "Effect": "Allow"
    },
    {
      "Action" : [
        "s3:GetObject"
      ],
      "Resource": "arn:${local.arn}:s3:::*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_instance" "worker" {
  count = "${var.instance_count}"
  ami   = "${var.ec2_ami}"

  iam_instance_profile   = "${aws_iam_instance_profile.worker.name}"
  instance_type          = "${var.ec2_type}"
  subnet_id              = "${element(var.subnet_ids, count.index)}"
  user_data              = "${var.user_data_ign}"
  vpc_security_group_ids = ["${var.sg_ids}"]

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new CoreOS Release comes
    # out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}-worker-${count.index}",
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

resource "aws_elb_attachment" "workers" {
  count    = "${length(var.load_balancers) * var.instance_count}"
  elb      = "${var.load_balancers[count.index / var.instance_count]}"
  instance = "${aws_instance.worker.*.id[count.index % var.instance_count]}"
}
