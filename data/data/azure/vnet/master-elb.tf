resource "azurerm_public_ip" "cluster_public_ip" {
  sku                 = "Standard"
  location            = "${var.region}"
  name                = "${var.cluster_id}-pip"
  resource_group_name = "${var.rg_name}"
  allocation_method   = "Static"
  domain_name_label   = "${var.dns_label}"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-pip",
  ), var.tags)}"
}

data "azurerm_public_ip" "cluster_public_ip" {
  name                = "${azurerm_public_ip.cluster_public_ip.name}"
  resource_group_name = "${var.rg_name}"
}

resource "azurerm_lb" "external_lb" {
  sku                 = "Standard"
  name                = "${var.cluster_id}-elb"
  resource_group_name = "${var.rg_name}"
  location            = "${var.region}"
  
  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.cluster_public_ip.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "master_elb_pool" {
  resource_group_name = "${var.rg_name}"
  loadbalancer_id = "${azurerm_lb.external_lb.id}"
  name = "${var.cluster_id}-elb-master"
}

resource "azurerm_lb_rule" "external_lb_rule_api_internal" {
  name = "api-internal"
  resource_group_name = "${var.rg_name}"
  protocol="Tcp"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.master_elb_pool.id}"
  loadbalancer_id = "${azurerm_lb.external_lb.id}"
  frontend_port = 6443
  backend_port = 6443
  frontend_ip_configuration_name = "PublicIPAddress"
  enable_floating_ip = false
  idle_timeout_in_minutes = 4
  load_distribution = "Default"
  probe_id = "${azurerm_lb_probe.external_lb_probe_api_internal.id}"
}

resource "azurerm_lb_probe" "external_lb_probe_api_internal" {
  name = "api-internal-probe"
  resource_group_name = "${var.rg_name}"
  interval_in_seconds = 15
  number_of_probes = 4
  loadbalancer_id = "${azurerm_lb.external_lb.id}"
  port = 6443
  protocol = "Tcp"
}


