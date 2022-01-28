
resource "alicloud_nat_gateway" "nat_gateway" {
  count = length(var.vswitch_ids) == 0 ? 1 : 0

  vpc_id               = local.vpc_id
  nat_gateway_name     = "${local.prefix}-ngw"
  vswitch_id           = alicloud_vswitch.vswitch_nat_gateway[0].id
  internet_charge_type = "PayByLcu"
  nat_type             = "Enhanced"
  description          = local.description
  tags = merge(
    {
      "Name" = "${local.prefix}-ngw"
    },
    var.tags,
  )
}

resource "alicloud_snat_entry" "snat_entrys" {
  count = length(var.vswitch_ids) == 0 ? length(local.vswitch_ids) : 0

  depends_on        = [alicloud_eip_association.eip_association]
  snat_table_id     = alicloud_nat_gateway.nat_gateway[0].snat_table_ids
  source_vswitch_id = local.vswitch_ids[count.index]
  snat_ip           = alicloud_eip_address.eip[0].ip_address
}
