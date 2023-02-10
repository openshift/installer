locals {
  // extracting "api.<clustername>" from <clusterdomain>
  api_external_name = "api.${replace(var.cluster_domain, ".${var.base_domain}", "")}"
}

resource "azurerm_private_dns_zone" "private" {
  name                = var.cluster_domain
  resource_group_name = var.resource_group_name
  depends_on          = [azurerm_dns_cname_record.api_external_v4, azurerm_dns_cname_record.api_external_v6]
  tags                = var.azure_extra_tags
}

# Sleep injected due to https://github.com/hashicorp/terraform-provider-azurerm/issues/18350
resource "time_sleep" "wait_30_seconds" {
  depends_on      = [azurerm_private_dns_zone.private]
  create_duration = "30s"
}

resource "azurerm_private_dns_zone_virtual_network_link" "network" {
  name                  = "${var.cluster_id}-network-link"
  resource_group_name   = var.resource_group_name
  private_dns_zone_name = azurerm_private_dns_zone.private.name
  virtual_network_id    = var.virtual_network_id
  depends_on            = [time_sleep.wait_30_seconds]
  tags                  = var.azure_extra_tags
}

resource "azurerm_private_dns_a_record" "apiint_internal" {
  count = var.use_ipv4 ? 1 : 0

  name                = "api-int"
  zone_name           = azurerm_private_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v4]
  tags                = var.azure_extra_tags
}

resource "azurerm_private_dns_aaaa_record" "apiint_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "api-int"
  zone_name           = azurerm_private_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v6]
  tags                = var.azure_extra_tags
}

resource "azurerm_private_dns_a_record" "api_internal" {
  count = var.use_ipv4 ? 1 : 0

  name                = "api"
  zone_name           = azurerm_private_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v4]
  tags                = var.azure_extra_tags
}

resource "azurerm_private_dns_aaaa_record" "api_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "api"
  zone_name           = azurerm_private_dns_zone.private.name
  resource_group_name = var.resource_group_name
  ttl                 = 300
  records             = [var.internal_lb_ipaddress_v6]
  tags                = var.azure_extra_tags
}

resource "azurerm_dns_cname_record" "api_external_v4" {
  count = var.private || ! var.use_ipv4 ? 0 : 1

  name                = local.api_external_name
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn_v4
  tags                = var.azure_extra_tags
}

resource "azurerm_dns_cname_record" "api_external_v6" {
  count = var.private || ! var.use_ipv6 ? 0 : 1

  name                = "v6-${local.api_external_name}"
  zone_name           = var.base_domain
  resource_group_name = var.base_domain_resource_group_name
  ttl                 = 300
  record              = var.external_lb_fqdn_v6
  tags                = var.azure_extra_tags
}


