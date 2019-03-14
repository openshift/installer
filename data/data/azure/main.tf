locals {
  tags = "${merge(map(
    "kubernetes.io_cluster.${var.cluster_id}", "owned"
  ), var.azure_extra_tags)}"
}

provider "azurerm" {
   # version = "~>1.5"
}

data "local_file" "ignition" {
    filename = "${var.ignition_bootstrap}"
}

//TODO :
module "bootstrap" {
  source = "./bootstrap"
  resource_group_name      = "${azurerm_resource_group.main.name}"
  region                   = "${var.azure_region}"
  
  vm_size                  = "${var.azure_bootstrap_vm_type}"
  cluster_id               = "${var.cluster_id}"
  ignition                 = "${var.ignition_bootstrap}"
  subnet_id                = "${module.vnet.public_subnet_id}"
  elb_backend_pool_id      = "${module.vnet.lb_backend_pool_id}"
  nsg_id                   = "${module.vnet.master_nsg_id}"
  tags = "${local.tags}"
}


//rename to vnet
module "vnet" {
  source = "./vnet"
  rg_name =   "${azurerm_resource_group.main.name}"
  cidr_block = "${var.machine_cidr}"
  cluster_id = "${var.cluster_id}"
  region     = "${var.azure_region}"

  tags = "${local.tags}"
}

resource "random_string" "resource_group_suffix" {
  length = 5
  upper = false
  special = false
}

resource "azurerm_resource_group" "main" {
  name = "${var.cluster_id}-${random_string.resource_group_suffix.result}-rg"
  location = "${var.azure_region}"
}

//INBOUND NAT SSH on 1st master + NAT rule for each master-> see aks-engine
//Standard LB
//Standard PIP
//ILB


