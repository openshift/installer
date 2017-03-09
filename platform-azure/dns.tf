resource "azurerm_dns_zone" "tectonic_azure_dns_zone" {
   name = "${var.tectonic_base_domain}"
   resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
}

resource "azurerm_dns_a_record" "tectonic-api" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}-k8s"
  ttl     = "60"
  records = ["${azurerm_public_ip.tectonic_master_ip.ip_address}"]
}

resource "azurerm_dns_a_record" "tectonic-console" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}"
  ttl     = "60"
  records = ["${azurerm_public_ip.tectonic_master_ip.ip_address}"]
}

resource "azurerm_dns_a_record" "etcd" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}-etc"
  ttl     = "60"
  records = ["${azurerm_public_ip.etcd_node.ip_address}"]
}

resource "azurerm_dns_a_record" "master_nodes" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  count   = "${var.tectonic_master_count}"
  name    = "${var.tectonic_cluster_name}-master-${count.index}"
  ttl     = "59"
  records = ["${azurerm_public_ip.tectonic_master_ip.ip_address[count.index]}"]
}

resource "azurerm_dns_a_record" "worker_nodes" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  count   = "${var.tectonic_worker_count}"
  name    = "${var.tectonic_cluster_name}-worker-${count.index}"
  ttl     = "59"
  records = ["${azurerm_public_ip.worker_node.ip_address[count.index]}"]
}
