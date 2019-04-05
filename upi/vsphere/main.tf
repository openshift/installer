provider "vsphere" {
  user                 = "${var.vsphere_user}"
  password             = "${var.vsphere_password}"
  vsphere_server       = "${var.vsphere_server}"
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = "${var.vsphere_datacenter}"
}

module "folder" {
  source = "./folder"

  path          = "${var.cluster_id}"
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
  datastore        = "${var.vsphere_datastore}"
  folder           = "${module.folder.path}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ipam             = "${var.ipam}"
  ipam_token       = "${var.ipam_token}"
  machine_cidr     = "${var.machine_cidr}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "control_plane" {
  source = "./machine"

  name             = "control-plane"
  instance_count   = "${var.control_plane_count}"
  ignition         = "${var.control_plane_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  folder           = "${module.folder.path}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ipam             = "${var.ipam}"
  ipam_token       = "${var.ipam_token}"
  machine_cidr     = "${var.machine_cidr}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "compute" {
  source = "./machine"

  name             = "compute"
  instance_count   = "${var.compute_count}"
  ignition         = "${var.compute_ignition}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  folder           = "${module.folder.path}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ipam             = "${var.ipam}"
  ipam_token       = "${var.ipam_token}"
  machine_cidr     = "${var.machine_cidr}"

  extra_user_names           = ["${var.extra_user_names}"]
  extra_user_password_hashes = ["${var.extra_user_password_hashes}"]
}

module "dns" {
  source = "./route53"

  base_domain       = "${var.base_domain}"
  cluster_domain    = "${var.cluster_domain}"
  bootstrap_ips     = ["${module.bootstrap.ip_addresses}"]
  control_plane_ips = ["${module.control_plane.ip_addresses}"]
  compute_ips       = ["${module.compute.ip_addresses}"]
}
