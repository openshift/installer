locals {
  // Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot use a dynamic list for count
  // and therefore are force to implicitly assume that the list is of aws_lb_target_group_arns_length - 1, in case there is no api_external
  target_group_arns_length = var.publish_strategy == "External" ? var.target_group_arns_length : var.target_group_arns_length - 1
  description              = "Created By OpenShift Installer"
}

data "aws_partition" "current" {}

data "aws_ebs_default_kms_key" "current" {}

resource "aws_iam_instance_profile" "master" {
  name = "${var.cluster_id}-master-profile"

  role = var.iam_role_name != "" ? var.iam_role_name : aws_iam_role.master_role[0].name

  tags = merge(
    {
      "Name" = "${var.cluster_id}-master-profile"
    },
    var.tags,
  )
}

resource "aws_iam_role" "master_role" {
  count       = var.iam_role_name == "" ? 1 : 0
  name        = "${var.cluster_id}-master-role"
  path        = "/"
  description = local.description

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
      "Name" = "${var.cluster_id}-master-role"
    },
    var.tags,
  )
}

resource "aws_iam_role_policy" "master_policy" {

  // List curated from https://github.com/kubernetes/cloud-provider-aws#readme, minus entries specific to EKS
  // integrations.
  // This list should not be updated any further without an operator owning migrating changes here for existing
  // clusters.
  // Please see: docs/dev/aws/iam_permissions.md

  count = var.iam_role_name == "" ? 1 : 0
  name = "${var.cluster_id}-master-policy"
  role = aws_iam_role.master_role[0].id

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:AttachVolume",
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:CreateSecurityGroup",
        "ec2:CreateTags",
        "ec2:CreateVolume",
        "ec2:DeleteSecurityGroup",
        "ec2:DeleteVolume",
        "ec2:Describe*",
        "ec2:DetachVolume",
        "ec2:ModifyInstanceAttribute",
        "ec2:ModifyVolume",
        "ec2:RevokeSecurityGroupIngress",
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:AttachLoadBalancerToSubnets",
        "elasticloadbalancing:ApplySecurityGroupsToLoadBalancer",
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateLoadBalancerPolicy",
        "elasticloadbalancing:CreateLoadBalancerListeners",
        "elasticloadbalancing:CreateTargetGroup",
        "elasticloadbalancing:ConfigureHealthCheck",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:DeleteLoadBalancerListeners",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:DeregisterInstancesFromLoadBalancer",
        "elasticloadbalancing:DeregisterTargets",
        "elasticloadbalancing:Describe*",
        "elasticloadbalancing:DetachLoadBalancerFromSubnets",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:ModifyTargetGroup",
        "elasticloadbalancing:ModifyTargetGroupAttributes",
        "elasticloadbalancing:RegisterInstancesWithLoadBalancer",
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:SetLoadBalancerPoliciesForBackendServer",
        "elasticloadbalancing:SetLoadBalancerPoliciesOfListener",
        "kms:DescribeKey"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF

}

resource "aws_network_interface" "master" {
  count       = var.instance_count
  description = local.description
  subnet_id   = var.az_to_subnet_id[var.availability_zones[count.index]]

  security_groups = var.master_sg_ids

  tags = merge(
    {
      "Name" = "${var.cluster_id}-master-${count.index}"
    },
    var.tags,
  )
}

resource "aws_instance" "master" {
  count = var.instance_count
  ami   = var.ec2_ami

  iam_instance_profile = aws_iam_instance_profile.master.name
  instance_type        = var.instance_type
  user_data            = var.user_data_ign

  network_interface {
    network_interface_id = aws_network_interface.master[count.index].id
    device_index         = 0
  }

  dynamic "instance_market_options" {
    for_each = var.use_spot_instance ? [1] : []
    content {
      market_type = "spot"
      spot_options {
        spot_instance_type = "one-time"
      }
    }
  }

  lifecycle {
    # Ignore changes in the AMI which force recreation of the resource. This
    # avoids accidental deletion of nodes whenever a new CoreOS Release comes
    # out.
    ignore_changes = [ami]
  }

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = var.instance_metadata_authentication
  }

  tags = merge(
    {
      "Name" = "${var.cluster_id}-master-${count.index}"
    },
    var.tags,
  )

  root_block_device {
    volume_type = var.root_volume_type
    volume_size = var.root_volume_size
    iops        = var.root_volume_iops
    encrypted   = var.root_volume_encrypted
    kms_key_id  = var.root_volume_kms_key_id == "" ? data.aws_ebs_default_kms_key.current.key_arn : var.root_volume_kms_key_id
  }

  volume_tags = merge(
    {
      "Name" = "${var.cluster_id}-master-${count.index}-vol"
    },
    var.tags,
  )

  depends_on = [
    aws_network_interface.master
  ]
}

resource "aws_lb_target_group_attachment" "master" {
  count = var.instance_count * local.target_group_arns_length

  target_group_arn = var.target_group_arns[count.index % local.target_group_arns_length]
  target_id        = aws_instance.master[floor(count.index / local.target_group_arns_length)].private_ip
}

