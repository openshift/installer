locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
  tags = merge(
    {
      "GISV"                                      = "ocp",
      "sigs.k8s.io/cloud-provider-alibaba/origin" = "ocp",
      "kubernetes.io/cluster/${var.cluster_id}"   = "owned"
    },
    var.ali_extra_tags,
  )
  system_disk_size     = 120
  system_disk_category = "cloud_essd"
}

provider "alicloud" {
  access_key = var.ali_access_key
  secret_key = var.ali_secret_key
  region     = var.ali_region_id
}

data "alicloud_instances" "bootstrap_data" {
  ids = [alicloud_instance.bootstrap.id]
}

# Using this data source can enable OSS service automatically.
data "alicloud_oss_service" "open" {
  enable = "On"
}

resource "alicloud_oss_bucket" "bucket" {
  bucket = var.ali_ignition_bucket
  acl    = "private"
  tags = merge(
    {
      "Name" = "${local.prefix}-bucket"
    },
    local.tags,
  )
}

resource "alicloud_oss_bucket_object" "ignition_file" {
  bucket = alicloud_oss_bucket.bucket.id
  key    = "bootstrap.ign"
  source = var.ignition_bootstrap_file
  acl    = "private"
}

resource "alicloud_ram_role" "role" {
  name     = "${local.prefix}-role-bootstrap"
  document = <<EOF
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
  policy_name = "${local.prefix}-policy-bootstrap"
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
  resource_group_id = var.ali_resource_group_id
  name              = "${local.prefix}_sg_bootstrap"
  description       = local.description
  vpc_id            = var.vpc_id
  tags = merge(
    {
      "Name" = "${local.prefix}-sg-bootstrap"
    },
    local.tags,
  )
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
  resource_group_id = var.ali_resource_group_id

  host_name                  = "${local.prefix}-bootstrap"
  instance_name              = "${local.prefix}-bootstrap"
  instance_type              = var.ali_bootstrap_instance_type
  image_id                   = var.ali_image_id
  vswitch_id                 = var.vswitch_ids[0]
  security_groups            = [alicloud_security_group.sg_bootstrap.id, var.sg_master_id]
  internet_max_bandwidth_out = 5
  role_name                  = alicloud_ram_role.role.name

  system_disk_name        = "${local.prefix}_sys_disk-bootstrap"
  system_disk_description = local.description
  system_disk_category    = local.system_disk_category
  system_disk_size        = local.system_disk_size

  user_data = var.ali_bootstrap_stub_ignition
  tags = merge(
    {
      "Name" = "${local.prefix}-bootstrap"
    },
    local.tags,
  )
}

resource "alicloud_slb_backend_server" "slb_attachment_bootstraps" {
  count = length(var.slb_ids)

  load_balancer_id = var.slb_ids[count.index]
  backend_servers {
    server_id = alicloud_instance.bootstrap.id
    weight    = 90
  }
}