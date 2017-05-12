resource "azurerm_dns_a_record" "etcd" {
  resource_group_name = "${var.resource_group_name}"
  zone_name           = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.cluster_name}-etcd"
  ttl     = "60"
  records = ["${var.etcd_ip_addresses}"]

  count = "${var.use_custom_fqdn ? 1 : 0}"
}
