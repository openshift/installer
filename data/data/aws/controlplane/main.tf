locals {
  arn = "aws"
}

resource "aws_iam_instance_profile" "control_plane" {
  name = "${var.cluster_name}-control-plane-profile"

  role = "${aws_iam_role.control_plane_role.name}"
}

resource "aws_iam_role" "control_plane_role" {
  name = "${var.cluster_name}-control-plane-role"
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

  tags = "${var.tags}"
}

resource "aws_iam_role_policy" "control_plane_policy" {
  name = "${var.cluster_name}_control_plane_policy"
  role = "${aws_iam_role.control_plane_role.id}"

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

resource "aws_instance" "control_plane" {
  count = "${var.instance_count}"
  ami   = "${var.ec2_ami}"

  iam_instance_profile = "${aws_iam_instance_profile.control_plane.name}"
  instance_type        = "${var.ec2_type}"
  subnet_id            = "${element(var.subnet_ids, count.index)}"
  user_data            = "${var.user_data_ign}"

  vpc_security_group_ids = ["${var.control_plane_sg_ids}"]

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new CoreOS Release comes
    # out.
    ignore_changes = ["ami"]
  }

  tags = "${merge(map(
      "Name", "${var.cluster_name}-${var.machine_pool_name}-${count.index}",
      "clusterid", "${var.cluster_name}"
    ), var.tags)}"

  root_block_device {
    volume_type = "${var.root_volume_type}"
    volume_size = "${var.root_volume_size}"
    iops        = "${var.root_volume_type == "io1" ? var.root_volume_iops : 0}"
  }

  volume_tags = "${merge(map(
    "Name", "${var.cluster_name}-${var.machine_pool_name}-${count.index}-vol",
  ), var.tags)}"
}

resource "aws_lb_target_group_attachment" "control_plane" {
  count = "${var.instance_count * var.target_group_arns_length}"

  target_group_arn = "${var.target_group_arns[count.index % var.target_group_arns_length]}"
  target_id        = "${aws_instance.control_plane.*.private_ip[count.index / var.target_group_arns_length]}"
}
