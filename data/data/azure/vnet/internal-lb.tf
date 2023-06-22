locals {
  internal_lb_frontend_ip_v4_configuration_name = "internal-lb-ip-v4"
  internal_lb_frontend_ip_v6_configuration_name = "internal-lb-ip-v6"
  internal_lb_frontend_ip_v4_configuration = [for idx in range(length(local.master_subnet_id)) : {
    name : "${local.internal_lb_frontend_ip_v4_configuration_name}-${idx}",
    ipv6 : false,
    include : var.use_ipv4,
    subnet_id : local.master_subnet_id[idx]
    private_ip_address : null
    zones : var.azure_outbound_routing_type == "NatGateway" ? [local.master_subnet_zone_map[local.master_subnet_id[idx]]] : null
  }]
  internal_lb_frontend_ip_v6_configuration = [for idx in range(length(local.master_subnet_id)) : {
    name : "${local.internal_lb_frontend_ip_v6_configuration_name}-${idx}",
    ipv6 : true,
    include : var.use_ipv6,
    subnet_id : local.master_subnet_id[idx]
    private_ip_address : var.use_ipv6 ? cidrhost(local.master_subnet_cidr_v6, -2 - idx) : null
    zones : var.azure_outbound_routing_type == "NatGateway" ? [local.master_subnet_zone_map[local.master_subnet_id[idx]]] : null
  }]
}

resource "azurerm_lb" "internal" {
  sku                 = "Standard"
  name                = "${var.cluster_id}-internal"
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region

  dynamic "frontend_ip_configuration" {
    for_each = [for ip in flatten([
      local.internal_lb_frontend_ip_v4_configuration,
      local.internal_lb_frontend_ip_v6_configuration,
      ]) : ip if ip.include
    ]

    content {
      name                       = frontend_ip_configuration.value.name
      subnet_id                  = frontend_ip_configuration.value.subnet_id
      private_ip_address_version = frontend_ip_configuration.value.ipv6 ? "IPv6" : "IPv4"
      # WORKAROUND: Allocate a high ipv6 internal LB address to avoid the race with NIC allocation (a master and the LB
      #   were being assigned the same IP dynamically). Issue is being tracked as a support ticket to Azure.
      private_ip_address_allocation = frontend_ip_configuration.value.ipv6 ? "Static" : "Dynamic"
      private_ip_address            = frontend_ip_configuration.value.private_ip_address
      zones                         = frontend_ip_configuration.value.zones
    }
  }

  tags = var.azure_extra_tags
}

resource "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool_v4" {
  count = var.use_ipv4 ? length(local.internal_lb_frontend_ip_v4_configuration) : 0

  loadbalancer_id = azurerm_lb.internal.id
  // WORKAROUND: the cloud provider is not happy when we change the name for private clusters
  name = count.index > 0 ? "${var.cluster_id}-${count.index}" : "${var.cluster_id}"
}

resource "azurerm_lb_backend_address_pool" "internal_lb_controlplane_pool_v6" {
  count = var.use_ipv6 ? length(local.internal_lb_frontend_ip_v6_configuration) : 0

  loadbalancer_id = azurerm_lb.internal.id
  // WORKAROUND: the cloud provider is not happy when we change the name for private clusters
  name = count.index > 0 ? "${var.cluster_id}-IPv6-${count.index}" : "${var.cluster_id}-IPv6"
}

resource "azurerm_lb_rule" "internal_lb_rule_api_internal_v4" {
  count = var.use_ipv4 ? length(local.internal_lb_frontend_ip_v4_configuration) : 0

  name                           = "api-internal-v4-${count.index}"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v4[count.index].id]
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v4_configuration[count.index].name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_api_internal.id
}

resource "azurerm_lb_rule" "internal_lb_rule_api_internal_v6" {
  count = var.use_ipv6 ? length(local.internal_lb_frontend_ip_v6_configuration) : 0

  name                           = "api-internal-v6-${count.index}"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[count.index].id]
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v6_configuration[count.index].name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_api_internal.id
}

resource "azurerm_lb_rule" "internal_lb_rule_sint_v4" {
  count = var.use_ipv4 ? length(local.internal_lb_frontend_ip_v4_configuration) : 0

  name                           = "sint-v4-${count.index}"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v4[count.index].id]
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 22623
  backend_port                   = 22623
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v4_configuration[count.index].name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_sint.id
}

resource "azurerm_lb_rule" "internal_lb_rule_sint_v6" {
  count = var.use_ipv6 ? length(local.internal_lb_frontend_ip_v6_configuration) : 0

  name                           = "sint-v6-${count.index}"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.internal_lb_controlplane_pool_v6[count.index].id]
  loadbalancer_id                = azurerm_lb.internal.id
  frontend_port                  = 22623
  backend_port                   = 22623
  frontend_ip_configuration_name = local.internal_lb_frontend_ip_v6_configuration[count.index].name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.internal_lb_probe_sint.id
}

resource "azurerm_lb_probe" "internal_lb_probe_sint" {
  name                = "sint-probe"
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurerm_lb.internal.id
  port                = 22623
  protocol            = "Https"
  request_path        = "/healthz"
}

resource "azurerm_lb_probe" "internal_lb_probe_api_internal" {
  name                = "api-internal-probe"
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurerm_lb.internal.id
  port                = 6443
  protocol            = "Https"
  request_path        = "/readyz"
}
