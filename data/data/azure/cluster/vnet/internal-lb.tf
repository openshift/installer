locals {
  internal_lb_frontend_ip_v4_configuration_name = "internal-lb-ip-v4"
  internal_lb_frontend_ip_v6_configuration_name = "internal-lb-ip-v6"
}

data "azurerm_lb" "internal" {
  name                = "${var.cluster_id}-internal"
  resource_group_name = var.resource_group_name
}

data "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool_v4" {
  count = var.use_ipv4 ? 1 : 0

  loadbalancer_id = data.azurerm_lb.internal.id
  name            = var.cluster_id
}

data "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool_v6" {
  count = var.use_ipv6 ? 1 : 0

  loadbalancer_id = data.azurerm_lb.internal.id
  name            = "${var.cluster_id}-IPv6"
}

data "azurerm_lb_rule" "internal_lb_rule_api_internal_v4" {
  count = var.use_ipv4 ? 1 : 0

  name                = "api-internal-v4"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.internal.id
}

data "azurerm_lb_rule" "internal_lb_rule_api_internal_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "api-internal-v6"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.internal.id
}

data "azurerm_lb_rule" "internal_lb_rule_sint_v4" {
  count = var.use_ipv4 ? 1 : 0

  name                = "sint-v4"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.internal.id
}

data "azurerm_lb_rule" "internal_lb_rule_sint_v6" {
  count = var.use_ipv6 ? 1 : 0

  name                = "sint-v6"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.internal.id
}
