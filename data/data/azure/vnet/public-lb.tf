locals {
  public_lb_frontend_ip_configuration_name    = "public-lb-ip"
  public_lb_frontend_ip_v6_configuration_name = "public-lb-ip-v6"
}

resource "azurerm_public_ip" "cluster_public_ip" {
  sku                 = "Standard"
  location            = var.region
  name                = "${var.cluster_id}-pip"
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  domain_name_label   = var.dns_label
}

data "azurerm_public_ip" "cluster_public_ip" {
  name                = azurerm_public_ip.cluster_public_ip.name
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
  name                = "${var.cluster_id}-public-lb"
  resource_group_name = var.resource_group_name
  location            = var.region

  dynamic "frontend_ip_configuration" {
    for_each = [for ip in [
      { name : local.public_lb_frontend_ip_configuration_name, value : azurerm_public_ip.cluster_public_ip.id, ipv6 : false },
      { name : local.public_lb_frontend_ip_v6_configuration_name, value : azurerm_public_ip.cluster_public_ip_v6[0].id, ipv6 : true },
      ] : {
      name : ip.name
      value : ip.value
      ipv6 : ip.ipv6
    } if ! ip.ipv6 || var.use_ipv6]

    content {
      name                          = frontend_ip_configuration.value.name
      public_ip_address_id          = frontend_ip_configuration.value.value
      private_ip_address_version    = frontend_ip_configuration.value.ipv6 ? "IPv6" : "IPv4"
      private_ip_address_allocation = "Dynamic"
    }
  }
}

resource "azurerm_lb_backend_address_pool" "master_public_lb_pool" {
  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.public.id
  name                = "${var.cluster_id}-public-lb-control-plane"
}

resource "azurerm_lb_backend_address_pool" "master_public_lb_pool_v6" {
  count = var.use_ipv6 ? 1 : 0

  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.public.id
  name                = "${var.cluster_id}-public-lb-control-plane-v6"
}

resource "azurerm_lb_rule" "public_lb_rule_api_internal" {
  count = var.private ? 0 : 1

  name                           = "api-internal"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.master_public_lb_pool.id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_configuration_name
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
  backend_address_pool_id        = azurerm_lb_backend_address_pool.master_public_lb_pool_v6[0].id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurerm_lb_rule" "internal_outbound_rule" {
  count = var.private ? 1 : 0

  name                           = "internal_outbound_rule"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.master_public_lb_pool.id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 27627
  backend_port                   = 27627
  frontend_ip_configuration_name = local.public_lb_frontend_ip_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
}

resource "azurerm_lb_rule" "internal_outbound_rule_v6" {
  count = var.private && var.use_ipv6 ? 1 : 0

  name                           = "internal_outbound_rule_v6"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.master_public_lb_pool_v6[0].id
  loadbalancer_id                = azurerm_lb.public.id
  frontend_port                  = 27627
  backend_port                   = 27627
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
}

resource "azurerm_lb_probe" "public_lb_probe_api_internal" {
  count = var.private ? 0 : 1

  name                = "api-internal-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 10
  number_of_probes    = 3
  loadbalancer_id     = azurerm_lb.public.id
  port                = 6443
  protocol            = "TCP"
}
