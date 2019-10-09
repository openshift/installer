locals {
  bootstrap_count = "${var.bootstrap_complete ? 0 : 1}"
}

// Request from phpIPAM a new IP address for the bootstrap node
module "ipam_bootstrap" {
  source = "./modules/ipam"

  // name of virtual machine
  name = "bootstrap"

  // Only have a single bootstrap virtual machine
  // And once bootstrap is complete remove it
  instance_count = "${local.bootstrap_count}"

  // Hostname or IP address of the phpIPAM server
  ipam = "${var.ipam}"

  // Access token for phpIPAM
  ipam_token = "${var.ipam_token}"

  // Subnet where we will request an ip address from phpIPAM
  machine_cidr = "${var.machine_cidr}"

  // If we already assigned addresses return those
  ip_addresses = ["${compact(list(var.bootstrap_ip_address))}"]

  // Full domain of the OpenShift cluster
  cluster_domain = "${var.cluster_domain}"
}

// Request from phpIPAM a new IP addresses for the control-plane nodes
module "ipam_control_plane" {
  source         = "./modules/ipam"
  name           = "controlplane"
  instance_count = "${var.control_plane_count}"
  ipam           = "${var.ipam}"
  ipam_token     = "${var.ipam_token}"
  machine_cidr   = "${var.machine_cidr}"
  ip_addresses   = "${var.control_plane_ip_addresses}"
  cluster_domain = "${var.cluster_domain}"
}

// Request from phpIPAM a new IP addresses for the compute nodes
module "ipam_compute" {
  source         = "./modules/ipam"
  name           = "compute"
  instance_count = "${var.compute_count}"
  ipam           = "${var.ipam}"
  ipam_token     = "${var.ipam_token}"
  machine_cidr   = "${var.machine_cidr}"
  ip_addresses   = "${var.compute_ip_addresses}"
  cluster_domain = "${var.cluster_domain}"
}

// Clones, creates and configures (via ignition) RHCOS virtual machines

module "rhcos_virtual_machines" {
  source = "./modules/rhcos_virtual_machines"

  // VMware vSphere variables
  vsphere_server     = "${var.vsphere_server}"
  vsphere_user       = "${var.vsphere_user}"
  vsphere_password   = "${var.vsphere_password}"
  vsphere_cluster    = "${var.vsphere_cluster}"
  vsphere_datastore  = "${var.vsphere_datastore}"
  vsphere_datacenter = "${var.vsphere_datacenter}"
  vm_network         = "${var.vm_network}"
  vm_template        = "${var.vm_template}"
  vm_dns_addresses   = "${var.vm_dns_addresses}"

  // Virtual Machine type variables
  bootstrap_ip_address    = "${module.ipam_bootstrap.ip_addresses[0]}"
  bootstrap_ignition_path = "${var.bootstrap_ignition_path}"
  bootstrap_count         = "${local.bootstrap_count}"

  control_plane_ip_addresses  = "${module.ipam_control_plane.ip_addresses}"
  control_plane_ignition_path = "${var.control_plane_ignition_path}"
  control_plane_count         = "${var.control_plane_count}"

  compute_ip_addresses  = "${module.ipam_compute.ip_addresses}"
  compute_ignition_path = "${var.compute_ignition_path}"
  compute_count         = "${var.compute_count}"

  // OpenShift variables
  base_domain    = "${var.base_domain}"
  cluster_domain = "${var.cluster_domain}"
  cluster_id     = "${var.cluster_id}"
  machine_cidr   = "${var.machine_cidr}"
}

module "dns" {
  source = "./modules/route53"

  base_domain    = "${var.base_domain}"
  cluster_domain = "${var.cluster_domain}"

  bootstrap_count        = "${local.bootstrap_count}"
  bootstrap_ip_addresses = "${module.ipam_bootstrap.ip_addresses}"

  control_plane_count        = "${var.control_plane_count}"
  control_plane_ip_addresses = "${module.ipam_control_plane.ip_addresses}"

  compute_count        = "${var.compute_count}"
  compute_ip_addresses = "${module.ipam_compute.ip_addresses}"
}
