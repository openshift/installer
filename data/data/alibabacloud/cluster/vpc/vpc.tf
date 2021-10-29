
locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
  newbits     = tonumber(split("/", var.vpc_cidr_block)[1]) < 16 ? 20 - tonumber(split("/", var.vpc_cidr_block)[1]) : 4
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

  vswitch_name = "${local.prefix}-vswitch-${var.zone_ids[count.index]}"
  description  = local.description
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = cidrsubnet(var.vpc_cidr_block, local.newbits, count.index)
  zone_id      = var.zone_ids[count.index]
  tags = merge(
    {
      "Name" = "${local.prefix}-vswitch-${var.zone_ids[count.index]}"
    },
    var.tags,
  )
}

resource "alicloud_vswitch" "vswitch_nat_gateway" {
  vswitch_name = "${local.prefix}-vswitch-nat-gateway"
  description  = local.description
  vpc_id       = alicloud_vpc.vpc.id
  cidr_block   = cidrsubnet(var.vpc_cidr_block, local.newbits, local.newbits)
  zone_id      = var.nat_gateway_zone_id
  tags = merge(
    {
      "Name" = "${local.prefix}-vswitch-nat-gateway"
    },
    var.tags,
  )
}
