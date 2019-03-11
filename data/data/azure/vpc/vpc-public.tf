resource "azurerm_subnet" "public_subnet" {
  resource_group_name = "${var.rg_name}"
  address_prefix = "${cidrsubnet(local.new_public_cidr_range, 3, 0)}"
  virtual_network_name = "${var.cluster_id}-public-subnet"

  tags = "${merge(map(
    "Name", "${var.cluster_id}-public-subnet",
  ), var.tags)}"
}

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