resource "azurerm_public_ip" "cluster_public_ip" {
  location= "${var.region}"
  name = "${var.cluster_id}-pip"
  resource_group_name="${var.rg_name}"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-pip",
  ), var.tags)}"

}

resource "azurerm_lb" "external_lb" {
  name                = "${var.cluster_id}-elb"
  resource_group_name = "${var.rg_name}"
  location            = "${var.region}"
  
  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.cluster_public_ip.id}"
  }

  tags = "${merge(map(
    "Name", "${var.cluster_id}-elb",
  ), var.tags)}"
}

resource "azurerm_lb_backend_address_pool" "master_elb_pool" {
  resource_group_name = "${var.rg_name}"
  location            = "${var.region}"
  loadbalancer_id = "${azurerm_lb.external_lb.id}"
  name = "${var.cluster_id}-elb-master"

   tags = "${merge(map(
    "Name", "${var.cluster_id}-elb-pool-master",
  ), var.tags)}"
}

resource "azurerm_network_interface_backend_address_pool_association" "master" {
  network_interface_id    = "${local.vnet_id}"
  ip_configuration_name   = "master-nic-association"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.master_elb_pool.id}"
}

resource "azurerm_lb_rule" "api_internal" {
  resource_group_name = "${var.rg_name}"
  location            = "${var.region}"
  backend_address_pool_id = "${azurerm_lb_backend_address_pool.master_elb_pool.id}"
  loadbalancer_id = "${azurerm_lb.external_lb.id}"
  frontend_port = 6443
}
