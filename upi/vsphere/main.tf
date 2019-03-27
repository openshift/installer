provider "vsphere" {
  user                 = "${var.vsphere_user}"
  password             = "${var.vsphere_password}"
  vsphere_server       = "${var.vsphere_server}"
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

module "network" {
  source = "./network"

  # TODO Pass in and check the base domain and cluster domain to look for A records
  machine_cidr   = "${var.machine_cidr}"
  cluster_domain = "${var.cluster_domain}"
  master_count   = "${var.control_plane_instance_count}"
  worker_count   = "${var.compute_instance_count}"
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
  machine_cidr     = "${var.machine_cidr}"
  ips              = "${module.network.bootstrap_ip}"
  instance_count   = "${var.bootstrap_complete ? 0 : 1}"
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
  machine_cidr     = "${var.machine_cidr}"
  ips              = "${module.network.master_ips}"
  instance_count   = "${var.control_plane_instance_count}"
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
  machine_cidr     = "${var.machine_cidr}"
  ips              = "${module.network.worker_ips}"
  instance_count   = "${var.compute_instance_count}"
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

  base_domain       = "${var.base_domain}"
  cluster_domain    = "${var.cluster_domain}"
  bootstrap_ip      = "${module.network.bootstrap_ip}"
  control_plane_ips = "${module.network.master_ips}"
}
