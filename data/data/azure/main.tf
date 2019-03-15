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
  elb_backend_pool_id      = "${module.vnet.lb_backend_pool_id}"
  elb_bootstrap_ssh_natrule_id = "${module.vnet.bootstrap_ssh_natrule_id}"
  nsg_id                   = "${module.vnet.master_nsg_id}"
  tags                     = "${local.tags}"
  boot_diag_blob_endpoint  = "${azurerm_storage_account.bootdiag.primary_blob_endpoint}"

}


//rename to vnet
module "vnet" {
  source = "./vnet"
  rg_name =   "${azurerm_resource_group.main.name}"
  cidr_block = "${var.machine_cidr}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.azure_region}"
  dns_label  = "${var.cluster_id}${random_string.resource_group_suffix.result}"
  tags = "${local.tags}"
}

resource "random_string" "resource_group_suffix" {
  length = 5
  upper = false
  special = false
}

resource "azurerm_resource_group" "main" {
  name = "${var.cluster_id}-rg-${random_string.resource_group_suffix.result}"
  location = "${var.azure_region}"
}

resource "azurerm_storage_account" "bootdiag" {
  name = "bootdiagmasters${random_string.resource_group_suffix.result}"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  location                 = "${var.azure_region}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

//ILB


