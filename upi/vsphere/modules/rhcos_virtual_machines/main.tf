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
  instance_count   = "${var.bootstrap_count}"
  ignition         = "${file("${var.bootstrap_ignition_path}")}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  datastore        = "${var.vsphere_datastore}"
  folder           = "${module.folder.path}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ip_addresses     = ["${compact(list(var.bootstrap_ip_address))}"]
  machine_cidr     = "${var.machine_cidr}"
  memory           = "8192"
  num_cpu          = "4"
  dns_addresses    = "${var.vm_dns_addresses}"
}

module "control-plane" {
  source = "./machine"

  name             = "controlplane"
  instance_count   = "${var.control_plane_count}"
  ignition         = "${file("${var.control_plane_ignition_path}")}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  folder           = "${module.folder.path}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ip_addresses     = "${var.control_plane_ip_addresses}"
  machine_cidr     = "${var.machine_cidr}"
  memory           = "16384"
  num_cpu          = "4"
  dns_addresses    = "${var.vm_dns_addresses}"
}

module "compute" {
  source = "./machine"

  name             = "compute"
  instance_count   = "${var.compute_count}"
  ignition         = "${file("${var.compute_ignition_path}")}"
  resource_pool_id = "${module.resource_pool.pool_id}"
  folder           = "${module.folder.path}"
  datastore        = "${var.vsphere_datastore}"
  network          = "${var.vm_network}"
  datacenter_id    = "${data.vsphere_datacenter.dc.id}"
  template         = "${var.vm_template}"
  cluster_domain   = "${var.cluster_domain}"
  ip_addresses     = "${var.compute_ip_addresses}"
  machine_cidr     = "${var.machine_cidr}"
  memory           = "8192"
  num_cpu          = "4"
  dns_addresses    = "${var.vm_dns_addresses}"
}
