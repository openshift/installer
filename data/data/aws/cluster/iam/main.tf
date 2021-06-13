locals {
  arn = "aws"
}

data "aws_partition" "current" {}

resource "aws_iam_instance_profile" "worker" {
  name = "${var.cluster_id}-worker-profile"

  role = var.worker_iam_role_name != "" ? var.worker_iam_role_name : aws_iam_role.worker_role[0].name
}

resource "aws_iam_role" "worker_role" {
  count = var.worker_iam_role_name == "" ? 1 : 0
  name  = "${var.cluster_id}-worker-role"
  path  = "/"

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
      "Name" = "${var.cluster_id}-worker-role"
    },
    var.tags,
  )
}

resource "aws_iam_role_policy" "worker_policy" {

  // List curated from https://github.com/kubernetes/cloud-provider-aws#readme, minus entries specific to EKS
  // integrations.
  // This list should not be updated any further without an operator owning migrating changes here for existing
  // clusters.
  // Please see: docs/dev/aws/iam_permissions.md

  count = var.worker_iam_role_name == "" ? 1 : 0
  name = "${var.cluster_id}-worker-policy"
  role = aws_iam_role.worker_role[0].id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeRegions"
      ],
      "Resource": "*"
    }
  ]
}
EOF

}

