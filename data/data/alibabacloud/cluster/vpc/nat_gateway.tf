
resource "alicloud_nat_gateway" "nat_gateway" {
  vpc_id               = local.vpc_id
  nat_gateway_name     = "${local.prefix}-ngw"
  vswitch_id           = alicloud_vswitch.vswitch_nat_gateway.id
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
  count = length(local.vswitch_ids)

  depends_on        = [alicloud_eip_association.eip_association]
  snat_table_id     = alicloud_nat_gateway.nat_gateway.snat_table_ids
  source_vswitch_id = local.vswitch_ids[count.index]
  snat_ip           = alicloud_eip_address.eip.ip_address
}
