
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

resource "alicloud_snat_entry" "snat_entry" {
  count = length(var.zone_ids)

  depends_on        = [alicloud_eip_association.eip_associations]
  snat_table_id     = alicloud_nat_gateway.nat_gateways[count.index].snat_table_ids
  source_vswitch_id = alicloud_vswitch.vswitchs[count.index].id
  snat_ip           = alicloud_eip.eips[count.index].ip_address
}