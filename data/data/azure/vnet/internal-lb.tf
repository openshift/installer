locals {
  internal_lb_frontend_ip_configuration_name    = "internal-lb-ip"
  internal_lb_frontend_ip_v6_configuration_name = "internal-lb-ip-v6"
}

resource "azurerm_lb" "internal" {
  sku                 = "Standard"
  name                = "${var.cluster_id}-internal-lb"
  resource_group_name = var.resource_group_name
  location            = var.region

  dynamic "frontend_ip_configuration" {
    for_each = [for ip in [
      { name : local.internal_lb_frontend_ip_configuration_name, ipv6 : false },
      { name : local.internal_lb_frontend_ip_v6_configuration_name, ipv6 : true },
      ] : {
      name : ip.name
      ipv6 : ip.ipv6
    } if ! ip.ipv6 || var.use_ipv6]

    content {
      name                       = frontend_ip_configuration.value.name
      subnet_id                  = local.master_subnet_id
      private_ip_address_version = frontend_ip_configuration.value.ipv6 ? "IPv6" : "IPv4"
      # WORKAROUND: Allocate a high ipv6 internal LB address to avoid the race with NIC allocation (a master and the LB
      #   were being assigned the same IP dynamically). Issue is being tracked as a support ticket to Azure.
      private_ip_address_allocation = frontend_ip_configuration.value.ipv6 ? "Static" : "Dynamic"
      private_ip_address            = frontend_ip_configuration.value.ipv6 ? cidrhost(local.master_subnet_cidr_v6, -2) : null
    }
  }
}

resource "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool" {
  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.internal.id
  name                = "${var.cluster_id}-internal-controlplane"
}

resource "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool_v6" {
  count = var.use_ipv6 ? 1 : 0

  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurerm_lb.internal.id
  name                = "${var.cluster_id}-internal-controlplane-v6"
}

resource "azurerm_lb_rule" "internal_lb_rule_api_internal" {
  name                           = "api-internal"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool.id
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_api_internal.id
}

resource "azurerm_lb_rule" "internal_lb_rule_api_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                           = "api-internal-v6"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[0].id
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_api_internal.id
}

resource "azurerm_lb_rule" "internal_lb_rule_sint" {
  name                           = "sint"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool.id
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 22623
  backend_port                   = 22623
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_sint.id
}

resource "azurerm_lb_rule" "internal_lb_rule_sint_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                           = "sint-v6"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[0].id
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 22623
  backend_port                   = 22623
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_sint.id
}

resource "azurerm_lb_probe" "internal_lb_probe_sint" {
  name                = "sint-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 10
  number_of_probes    = 3
  loadbalancer_id     = azurerm_lb.internal.id
  port                = 22623
  protocol            = "TCP"
}

resource "azurerm_lb_probe" "internal_lb_probe_api_internal" {
  name                = "api-internal-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 10
  number_of_probes    = 3
  loadbalancer_id     = azurerm_lb.internal.id
  port                = 6443
  protocol            = "TCP"
}
