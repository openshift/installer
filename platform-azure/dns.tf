resource "azurerm_dns_zone" "tectonic_azure_dns_zone" {
   name = ""${var.tectonic_base_domain}"
   resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_dns_a_record" "tectonic-api" {
  resource_group_name = "tectonic_azure_cluster_resource_group"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}-k8s"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4}"]
}

resource "azurerm_dns_a_recard" "tectonic-console" {
  resource_group_name = "tectonic_azure_cluster_resource_group"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4}"]
}

resource "azurerm_dns_a_recard" "etcd" {
  resource_group_name = "tectonic_azure_cluster_resource_group"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  name    = "${var.tectonic_cluster_name}-etc"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.etcd_node.*.access_ip_v4}"]
}

resource "azurerm_dns_a_recard" "master_nodes" {
  resource_group_name = "tectonic_azure_cluster_resource_group"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  count   = "${var.tectonic_master_count}"
  name    = "${var.tectonic_cluster_name}-master-${count.index}"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.master_node.*.access_ip_v4[count.index]}"]
}

resource "azurerm_dns_a_recard" "worker_nodes" {
  resource_group_name = "tectonic_azure_cluster_resource_group"
  zone_name = "${azurerm_dns_zone.tectonic_azure_dns_zone.name}"

  count   = "${var.tectonic_worker_count}"
  name    = "${var.tectonic_cluster_name}-worker-${count.index}"
  ttl     = "60"
  records = ["${openstack_compute_instance_v2.worker_node.*.access_ip_v4[count.index]}"]
}
