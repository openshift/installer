resource "azurerm_lb" "tectonic_k8" {
  name                = "k8-lb"
  location            = "${var.azurerm_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"

  frontend_ip_configuration {
    name                          = "default"
    public_ip_address_id          = "${azurerm_public_ip.tectonic_azure_k8_public_ip.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_lb_rule" "k8-lb" {
  name                    = "k8-lb-rule-80-80"
  resource_group_name     = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  loadbalancer_id         = "${azurerm_lb.k8-lb.id}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.k8-lb.id}"
  probe_id                = "${azurerm_lb_probe.k8-lb.id}"

  protocol                       = "tcp"
  frontend_port                  = 80
  backend_port                   = 80
  frontend_ip_configuration_name = "default"

}

resource "azurerm_lb_probe" "k8-lb" {
  name                = "k8-lb-probe-80-up"
  loadbalancer_id     = "${azurerm_lb.k8-lb.id}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  protocol            = "Http"
  request_path        = "/"
  port                = 80
}

resource "azurerm_lb_backend_address_pool" "k8-lb" {
  name                = "k8-lb-pool"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  loadbalancer_id     = "${azurerm_lb.k8-lb.id}"
}