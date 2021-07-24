locals {
  public_lb_frontend_ip_v4_configuration_name = "public-lb-ip-v4"
}


resource "azurestack_public_ip" "cluster_public_ip_v4" {
  count = ! var.azure_private ? 1 : 0

  location                     = var.azure_region
  name                         = "${var.cluster_id}-pip-v4"
  resource_group_name          = data.azurestack_resource_group.main.name
  public_ip_address_allocation = "Static"
  domain_name_label            = var.cluster_id
}

resource "azurestack_lb" "public" {
  count = ! var.azure_private ? 1 : 0

  name                = var.cluster_id
  resource_group_name = data.azurestack_resource_group.main.name
  location            = var.azure_region

  frontend_ip_configuration {
    name                          = local.public_lb_frontend_ip_v4_configuration_name
    public_ip_address_id          = azurestack_public_ip.cluster_public_ip_v4[0].id
    private_ip_address_allocation = "Dynamic"
  }
}

resource "azurestack_lb_backend_address_pool" "public_lb_pool_v4" {
  count = ! var.azure_private ? 1 : 0

  resource_group_name = data.azurestack_resource_group.main.name
  loadbalancer_id     = azurestack_lb.public[0].id
  name                = var.cluster_id
}

resource "azurestack_lb_rule" "public_lb_rule_api_internal_v4" {
  count = ! var.azure_private ? 1 : 0

  name                           = "api-internal-v4"
  resource_group_name            = data.azurestack_resource_group.main.name
  protocol                       = "Tcp"
  backend_address_pool_id        = azurestack_lb_backend_address_pool.public_lb_pool_v4[0].id
  loadbalancer_id                = azurestack_lb.public[0].id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v4_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurestack_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurestack_lb_probe" "public_lb_probe_api_internal" {
  count = ! var.azure_private ? 1 : 0

  name                = "api-internal-probe"
  resource_group_name = data.azurestack_resource_group.main.name
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurestack_lb.public[0].id
  port                = 6443
  protocol            = "TCP"
}
