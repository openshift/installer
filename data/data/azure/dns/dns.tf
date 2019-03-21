resource "azurerm_dns_zone" "private" {
  name                = "${var.base_domain}"
  resource_group_name = "${var.resource_group_name}"
  zone_type           = "Private"
  resolution_virtual_network_ids = ["${var.internal_dns_resolution_vnet_id}"]
}

resource "azurerm_dns_zone" "public" {
  name                = "${var.cluster_domain}"
  resource_group_name = "${var.resource_group_name}"
  zone_type           = "Public"
}

resource "azurerm_dns_a_record" "api_internal" {
  name                = "api.${var.cluster_name}"
  zone_name           = "${azurerm_dns_zone.private.name}"
  resource_group_name = "${var.resource_group_name}"
  ttl                 = 300
  records             = ["${var.internal_lb_ipaddress}"]
}

resource "azurerm_dns_cname_record" "api_external" {
  name                = "api"
  zone_name           = "${azurerm_dns_zone.public.name}"
  resource_group_name = "${var.resource_group_name}"
  ttl                 = 300
  record              = "${var.external_lb_dns_label}.${var.base_domain}"
}