resource "azurerm_dns_a_record" "tectonic-api" {
  resource_group_name = "${var.resource_group_name}"
  zone_name           = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.cluster_name}-k8s"
  ttl     = "60"
  records = ["${var.master_ip_addresses}"]

  count = "${var.create_dns_zone ? 1 : 0}"
}

resource "azurerm_dns_a_record" "tectonic-console" {
  resource_group_name = "${var.resource_group_name}"
  zone_name           = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.cluster_name}"
  ttl     = "60"
  records = ["${var.console_ip_address}"]

  count = "${var.create_dns_zone ? 1 : 0}"
}

resource "azurerm_dns_a_record" "master_nodes" {
  resource_group_name = "${var.resource_group_name}"
  zone_name           = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.cluster_name}-master"
  ttl     = "59"
  records = ["${var.master_ip_addresses}"]

  count = "${var.create_dns_zone ? 1 : 0}"
}
