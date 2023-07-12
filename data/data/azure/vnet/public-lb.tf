locals {
  public_lb_frontend_ip_v4_configuration_name = "public-lb-ip-v4"
  public_lb_frontend_ip_v6_configuration_name = "public-lb-ip-v6"
}

locals {
  // DEBUG: Azure apparently requires dual stack LB for v6
  need_public_ipv4 = ! var.azure_private || var.azure_outbound_routing_type != "UserDefinedRouting"

  need_public_ipv6 = var.use_ipv6 && (! var.azure_private || var.azure_outbound_routing_type != "UserDefinedRouting")
}


resource "azurerm_public_ip" "cluster_public_ip_v4" {
  count = local.need_public_ipv4 ? 1 : 0

  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-pip-v4"
  resource_group_name = data.azurerm_resource_group.main.name
  allocation_method   = "Static"
  domain_name_label   = var.cluster_id
  tags                = var.azure_extra_tags
}

data "azurerm_public_ip" "cluster_public_ip_v4" {
  // DEBUG: Azure apparently requires dual stack LB for v6
  count = local.need_public_ipv4 ? 1 : 0

  name                = azurerm_public_ip.cluster_public_ip_v4[0].name
  resource_group_name = data.azurerm_resource_group.main.name
}


resource "azurerm_public_ip" "cluster_public_ip_v6" {
  count = local.need_public_ipv6 ? 1 : 0

  ip_version          = "IPv6"
  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-pip-v6"
  resource_group_name = data.azurerm_resource_group.main.name
  allocation_method   = "Static"
  domain_name_label   = var.cluster_id
  tags                = var.azure_extra_tags
}

data "azurerm_public_ip" "cluster_public_ip_v6" {
  count = local.need_public_ipv6 ? 1 : 0

  name                = azurerm_public_ip.cluster_public_ip_v6[0].name
  resource_group_name = data.azurerm_resource_group.main.name
}

