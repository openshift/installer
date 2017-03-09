resource "azurerm_lb" "tectonic_console_lb" {
  name                = "console-lb"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  frontend_ip_configuration {
    name                          = "default"
    public_ip_address_id          = "${azurerm_public_ip.tectonic_console_ip.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_public_ip" "tectonic_console_ip" {
  name                         = "tectonic_console_ip"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  public_ip_address_allocation = "static"

  tags {
    environment = "staging"
  }
}

resource "azurerm_lb_rule" "console-lb-https" {
  name                    = "console-lb-rule-443-32000"
  resource_group_name     = "${var.resource_group_name}"
  loadbalancer_id         = "${azurerm_lb.tectonic_console_lb.id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.k8-lb.id}"
  probe_id                = "${azurerm_lb_probe.console-lb.id}"

  protocol                       = "tcp"
  frontend_port                  = 443
  backend_port                   = 32000
  frontend_ip_configuration_name = "default"
}

resource "azurerm_lb_rule" "console-lb-identity" {
  name                    = "console-lb-rule-80-32001"
  resource_group_name     = "${var.resource_group_name}"
  loadbalancer_id         = "${azurerm_lb.tectonic_console_lb.id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.k8-lb.id}"
  probe_id                = "${azurerm_lb_probe.console-lb.id}"

  protocol                       = "tcp"
  frontend_port                  = 80
  backend_port                   = 32001
  frontend_ip_configuration_name = "default"
}

resource "azurerm_lb_probe" "console-lb" {
  name                = "console-lb-probe-443-up"
  loadbalancer_id     = "${azurerm_lb.tectonic_console_lb.id}"
  resource_group_name = "${var.resource_group_name}"
  protocol            = "tcp"
  port                = 32000
}
