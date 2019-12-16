locals {
  // extracting "api.<clustername>" from <clusterdomain>
  api_external_name = "api.${replace(var.cluster_domain, ".${var.base_domain}", "")}"
}

resource "azureprivatedns_zone" "private" {
  name                = var.cluster_domain
  resource_group_name = var.resource_group_name
}

resource "azureprivatedns_zone_virtual_network_link" "network" {
  name                  = "${var.cluster_id}-network-link"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azureprivatedns_zone.private.name
  virtual_network_id    = var.virtual_network_id
}

resource "azureprivatedns_a_record" "apiint_internal" {
  name                = "api-int"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress]
}

resource "azureprivatedns_aaaa_record" "apiint_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "api-int"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v6]
}

resource "azureprivatedns_a_record" "api_internal" {
  name                = "api"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress]
}

resource "azureprivatedns_aaaa_record" "api_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "api"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v6]
}

resource "azurerm_dns_cname_record" "api_external" {
  count = var.private ? 0 : 1

  name                = local.api_external_name
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn
}

resource "azurerm_dns_cname_record" "api_external_v6" {
  count = var.private || ! var.use_ipv6 ? 0 : 1

  name                = "v6-${local.api_external_name}"
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn_v6
}

resource "azureprivatedns_a_record" "etcd_a_nodes" {
  count               = var.etcd_count
  name                = "etcd-${count.index}"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 60
  records             = [var.etcd_ip_addresses[count.index]]
}

resource "azureprivatedns_aaaa_record" "etcd_aaaa_nodes" {
  count               = var.use_ipv6 ? var.etcd_count : 0
  name                = "etcd-${count.index}"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 60
  records             = [var.etcd_ip_v6_addresses[count.index]]
}

resource "azureprivatedns_srv_record" "etcd_cluster" {
  name                = "_etcd-server-ssl._tcp"
  zone_name           = azureprivatedns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 60

  dynamic "record" {
    for_each = var.use_ipv6 ? concat(azureprivatedns_a_record.etcd_a_nodes.*.name, azureprivatedns_aaaa_record.etcd_aaaa_nodes.*.name) : azureprivatedns_a_record.etcd_a_nodes.*.name
    iterator = name
    content {
      target   = "${name.value}.${azureprivatedns_zone.private.name}"
      priority = 10
      weight   = 10
      port     = 2380
    }
  }
}

