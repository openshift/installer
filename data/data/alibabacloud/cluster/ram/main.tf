locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

resource "alicloud_ram_role" "role_master" {
  name        = "${local.prefix}-role-master"
  description = local.description
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_ram_policy" "role_policy_master" {
  policy_name = "${local.prefix}-policy-master"
  policy_document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": [
        "ecs:Describe*",
        "ecs:AttachDisk",
        "ecs:CreateDisk",
        "ecs:ResizeDisk",
        "ecs:CreateSnapshot",
        "ecs:CreateRouteEntry",
        "ecs:DeleteDisk",
        "ecs:DeleteSnapshot",
        "ecs:DeleteRouteEntry",
        "ecs:DetachDisk",
        "ecs:ModifyAutoSnapshotPolicyEx",
        "ecs:ModifyDiskAttribute",
        "ecs:CreateNetworkInterface",
        "ecs:DescribeNetworkInterfaces",
        "ecs:AttachNetworkInterface",
        "ecs:DetachNetworkInterface",
        "ecs:DeleteNetworkInterface",
        "ecs:DescribeInstanceAttribute",
        "ecs:AssignPrivateIpAddresses",
        "ecs:UnassignPrivateIpAddresses",
        "ecs:DescribeInstances"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "pvtz:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "cr:GetAuthorizationToken",
        "cr:ListInstanceEndpoint",
        "cr:PullRepository",
        "cr:GetInstanceVpcEndpoint"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "slb:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "cms:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "vpc:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "log:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    }
  ]
}
  EOF
}

resource "alicloud_ram_role_policy_attachment" "attach_policy_master" {
  policy_name = alicloud_ram_policy.role_policy_master.name
  policy_type = alicloud_ram_policy.role_policy_master.type
  role_name   = alicloud_ram_role.role_master.name
}

resource "alicloud_ram_role" "role_worker" {
  name        = "${local.prefix}-role-worker"
  description = local.description
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  EOF
}

resource "alicloud_ram_policy" "role_policy_worker" {
  policy_name = "${local.prefix}-policy-worker"
  policy_document = <<EOF
{
  "Version": "1",
  "Statement": [
    {
      "Action": [
        "ecs:AttachDisk",
        "ecs:DetachDisk",
        "ecs:DescribeDisks",
        "ecs:CreateDisk",
        "ecs:CreateSnapshot",
        "ecs:DeleteDisk",
        "ecs:AddTags",
        "ecs:RemoveTags",
        "ecs:CreateNetworkInterface",
        "ecs:DescribeNetworkInterfaces",
        "ecs:AttachNetworkInterface",
        "ecs:DetachNetworkInterface",
        "ecs:DeleteNetworkInterface",
        "ecs:DescribeInstanceAttribute",
        "ecs:AssignPrivateIpAddresses",
        "ecs:UnassignPrivateIpAddresses",
        "ecs:DescribeInstances"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "pvtz:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "cr:GetAuthorizationToken",
        "cr:ListInstanceEndpoint",
        "cr:PullRepository",
        "cr:GetInstanceVpcEndpoint"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "log:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "cms:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "vpc:*"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    }
  ]
}
  EOF
}

resource "alicloud_ram_role_policy_attachment" "attach_policy_worker" {
  policy_name = alicloud_ram_policy.role_policy_worker.name
  policy_type = alicloud_ram_policy.role_policy_worker.type
  role_name   = alicloud_ram_role.role_worker.name
}
