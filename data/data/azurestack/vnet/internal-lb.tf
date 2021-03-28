locals {
  internal_lb_frontend_ip_v4_configuration_name = "internal-lb-ip-v4"
}

resource "azurestack_lb" "internal" {
  name                = "${var.cluster_id}-internal"
  resource_group_name = var.resource_group_name
  location            = var.region

  frontend_ip_configuration {
    name                          = local.internal_lb_frontend_ip_v4_configuration_name
    subnet_id                     = local.master_subnet_id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurestack_lb_backend_address_pool" "internal_lb_controlplane_pool_v4" {
  resource_group_name = var.resource_group_name
  loadbalancer_id     = azurestack_lb.internal.id
  name                = var.cluster_id
}

resource "azurestack_lb_rule" "internal_lb_rule_api_internal_v4" {
  name                           = "api-internal-v4"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurestack_lb_backend_address_pool.internal_lb_controlplane_pool_v4.id
  loadbalancer_id                = azurestack_lb.internal.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v4_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurestack_lb_probe.internal_lb_probe_api_internal.id
}

resource "azurestack_lb_rule" "internal_lb_rule_sint_v4" {
  name                           = "sint-v4"
  resource_group_name            = var.resource_group_name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurestack_lb_backend_address_pool.internal_lb_controlplane_pool_v4.id
  loadbalancer_id                = azurestack_lb.internal.id
  frontend_port                  = 22623
  backend_port                   = 22623
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v4_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurestack_lb_probe.internal_lb_probe_sint.id
}

resource "azurestack_lb_probe" "internal_lb_probe_sint" {
  name                = "sint-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurestack_lb.internal.id
  port                = 22623
  protocol            = "TCP"
}

resource "azurestack_lb_probe" "internal_lb_probe_api_internal" {
  name                = "api-internal-probe"
  resource_group_name = var.resource_group_name
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurestack_lb.internal.id
  port                = 6443
  protocol            = "TCP"
}
