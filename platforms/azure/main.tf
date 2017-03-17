# Configure the Microsoft Azure Provider
provider "azurerm" {
  subscription_id = "***REMOVED***"
  client_id       = "***REMOVED***"
  client_secret   = "***REMOVED***"
  tenant_id       = "***REMOVED***"
}

resource "azurerm_resource_group" "tectonic_azure_cluster_resource_group" {
  location = "${var.tectonic_azure_location}"
  name     = "tectonic-cluster-${var.tectonic_cluster_name}-group-alt"
}

resource "azurerm_virtual_network" "tectonic_vnet" {
  name                = "tectonic_vnet"
  address_space       = ["10.0.0.0/16"]
  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
}

module "etcd" {
  source = "../../modules/azure/etcd"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  image_reference     = "${var.tectonic_azure_image_reference}"
  vm_size             = "${var.tectonic_azure_vm_size}"

  count           = "${var.tectonic_master_count}"
  base_domain     = "${var.tectonic_base_domain}"
  cluster_name    = "${var.tectonic_cluster_name}"
  ssh_key         = "${var.tectonic_ssh_key}"
  virtual_network = "${azurerm_virtual_network.tectonic_vnet.name}"
}

module "masters" {
  source = "../../modules/azure/master"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  image_reference     = "${var.tectonic_azure_image_reference}"
  vm_size             = "${var.tectonic_azure_vm_size}"

  count           = "${var.tectonic_master_count}"
  base_domain     = "${var.tectonic_base_domain}"
  cluster_name    = "${var.tectonic_cluster_name}"
  ssh_key         = "${var.tectonic_ssh_key}"
  virtual_network = "${azurerm_virtual_network.tectonic_vnet.name}"
  kube_image_url  = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
  kube_image_tag  = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
}

module "workers" {
  source = "../../modules/azure/worker"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  image_reference     = "${var.tectonic_azure_image_reference}"
  vm_size             = "${var.tectonic_azure_vm_size}"

  worker_count    = "${var.tectonic_worker_count}"
  base_domain     = "${var.tectonic_base_domain}"
  cluster_name    = "${var.tectonic_cluster_name}"
  ssh_key         = "${var.tectonic_ssh_key}"
  virtual_network = "${azurerm_virtual_network.tectonic_vnet.name}"
  kube_image_url  = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
  kube_image_tag  = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
}

module "dns" {
  source = "../../modules/azure/dns"

  master_ip_addresses = "${module.masters.ip_address}"
  console_ip_address  = "${module.masters.console_ip_address}"
  etcd_ip_addresses   = "${module.etcd.ip_address}"

  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  location            = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_dns_resource_group}"

  // TODO etcd list
  // TODO worker list
}
