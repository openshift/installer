locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = var.ignition_bucket
  acl    = "private"
  tags = merge(
    {
      "Name" = "${local.prefix}-bucket"
    },
    var.tags,
  )
}

resource "alicloud_oss_bucket_object" "ignition_file" {
  bucket = alicloud_oss_bucket.bucket.id
  key    = "bootstrap.ign"
  source = var.ignition
  acl    = "private"
}

resource "alicloud_ram_role" "role" {
  name        = "${local.prefix}-role-bootstrap"
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
  description = local.description
}

resource "alicloud_ram_policy" "role_policy" {
  policy_name     = "${local.prefix}-bootstrap-policy"
  policy_document = <<EOF
  {
    "Statement": [
      {
        "Action": [
          "ecs:Describe*",
          "ecs:AttachDisk",
          "ecs:DetachDisk"
        ],
        "Effect": "Allow",
        "Resource": [
          "*"
        ]
      }
    ],
      "Version": "1"
  }
  EOF
}

resource "alicloud_ram_role_policy_attachment" "attach" {
  policy_name = alicloud_ram_policy.role_policy.name
  policy_type = alicloud_ram_policy.role_policy.type
  role_name   = alicloud_ram_role.role.name
}

resource "alicloud_security_group" "sg_bootstrap" {
  resource_group_id = var.resource_group_id
  name              = "${local.prefix}_sg_bootstrap"
  description       = local.description
  vpc_id            = var.vpc_id
}

resource "alicloud_security_group_rule" "sg_rule_ssh" {
  description       = local.description
  security_group_id = alicloud_security_group.sg_bootstrap.id
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "sg_rule_journald_gateway" {
  description       = local.description
  security_group_id = alicloud_security_group.sg_bootstrap.id
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "19531/19531"
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_instance" "bootstrap" {
  resource_group_id = var.resource_group_id

  instance_name              = "${local.prefix}_bootstrap"
  instance_type              = var.instance_type
  image_id                   = var.image_id
  internet_max_bandwidth_out = 10

  vswitch_id      = var.vswitch_id
  security_groups = [alicloud_security_group.sg_bootstrap.id]
  role_name       = alicloud_ram_role.role.name

  system_disk_name        = "${local.prefix}_sys_disk-bootstrap"
  system_disk_description = local.description
  system_disk_category    = var.system_disk_category
  system_disk_size        = var.system_disk_size

  data_disks {
    name        = "${local.prefix}_data_disk-bootstrap"
    category    = var.data_disk_category
    size        = var.data_disk_size
    description = local.description
  }

  user_data = var.ignition_stub
  key_name  = var.key_name
  tags = merge(
    {
      "Name" = "${local.prefix}-bootstrap"
    },
    var.tags,
  )
}

resource "alicloud_slb_attachment" "slb_attachment_bootstrap" {
  load_balancer_id = var.slb_id
  instance_ids     = [alicloud_instance.bootstrap.id]
  weight           = 90
}
