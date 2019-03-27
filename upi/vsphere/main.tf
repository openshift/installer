locals {
  bootstrap_needed = "${var.step < 3}"
  nodes_needed     = "${var.step > 1}"
  dns_needed       = "${var.step >= 2}"

  control_plane_count = "${local.nodes_needed ? var.control_plane_instance_count: 0}"
  compute_count       = "${local.nodes_needed ? var.compute_instance_count: 0}"
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

module "resource_pool" {
  source = "./resource_pool"

  name            = "${var.cluster_id}"
  datacenter_id   = "${data.vsphere_datacenter.dc.id}"
  vsphere_cluster = "${var.vsphere_cluster}"
}

module "bootstrap" {
  source = "./machine"

  name             = "bootstrap"
  instance_count   = "${local.bootstrap_needed ? 1 : 0}"
  ignition_url     = "${var.bootstrap_ignition_url}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  cluster_id       = "${var.cluster_id}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "control_plane" {
  source = "./machine"

  name             = "control-plane"
  instance_count   = "${local.nodes_needed ? var.control_plane_instance_count: 0}"
  ignition         = "${var.control_plane_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  cluster_id       = "${var.cluster_id}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "compute" {
  source = "./machine"

  name             = "compute"
  instance_count   = "${local.nodes_needed ? var.compute_instance_count: 0}"
  ignition         = "${var.compute_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  cluster_id       = "${var.cluster_id}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "dns" {
  source = "./route53"

  base_domain                  = "${var.base_domain}"
  cluster_domain               = "${var.cluster_domain}"
  bootstrap_ip                 = ["${module.bootstrap.ip_addresses}"]
  control_plane_instance_count = "${local.dns_needed ? var.control_plane_instance_count: 0}"
  control_plane_ips            = ["${module.control_plane.ip_addresses}"]
  compute_instance_count       = "${local.dns_needed ? var.compute_instance_count: 0}"
  compute_ips                  = ["${module.compute.ip_addresses}"]
}
