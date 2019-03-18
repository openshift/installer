locals {
  network_prefix = "${element(split("/", var.machine_cidr), 1)}"
  gateway        = "${cidrhost(var.machine_cidr,1)}"
}

provider "vsphere" {
  user                 = "${var.vsphere_user}"
  password             = "${var.vsphere_password}"
  vsphere_server       = "${var.vsphere_server}"
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

data "vsphere_compute_cluster" "compute_cluster" {
  name          = "${var.vsphere_cluster}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_datastore" "datastore" {
  name          = "${var.vsphere_datastore}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_network" "network" {
  name          = "${var.vm_network}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

data "vsphere_virtual_machine" "template" {
  name          = "${var.vm_template}"
  datacenter_id = "${data.vsphere_datacenter.dc.id}"
}

module "resource_pool" {
  source = "./resource_pool"

  name            = "${var.cluster_id}"
  datacenter_id   = "${data.vsphere_datacenter.dc.id}"
  vsphere_cluster = "${var.vsphere_cluster}"
}

module "bootstrap" {
  source = "./machine"

  name             = "bootstrap"
  instance_count   = "${var.bootstrap_complete ? 0 : 1}"
  ignition_url     = "${var.bootstrap_ignition_url}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  network_id       = "${data.vsphere_network.network.id}"
  vm_template_id   = "${data.vsphere_virtual_machine.template.id}"
  cluster_domain   = "${var.cluster_domain}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "control_plane" {
  source = "./machine"

  name             = "control-plane"
  instance_count   = "${var.control_plane_instance_count}"
  ignition         = "${var.control_plane_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  network_id       = "${data.vsphere_network.network.id}"
  vm_template_id   = "${data.vsphere_virtual_machine.template.id}"
  cluster_domain   = "${var.cluster_domain}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "compute" {
  source = "./machine"

  name             = "compute"
  instance_count   = "${var.compute_instance_count}"
  ignition         = "${var.compute_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore_id     = "${data.vsphere_datastore.datastore.id}"
  network_id       = "${data.vsphere_network.network.id}"
  vm_template_id   = "${data.vsphere_virtual_machine.template.id}"
  cluster_domain   = "${var.cluster_domain}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "dns" {
  source = "./route53"

  base_domain       = "${var.base_domain}"
  cluster_domain    = "${var.cluster_domain}"
  bootstrap_ip      = "${var.bootstrap_complete ? "" : var.bootstrap_ip}"
  control_plane_ips = "${var.control_plane_ips}"
}
