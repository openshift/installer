locals {
  // extracting "api.<clustername>" from <clusterdomain>
  api_external_name = "api.${replace(var.cluster_domain, ".${var.base_domain}", "")}"
}

resource "azurerm_private_dns_zone_virtual_network_link" "network" {
  count = var.azure_supports_private_dns ? 0 : 1

  name                  = "${var.cluster_id}-network-link"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = var.azure_cluster_dns_zone_name
  virtual_network_id    = var.virtual_network_id
}

resource "azurerm_dns_a_record" "apiint_internal" {
  name                = "api-int"
  zone_name           = var.azure_cluster_dns_zone_name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress]
}

resource "azurerm_dns_a_record" "api_internal" {
  name                = "api"
  zone_name           = var.azure_cluster_dns_zone_name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress]
}

resource "azurerm_dns_cname_record" "api_external" {
  count = var.private ? 0 : 1

  name                = local.api_external_name
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn
}

resource "azurerm_dns_a_record" "etcd_a_nodes" {

  count               = var.etcd_count
  name                = "etcd-${count.index}"
  zone_name           = var.azure_cluster_dns_zone_name
  resource_group_name = var.resource_group_name
  ttl                 = 60
  records             = [var.etcd_ip_addresses[count.index]]
}

resource "azurerm_dns_srv_record" "etcd_cluster" {
  name                = "_etcd-server-ssl._tcp"
  zone_name           = var.azure_cluster_dns_zone_name
  resource_group_name = var.resource_group_name
  ttl                 = 60

  dynamic "record" {
    for_each = azurerm_dns_a_record.etcd_a_nodes.*.name
    iterator = name
    content {
      target   = "${name.value}.${var.azure_cluster_dns_zone_name}"
      priority = 10
      weight   = 10
      port     = 2380
    }
  }
}

