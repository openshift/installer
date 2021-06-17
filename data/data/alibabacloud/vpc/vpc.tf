
locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

resource "alicloud_vpc" "vpc" {
  resource_group_id = var.resource_group_id
  vpc_name          = "${local.prefix}-vpc"
  cidr_block        = var.vpc_cidr_block
  description       = local.description
  tags = merge(
    {
      "Name" = "${local.prefix}-vpc"
    },
    var.tags,
  )
}

resource "alicloud_vswitch" "vswitchs" {
  count = length(var.zone_ids)

  vswitch_name = "${local.prefix}-vswitch-${count.index}"
  description  = local.description
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = var.vswitch_cidr_blocks[count.index]
  zone_id      = var.zone_ids[count.index]
  tags = merge(
    {
      "Name" = "${local.prefix}-vswitch"
    },
    var.tags,
  )
}
