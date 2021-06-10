locals {
  public_lb_frontend_ip_v4_configuration_name = "public-lb-ip-v4"
  public_lb_frontend_ip_v6_configuration_name = "public-lb-ip-v6"
}

locals {
  // DEBUG: Azure apparently requires dual stack LB for v6
  need_public_ipv4 = ! var.private || ! var.outbound_udr

  need_public_ipv6 = var.use_ipv6 && (! var.private || ! var.outbound_udr)
}

data "azurerm_public_ip" "cluster_public_ip_v4" {
  // DEBUG: Azure apparently requires dual stack LB for v6
  count = local.need_public_ipv4 ? 1 : 0

  name                = "${var.cluster_id}-pip-v4"
  resource_group_name = var.resource_group_name
}

data "azurerm_public_ip" "cluster_public_ip_v6" {
  count = local.need_public_ipv6 ? 1 : 0

  name                = "${var.cluster_id}-pip-v6"
  resource_group_name = var.resource_group_name
}

data "azurerm_lb" "public" {
  name                = var.cluster_id
  resource_group_name = var.resource_group_name
}

// The backends are only created when frontend configuration exists, because of the following error from Azure API;
// ```
// Load Balancer /subscriptions/xx/resourceGroups/xx/providers/Microsoft.Network/loadBalancers/xx-public-lb does not have Frontend IP Configuration, 
// but it has other child resources. This setup is not supported.
// ```
data "azurerm_lb_backend_address_pool" "public_lb_pool_v4" {
  count = local.need_public_ipv4 ? 1 : 0

  loadbalancer_id = data.azurerm_lb.public.id
  name            = var.cluster_id
}

data "azurerm_lb_backend_address_pool" "public_lb_pool_v6" {
  count = local.need_public_ipv6 ? 1 : 0

  loadbalancer_id = data.azurerm_lb.public.id
  name            = "${var.cluster_id}-IPv6"
}

data "azurerm_lb_rule" "public_lb_rule_api_internal_v4" {
  count = var.use_ipv4 && ! var.private ? 1 : 0

  name                = "api-internal-v4"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.public.id
}

data "azurerm_lb_rule" "public_lb_rule_api_internal_v6" {
  count = var.use_ipv6 && ! var.private ? 1 : 0

  name                = "api-internal-v6"
  resource_group_name = var.resource_group_name
  loadbalancer_id     = data.azurerm_lb.public.id
}
