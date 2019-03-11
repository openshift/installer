locals {
  tags = "${merge(map(
    "kubernetes.io/cluster/${var.cluster_id}", "owned"
  ), var.azure_extra_tags)}"
}

provider "azurerm" {
  # whilst the `version` attribute is optional, we recommend pinning to a given version of the Provider
  #version = "=1.22.0"
}


//TODO :
module "bootstrap" {
  source = "./bootstrap"
  resource_group_name      = "${var.azure_rg_name}"
  
  instance_type            = "${var.azure_bootstrap_vm_type}"
  cluster_id               = "${var.cluster_id}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vpc.vnet_id}"
  vpc_security_group_ids   = "${list(module.vpc.master_nsg_id)}"

  tags = "${local.tags}"
}


//rename to vnet
module "vpc" {
  source = "./vpc"
  rg_name =   "${azurerm_resource_group.main.name}"
  //cidr_block = "${var.machine_cidr}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.azure_region}"

  tags = "${local.tags}"
}

resource "azurerm_resource_group" "main" {
  name = "${var.azure_rg_name}"
  location = "${var.azure_region}"
}