resource "azurerm_lb" "public" {
  count               = local.need_public_ipv4 ? 1 : 0
  sku                 = "Standard"
  name                = var.cluster_id
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region

  dynamic "frontend_ip_configuration" {
    for_each = [for ip in [
      // DEBUG: Azure apparently requires dual stack LB for external load balancers v6
      {
        name : local.public_lb_frontend_ip_v4_configuration_name,
        value : local.need_public_ipv4 ? azurerm_public_ip.cluster_public_ip_v4[0].id : null,
        include : local.need_public_ipv4,
        ipv6 : false,
      },
      {
        name : local.public_lb_frontend_ip_v6_configuration_name,
        value : local.need_public_ipv6 ? azurerm_public_ip.cluster_public_ip_v6[0].id : null,
        include : local.need_public_ipv6,
        ipv6 : true,
      },
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

  tags = var.azure_extra_tags
}

// The backends are only created when frontend configuration exists, because of the following error from Azure API;
// ```
// Load Balancer /subscriptions/xx/resourceGroups/xx/providers/Microsoft.Network/loadBalancers/xx-public-lb does not have Frontend IP Configuration, 
// but it has other child resources. This setup is not supported.
// ```
resource "azurerm_lb_backend_address_pool" "public_lb_pool_v4" {
  count = local.need_public_ipv4 ? 1 : 0

  loadbalancer_id = azurerm_lb.public[0].id
  name            = var.cluster_id
}

resource "azurerm_lb_backend_address_pool" "public_lb_pool_v6" {
  count = local.need_public_ipv6 ? 1 : 0

  loadbalancer_id = azurerm_lb.public[0].id
  name            = "${var.cluster_id}-IPv6"
}

resource "azurerm_lb_rule" "public_lb_rule_api_internal_v4" {
  count = var.use_ipv4 && ! var.azure_private ? 1 : 0

  name                           = "api-internal-v4"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id]
  loadbalancer_id                = azurerm_lb.public[0].id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v4_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurerm_lb_rule" "public_lb_rule_api_internal_v6" {
  count = var.use_ipv6 && ! var.azure_private ? 1 : 0

  name                           = "api-internal-v6"
  protocol                       = "Tcp"
  backend_address_pool_ids       = [azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id]
  loadbalancer_id                = azurerm_lb.public[0].id
  frontend_port                  = 6443
  backend_port                   = 6443
  frontend_ip_configuration_name = local.public_lb_frontend_ip_v6_configuration_name
  enable_floating_ip             = false
  idle_timeout_in_minutes        = 30
  load_distribution              = "Default"
  probe_id                       = azurerm_lb_probe.public_lb_probe_api_internal[0].id
}

resource "azurerm_lb_outbound_rule" "public_lb_outbound_rule_v4" {
  count = var.use_ipv4 && var.azure_private && var.azure_outbound_routing_type != "UserDefinedRouting" ? 1 : 0

  name                    = "outbound-rule-v4"
  loadbalancer_id         = azurerm_lb.public[0].id
  backend_address_pool_id = azurerm_lb_backend_address_pool.public_lb_pool_v4[0].id
  protocol                = "All"

  frontend_ip_configuration {
    name = local.public_lb_frontend_ip_v4_configuration_name
  }
}

resource "azurerm_lb_outbound_rule" "public_lb_outbound_rule_v6" {
  count = var.use_ipv6 && var.azure_private && var.azure_outbound_routing_type != "UserDefinedRouting" ? 1 : 0

  name                    = "outbound-rule-v6"
  loadbalancer_id         = azurerm_lb.public[0].id
  backend_address_pool_id = azurerm_lb_backend_address_pool.public_lb_pool_v6[0].id
  protocol                = "All"

  frontend_ip_configuration {
    name = local.public_lb_frontend_ip_v6_configuration_name
  }
}

resource "azurerm_lb_probe" "public_lb_probe_api_internal" {
  count = var.azure_private ? 0 : 1

  name                = "api-internal-probe"
  interval_in_seconds = 5
  number_of_probes    = 2
  loadbalancer_id     = azurerm_lb.public[0].id
  port                = 6443
  protocol            = "Https"
  request_path        = "/readyz"
}

resource "azurerm_public_ip" "ngw_public_ip_v4" {
  count = local.need_public_ipv4 && var.azure_outbound_routing_type == "NatGateway" ? 1 : 0

  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-natgw-pip-v4"
  resource_group_name = data.azurerm_resource_group.main.name
  allocation_method   = "Static"
  tags                = var.azure_extra_tags
}

resource "azurerm_public_ip" "ngw_public_ip_v6" {
  count = local.need_public_ipv6 && var.azure_outbound_routing_type == "NatGateway" ? 1 : 0

  ip_version          = "IPv6"
  sku                 = "Standard"
  location            = var.azure_region
  name                = "${var.cluster_id}-natgw-pip-v6"
  resource_group_name = data.azurerm_resource_group.main.name
  allocation_method   = "Static"
  tags                = var.azure_extra_tags
}

resource "azurerm_nat_gateway" "nat_gw" {
  count                   = var.azure_outbound_routing_type == "NatGateway" ? 1 : 0
  name                    = "${var.cluster_id}-natgw"
  location                = var.azure_region
  resource_group_name     = data.azurerm_resource_group.main.name
  sku_name                = "Standard"
  idle_timeout_in_minutes = 10
  tags                    = var.azure_extra_tags
  # By not specifying zones here, we make the NAT non-zonal
  #zones = ["1"]
}

data "azurerm_nat_gateway" "nat_gw" {
  count               = var.azure_outbound_routing_type == "NatGateway" ? 1 : 0
  name                = azurerm_nat_gateway.nat_gw[0].name
  resource_group_name = data.azurerm_resource_group.main.name
}

resource "azurerm_nat_gateway_public_ip_association" "nat_ip_assoc" {
  count                = var.azure_outbound_routing_type == "NatGateway" ? 1 : 0
  nat_gateway_id       = azurerm_nat_gateway.nat_gw[0].id
  public_ip_address_id = var.use_ipv6 ? azurerm_public_ip.ngw_public_ip_v6[0].id : azurerm_public_ip.ngw_public_ip_v4[0].id
}
