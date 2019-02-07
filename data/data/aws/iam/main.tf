locals {
  arn = "aws"
}

resource "aws_iam_instance_profile" "compute" {
  name = "${var.cluster_name}-compute-profile"

  role = "${aws_iam_role.compute_role.name}"
}

resource "aws_iam_role" "compute_role" {
  name = "${var.cluster_name}-compute-role"
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

resource "aws_iam_role_policy" "compute_policy" {
  name = "${var.cluster_name}_compute_policy"
  role = "${aws_iam_role.compute_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "ec2:Describe*",
      "Resource": "*"
    }
  ]
}
EOF
}
