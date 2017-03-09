resource "azurerm_dns_a_record" "etcd" {
  resource_group_name = "${azurerm_resource_group.tectonic_azure_dns_resource_group.name}"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}-etc"
  ttl     = "60"
  records = ["${azurerm_public_ip.etcd_node.ip_address}"]
}
