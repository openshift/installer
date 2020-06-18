locals {
  public_lb_frontend_ip_v4_configuration_name = "public-lb-ip-v4"
  public_lb_frontend_ip_v6_configuration_name = "public-lb-ip-v6"
}

resource "azurerm_public_ip" "cluster_public_ip_v4" {
  // DEBUG: Azure apparently requires dual stack LB for v6
  count = var.use_ipv4 || true ? 1 : 0

  sku                 = "Standard"
  location            = var.region
  name                = "${var.cluster_id}-pip-v4"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  domain_name_label   = var.dns_label
}

data "azurerm_public_ip" "cluster_public_ip_v4" {
  // DEBUG: Azure apparently requires dual stack LB for v6
  count = var.use_ipv4 || true ? 1 : 0

  name                = azurerm_public_ip.cluster_public_ip_v4[0].name
  resource_group_name = var.resource_group_name
}

resource "azurerm_public_ip" "cluster_public_ip_v6" {
  count = var.use_ipv6 ? 1 : 0

  ip_version          = "IPv6"
  sku                 = "Standard"
  location            = var.region
  name                = "${var.cluster_id}-pip-v6"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  domain_name_label   = var.dns_label
}

data "azurerm_public_ip" "cluster_public_ip_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = azurerm_public_ip.cluster_public_ip_v6[0].name
  resource_group_name = var.resource_group_name
}

resource "azurerm_lb" "public" {
  sku                 = "Standard"
  name                = var.cluster_id
  resource_group_name = var.resource_group_name
  location            = var.region

  dynamic "frontend_ip_configuration" {
    for_each = [for ip in [
      // DEBUG: Azure apparently requires dual stack LB for external load balancers v6
      { name : local.public_lb_frontend_ip_v4_configuration_name, value : azurerm_public_ip.cluster_public_ip_v4[0].id, include : true, ipv6 : false },
      { name : local.public_lb_frontend_ip_v6_configuration_name, value : var.use_ipv6 ? azurerm_public_ip.cluster_public_ip_v6[0].id : null, include : var.use_ipv6, ipv6 : true },
      ] : {
      name : ip.name
      value : ip.value
      ipv6 : ip.ipv6
      include : ip.include
      } if ip.include
    ]

    content {
      name                          = frontend_ip_configuration.value.name
      public_ip_address_id          = frontend_ip_configuration.value.value
      private_ip_address_version    = frontend_ip_configuration.value.ipv6 ? "IPv6" : "IPv4"
      private_ip_address_allocation = "Dynamic"
    }
  }
}

resource "azurerm_lb_backend_address_pool" "public_lb_pool_v4" {
  count = var.use_ipv4 ? 1 : 0

  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.public.id
  name                = var.cluster_id
}

resource "azurerm_lb_backend_address_pool" "public_lb_pool_v6" {
  count = var.use_ipv6 ? 1 : 0

  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.public.id
  name                = "${var.cluster_id}-IPv6"
}

resource "azurerm_lb_rule" "public_lb_rule_api_internal_v4" {
  count = var.private || ! var.use_ipv4 ? 0 : 1

  name                           = "api-internal-v4"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v4_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurerm_lb_rule" "public_lb_rule_api_internal_v6" {
  count = var.private || ! var.use_ipv6 ? 0 : 1

  name                           = "api-internal-v6"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurerm_lb_outbound_rule" "public_lb_outbound_rule_v4" {
  count = var.private && var.use_ipv4 ? 1 : 0

  name                    = "outbound-rule-v4"
  resource_group_name     = var.resource_group_name
  loadbalancer_id         = azurerm_lb.public.id
  backend_address_pool_id = azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id
  protocol                = "All"

  frontend_ip_configuration {
    name = local.public_lb_frontend_ip_v4_configuration_name
  }
}

resource "azurerm_lb_outbound_rule" "public_lb_outbound_rule_v6" {
  count = var.private && var.use_ipv6 ? 1 : 0

  name                    = "outbound-rule-v6"
  resource_group_name     = var.resource_group_name
  loadbalancer_id         = azurerm_lb.public.id
  backend_address_pool_id = azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id
  protocol                = "All"

  frontend_ip_configuration {
    name = local.public_lb_frontend_ip_v6_configuration_name
  }
}

resource "azurerm_lb_probe" "public_lb_probe_api_internal" {
  count = var.private ? 0 : 1

  name                = "api-internal-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurerm_lb.public.id
  port                = 6443
  protocol            = "HTTPS"
  request_path        = "/readyz"
}
