
resource "alicloud_nat_gateway" "nat_gateway" {
  vpc_id           = alicloud_vpc.vpc.id
  specification    = "Small"
  nat_gateway_name = "${local.prefix}-ngw"
  vswitch_id       = alicloud_vswitch.vswitch_nat_gateway.id
  nat_type         = "Enhanced"
  description      = local.description
  tags = merge(
    {
      "Name" = "${local.prefix}-ngw"
    },
    var.tags,
  )
}

resource "alicloud_snat_entry" "snat_entrys" {
  count = length(var.zone_ids)

  depends_on        = [alicloud_eip_association.eip_association]
  snat_table_id     = alicloud_nat_gateway.nat_gateway.snat_table_ids
  source_vswitch_id = alicloud_vswitch.vswitchs[count.index].id
  snat_ip           = alicloud_eip_address.eip.ip_address
}
