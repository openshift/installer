locals {
  description = "Created By OpenShift Installer"
  prefix      = var.cluster_id
}

resource "alicloud_pvtz_zone" "pvtz_internal" {
  resource_group_id = var.resource_group_id
  zone_name         = var.cluster_domain
}

resource "alicloud_pvtz_zone_attachment" "pvtz_internal_attachment" {
  zone_id = alicloud_pvtz_zone.pvtz_internal.id
  vpc_ids = [var.vpc_id]
}

resource "alicloud_pvtz_zone_record" "pvtz_internal_record_A" {
  zone_id = alicloud_pvtz_zone.pvtz_internal.id
  type    = "A"
  rr      = "api-int"
  value   = var.slb_internal_ip
  ttl     = 60
}

resource "alicloud_pvtz_zone" "pvtz_external" {
  resource_group_id = var.resource_group_id
  zone_name         = var.base_domain
}

resource "alicloud_pvtz_zone_attachment" "pvtz_external_attachment" {
  zone_id = alicloud_pvtz_zone.pvtz_external.id
  vpc_ids = [var.vpc_id]
}

resource "alicloud_pvtz_zone_record" "pvtz_external_record_A" {
  zone_id = alicloud_pvtz_zone.pvtz_external.id
  type    = "A"
  rr      = "api"
  value   = var.slb_external_ip
  ttl     = 60
}
