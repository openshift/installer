
resource "alicloud_nat_gateway" "nat_gateways" {
  count = length(var.zone_ids)

  vpc_id           = alicloud_vpc.vpc.id
  specification    = "Small"
  nat_gateway_name = "${local.prefix}-ngw-${count.index}"
  vswitch_id       = alicloud_vswitch.vswitchs[count.index].id
  nat_type         = "Enhanced"
  description      = local.description
  tags = merge(
    {
      "Name" = "${local.prefix}-ngw-${count.index}"
    },
    var.tags,
  )
}
