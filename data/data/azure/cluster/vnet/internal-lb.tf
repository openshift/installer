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
