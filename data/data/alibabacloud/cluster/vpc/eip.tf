resource "alicloud_eip_address" "eip" {
  description          = local.description
  address_name         = "${local.prefix}-eip"
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  resource_group_id    = var.resource_group_id
  tags = merge(
    {
      "Name" = "${local.prefix}-eip"
    },
    var.tags,
  )
}

resource "alicloud_eip_association" "eip_association" {
  allocation_id = alicloud_eip_address.eip.id
  instance_id   = alicloud_nat_gateway.nat_gateway.id
  instance_type = "Nat"
}
