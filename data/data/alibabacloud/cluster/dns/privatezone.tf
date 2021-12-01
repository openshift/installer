locals {
  description  = "Created By OpenShift Installer"
  prefix       = var.cluster_id
  cluster_name = split(".", var.cluster_domain)[0]
}

// Using this data source can open Private Zone service automatically.
data "alicloud_pvtz_service" "open" {
  enable = "On"
}

data "alicloud_alidns_domains" "dns_public" {
  domain_name_regex = "^${var.base_domain}$"
}

resource "alicloud_alidns_record" "dns_public_record" {
  count = length(data.alicloud_alidns_domains.dns_public.domains) == 0 ? 0 : 1

  domain_name = var.base_domain
  rr          = "api.${local.cluster_name}"
  type        = "A"
  value       = var.slb_external_ip
  status      = "ENABLE"
}

resource "alicloud_pvtz_zone" "pvtz" {
  resource_group_id = var.resource_group_id
  zone_name         = var.cluster_domain
}

resource "alicloud_pvtz_zone_attachment" "pvtz_attachment" {
  zone_id = alicloud_pvtz_zone.pvtz.id
  vpc_ids = [var.vpc_id]
}

resource "alicloud_pvtz_zone_record" "pvtz_record_api_int" {
  zone_id = alicloud_pvtz_zone.pvtz.id
  type    = "A"
  rr      = "api-int"
  value   = var.slb_internal_ip
  ttl     = 60
}

resource "alicloud_pvtz_zone_record" "pvtz_record_api" {
  zone_id = alicloud_pvtz_zone.pvtz.id
  type    = "A"
  rr      = "api"
  value   = var.slb_internal_ip
  ttl     = 60
}
