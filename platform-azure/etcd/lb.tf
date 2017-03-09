resource "azurerm_lb" "tectonic_etcd_lb" {
  name                = "etcd-lb"
  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"

  frontend_ip_configuration {
    name                          = "default"
    public_ip_address_id          = "${azurerm_public_ip.etcd_publicip.id}"
    private_ip_address_allocation = "dynamic"
  }
}

resource "azurerm_public_ip" "etcd_publicip" {
  name                         = "${var.tectonic_cluster_name}_etcd_publicip"
  location                     = "${var.tectonic_azure_location}"
  resource_group_name          = "${var.tectonic_azure_resource_group_name}"
  public_ip_address_allocation = "dynamic"
}

resource "azurerm_lb_rule" "etcd-lb" {
  name                           = "etcd-lb-rule-2379-2379"
  resource_group_name            = "${var.tectonic_azure_resource_group_name}"
  loadbalancer_id                = "${azurerm_lb.tectonic_etcd_lb.id}"
  backend_address_pool_id        = "${azurerm_lb_backend_address_pool.etcd-lb.id}"
  probe_id                       = "${azurerm_lb_probe.etcd-lb.id}"
  protocol                       = "tcp"
  frontend_port                  = 2379
  backend_port                   = 2379
  frontend_ip_configuration_name = "default"
}

resource "azurerm_lb_probe" "etcd-lb" {
  name                = "etcd-lb-probe"
  loadbalancer_id     = "${azurerm_lb.tectonic_etcd_lb.id}"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
  protocol            = "Tcp"
  port                = 2379
}

resource "azurerm_lb_backend_address_pool" "etcd-lb" {
  name                = "etcd-lb-pool"
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
  loadbalancer_id     = "${azurerm_lb.tectonic_etcd_lb.id}"
}

resource "azurerm_dns_a_record" "tectonic-etcd" {
  resource_group_name = "${var.tectonic_azure_resource_group_name}"
  zone_name           = "${var.tectonic_azure_dns_zone_name}"
  name                = "etcd-${var.tectonic_cluster_name}"
  ttl                 = "60"
  depends_on          = ["azurerm_public_ip.etcd_publicip"]
  records             = ["${azurerm_public_ip.etcd_publicip.ip_address}"]
}
