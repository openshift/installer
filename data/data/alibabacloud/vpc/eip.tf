resource "alicloud_eip" "eips" {
  count                = length(var.zone_ids)
  description          = local.description
  name                 = "${local.prefix}-eip-${count.index}"
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
  resource_group_id    = var.resource_group_id
  tags = merge(
    {
      "Name" = "${local.prefix}-eip-${count.index}"
    },
    var.tags,
  )
}

//Provides an Alicloud EIP Association resource for associating Elastic IP to Nat Gateway.
resource "alicloud_eip_association" "eip_associations" {
  count         = length(var.zone_ids)
  allocation_id = alicloud_eip.eips[count.index].id
  instance_id   = alicloud_nat_gateway.nat_gateways[count.index].id
  instance_type = "Nat"
}
