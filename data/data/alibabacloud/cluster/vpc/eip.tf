resource "alicloud_eip_address" "eip" {
  count = length(var.vswitch_ids) == 0 ? 1 : 0

  description          = local.description
  address_name         = "${local.prefix}-eip"
  payment_type         = "PayAsYouGo"
  internet_charge_type = "PayByTraffic"
  resource_group_id    = var.resource_group_id
  tags = merge(
    {
      "Name" = "${local.prefix}-eip"
    },
    var.tags,
  )
}

resource "alicloud_eip_association" "eip_association" {
  count = length(var.vswitch_ids) == 0 ? 1 : 0

  allocation_id = alicloud_eip_address.eip[0].id
  instance_id   = alicloud_nat_gateway.nat_gateway[0].id
  instance_type = "Nat"
}
