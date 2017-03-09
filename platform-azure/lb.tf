resource "azurerm_lb" "tectonic_api_lb" {
  name                = "k8-lb"
  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"

  frontend_ip_configuration {
    name                          = "default"
    public_ip_address_id          = "${azurerm_public_ip.tectonic_master_ip.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_lb_backend_address_pool" "k8-lb" {
  name                = "k8-lb-pool"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  loadbalancer_id     = "${azurerm_lb.tectonic_api_lb.id}"
}

resource "azurerm_lb_rule" "k8-lb" {
  name                    = "k8-lb-rule-443-443"
  resource_group_name     = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  loadbalancer_id         = "${azurerm_lb.tectonic_api_lb.id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.k8-lb.id}"
  probe_id                = "${azurerm_lb_probe.k8-lb.id}"

  protocol                       = "tcp"
  frontend_port                  = 443
  backend_port                   = 443
  frontend_ip_configuration_name = "default"

}

resource "azurerm_lb_probe" "k8-lb" {
  name                = "k8-lb-probe-443-up"
  loadbalancer_id     = "${azurerm_lb.tectonic_api_lb.id}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  protocol            = "tcp"
  port                = 443
}

resource "azurerm_lb_rule" "ssh-lb" {
  name                    = "ssh-lb"
  resource_group_name     = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  loadbalancer_id         = "${azurerm_lb.tectonic_api_lb.id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.k8-lb.id}"
  probe_id                = "${azurerm_lb_probe.ssh-lb.id}"

  protocol                       = "tcp"
  frontend_port                  = 22
  backend_port                   = 22
  frontend_ip_configuration_name = "default"
}

resource "azurerm_lb_probe" "ssh-lb" {
  name                = "ssh-lb-22-up"
  loadbalancer_id     = "${azurerm_lb.tectonic_api_lb.id}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  protocol            = "tcp"
  port                = 22
}


