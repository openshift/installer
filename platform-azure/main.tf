provider "azurerm" {}

resource "ignition_user" "core" {
  name = "core"

  ssh_authorized_keys = [
    "${file(var.tectonic_ssh_key)}",
  ]
}

resource "azurerm_resource_group" "tectonic_azure_cluster_resource_group" {
   location = "${var.tectonic_region}"
   name = "tectonic-cluster-${var.tectonic_cluster_name}-group"
}

//module "etcd" {
//  source = "./etcd"
//}

module "master" {
  source = "./master"

  location = "${var.tectonic_azure_location}"
  resource_group_name = "${azurerm_resource_group.tectonic_azure_cluster_resource_group.name}"
  image_reference = "${var.tectonic_azure_image_reference}"
  vm_size = "${var.tectonic_azure_vm_size}"

  kubelet_version = "${var.tectonic_kubelet_version}"
  count = "${var.tectonic_master_count}"
  base_domain = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"
  ignition_user_id = "${ignition_user.core.id}"
  ssh_key = "${var.tectonic_ssh_key}"
}

module "dns" {
  source = "./dns"

  master_ip_addresses  = "${module.master.ip_address}"

  base_domain = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  location = "${var.tectonic_azure_location}"
  resource_group_name = "${var.tectonic_azure_dns_resource_group}"

// TODO etcd list
// TODO worker list
}
