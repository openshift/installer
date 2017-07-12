resource "azurerm_dns_a_record" "tectonic-api" {
  count = "${var.base_domain != "" ? 1 : 0}"

  resource_group_name = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$1")}"
  zone_name           = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$2")}"

  name    = "${var.cluster_name}-api"
  ttl     = "60"
  records = ["${var.api_ip_addresses}"]
}

resource "azurerm_dns_a_record" "tectonic-console" {
  count = "${var.base_domain != "" ? 1 : 0}"

  resource_group_name = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$1")}"
  zone_name           = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$2")}"

  name    = "${var.cluster_name}"
  ttl     = "60"
  records = ["${var.console_ip_addresses}"]
}

resource "azurerm_dns_a_record" "tectonic-etcd" {
  count = "${var.base_domain != "" ? var.etcd_count : 0}"

  resource_group_name = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$1")}"
  zone_name           = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$2")}"

  name    = "${var.cluster_name}-etcd-${count.index}"
  ttl     = "60"
  records = ["${var.etcd_ip_addresses[count.index]}"]
}

resource "azurerm_dns_a_record" "tectonic-master" {
  count = "${var.base_domain != "" ? var.master_count : 0}"

  resource_group_name = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$1")}"
  zone_name           = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$2")}"

  name    = "${var.cluster_name}-master-${count.index}"
  ttl     = "60"
  records = ["${var.master_ip_addresses[count.index]}"]
}

resource "azurerm_dns_a_record" "tectonic-worker" {
  count = "${var.base_domain != "" ? var.worker_count : 0}"

  resource_group_name = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$1")}"
  zone_name           = "${replace(var.external_dns_zone_id, "${var.const_id_to_group_name_regex}", "$2")}"

  name    = "${var.cluster_name}-worker-${count.index}"
  ttl     = "60"
  records = ["${var.worker_ip_addresses[count.index]}"]
}
