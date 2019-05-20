locals {
  // extracting "api.<clustername>" from <clusterdomain>
  api_external_name = "api.${replace(var.cluster_domain, ".${var.base_domain}", "")}"
}

resource "azurerm_dns_zone" "private" {
  name                           = var.cluster_domain
  resource_group_name            = var.resource_group_name
  zone_type                      = "Private"
  resolution_virtual_network_ids = [var.internal_dns_resolution_vnet_id]
}

resource "azurerm_dns_cname_record" "apiint_internal" {
  name                = "api-int"
  zone_name           = azurerm_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn
}

resource "azurerm_dns_cname_record" "api_internal" {
  name                = "api"
  zone_name           = azurerm_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn
}

resource "azurerm_dns_cname_record" "api_external" {
  name                = local.api_external_name
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn
}

resource "azurerm_dns_a_record" "etcd_a_nodes" {
  count               = var.etcd_count
  name                = "etcd-${count.index}"
  zone_name           = azurerm_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 60
  records             = [var.etcd_ip_addresses[count.index]]
}

resource "azurerm_dns_srv_record" "etcd_cluster" {
  name                = "_etcd-server-ssl._tcp"
  zone_name           = azurerm_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 60

  dynamic "record" {
    for_each = azurerm_dns_a_record.etcd_a_nodes.*.name
    iterator = name
    content {
      target   = "${name.value}.${azurerm_dns_zone.private.name}"
      priority = 10
      weight   = 10
      port     = 2380
    }
  }
}

