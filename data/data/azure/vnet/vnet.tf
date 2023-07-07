resource "azurerm_virtual_network" "cluster_vnet" {
  count = var.azure_preexisting_network ? 0 : 1

  name                = var.azure_virtual_network
  resource_group_name = data.azurerm_resource_group.main.name
  location            = var.azure_region
  address_space       = concat(var.machine_v4_cidrs, var.machine_v6_cidrs)
  tags                = var.azure_extra_tags
}

resource "azurerm_subnet" "master_subnet" {
  count = var.azure_preexisting_network ? 0 : 1

  resource_group_name = data.azurerm_resource_group.main.name
  address_prefixes = [for cidr in [
    { value : local.master_subnet_cidr_v4, include : var.use_ipv4 },
    { value : local.master_subnet_cidr_v6, include : var.use_ipv6 }
  ] : cidr.value if cidr.include]
  virtual_network_name = local.virtual_network
  name                 = var.azure_control_plane_subnet
}

resource "azurerm_subnet" "worker_subnet" {
  count = var.azure_preexisting_network ? 0 : 1

  resource_group_name = data.azurerm_resource_group.main.name
  address_prefixes = [for cidr in [
    { value : local.worker_subnet_cidr_v4, include : var.use_ipv4 },
    { value : local.worker_subnet_cidr_v6, include : var.use_ipv6 }
  ] : cidr.value if cidr.include]
  virtual_network_name = local.virtual_network
  name                 = var.azure_compute_subnet
}

resource "azurerm_subnet_nat_gateway_association" "nat_master_assoc" {
  count          = var.azure_outbound_routing_type == "NatGateway" ? 1 : 0
  subnet_id      = local.master_subnet_id
  nat_gateway_id = azurerm_nat_gateway.nat_gw[0].id
}

resource "azurerm_subnet_nat_gateway_association" "nat_worker_assoc" {
  count          = var.azure_outbound_routing_type == "NatGateway" ? 1 : 0
  subnet_id      = local.worker_subnet_id
  nat_gateway_id = azurerm_nat_gateway.nat_gw[0].id
}
