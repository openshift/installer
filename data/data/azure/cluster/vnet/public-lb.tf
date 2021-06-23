locals {
  need_public_ipv4 = ! var.private || ! var.outbound_udr
  need_public_ipv6 = var.use_ipv6 && (! var.private || ! var.outbound_udr)
}

data "azurerm_public_ip" "cluster_public_ip_v4" {
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
