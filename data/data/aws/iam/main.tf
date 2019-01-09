locals {
  arn = "aws"
}

resource "aws_iam_instance_profile" "compute" {
  name = "${var.cluster_name}-compute-profile"

  role = "${var.compute_iam_role == "" ?
    join("|", aws_iam_role.compute_role.*.name) :
    join("|", data.aws_iam_role.compute_role.*.name)
  }"
}

data "aws_iam_role" "compute_role" {
  count = "${var.compute_iam_role == "" ? 0 : 1}"
  name  = "${var.compute_iam_role}"
}

resource "aws_iam_role" "compute_role" {
  count = "${var.compute_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}-compute-role"
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

resource "aws_iam_role_policy" "compute_policy" {
  count = "${var.compute_iam_role == "" ? 1 : 0}"
  name  = "${var.cluster_name}_compute_policy"
  role  = "${aws_iam_role.compute_role.id}"

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
