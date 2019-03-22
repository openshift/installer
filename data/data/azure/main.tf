locals {
  tags = "${merge(map(
    "kubernetes.io_cluster.${var.cluster_id}", "owned"
  ), var.azure_extra_tags)}"
}

provider "azurerm" {
   # version = "~>1.5"
}

module "bootstrap" {
  source = "./bootstrap"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  region                   = "${var.azure_region}"
  vm_size                  = "${var.azure_bootstrap_vm_type}"
  
  cluster_id               = "${var.cluster_id}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vnet.public_subnet_id}"
  elb_backend_pool_id      = "${module.vnet.elb_backend_pool_id}"
  ilb_backend_pool_id      = "${module.vnet.ilb_backend_pool_id}"
  tags                     = "${local.tags}"
  boot_diag_blob_endpoint  = "${azurerm_storage_account.bootdiag.primary_blob_endpoint}"
}

module "vnet" {
  source = "./vnet"
  resource_group_name =   "${azurerm_resource_group.main.name}"
  cidr_block = "${var.machine_cidr}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.azure_region}"
  dns_label  = "${var.cluster_id}"
  tags = "${local.tags}"
}

module "master" {
  source = "./master"
  resource_group_name = "${azurerm_resource_group.main.name}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.azure_region}"
  vm_size    = "${var.azure_bootstrap_vm_type}"
  ignition   = "${var.ignition_master}"
  external_lb_id = "${module.vnet.external_lb_id}"
  elb_backend_pool_id = "${module.vnet.elb_backend_pool_id}"
  ilb_backend_pool_id = "${module.vnet.ilb_backend_pool_id}"
  subnet_id = "${module.vnet.public_subnet_id}"
  instance_count = "${var.master_count}"
  tags = "${local.tags}"
  boot_diag_blob_endpoint  = "${azurerm_storage_account.bootdiag.primary_blob_endpoint}"
}

module "dns" {
  source = "./dns"
  cluster_domain = "${var.cluster_domain}"
  base_domain = "${var.base_domain}"
  cluster_name = "${var.cluster_name}"
  external_lb_dns_label = "${module.vnet.external_lb_dns_label}"
  internal_lb_ipaddress = "${module.vnet.internal_lb_ip_address}"
  resource_group_name = "${azurerm_resource_group.main.name}"
  internal_dns_resolution_vnet_id = "${module.vnet.vnet_id}"
}

resource "random_string" "resource_group_suffix" {
  length = 5
  upper = false
  special = false
}

resource "azurerm_resource_group" "main" {
  name = "${var.cluster_id}-rg"
  location = "${var.azure_region}"
}

resource "azurerm_storage_account" "bootdiag" {
  name = "bootdiagmasters${random_string.resource_group_suffix.result}"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${var.azure_region}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}


