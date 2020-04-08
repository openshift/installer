locals {
  bootstrap_fqdns     = ["bootstrap-0.${var.cluster_domain}"]
  lb_fqdns            = ["lb-0.${var.cluster_domain}"]
  api_lb_fqdns        = formatlist("%s.%s", ["api", "api-int", "*.apps"], var.cluster_domain)
  control_plane_fqdns = [for idx in range(var.control_plane_count) : "control-plane-${idx}.${var.cluster_domain}"]
  compute_fqdns       = [for idx in range(var.compute_count) : "compute-${idx}.${var.cluster_domain}"]
}

provider "vsphere" {
  user                 = var.vsphere_user
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_server
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
  name = var.vsphere_datacenter
}

data "vsphere_compute_cluster" "compute_cluster" {
  name          = var.vsphere_cluster
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_datastore" "datastore" {
  name          = var.vsphere_datastore
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network" {
  name          = var.vm_network
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "template" {
  name          = var.vm_template
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_resource_pool" "resource_pool" {
  name                    = var.cluster_id
  parent_resource_pool_id = data.vsphere_compute_cluster.compute_cluster.resource_pool_id
}

resource "vsphere_folder" "folder" {
  path          = var.cluster_id
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.dc.id
}

// Request from phpIPAM a new IP address for the bootstrap node
module "ipam_bootstrap" {
  source = "./ipam"

  // The hostname that will be added to phpIPAM when requesting an ip address
  hostnames = local.bootstrap_fqdns

  // Hostname or IP address of the phpIPAM server
  ipam = var.ipam

  // Access token for phpIPAM
  ipam_token = var.ipam_token

  // Subnet where we will request an ip address from phpIPAM
  machine_cidr = var.machine_cidr

  static_ip_addresses = var.bootstrap_ip_address == "" ? [] : [var.bootstrap_ip_address]

}

// Request from phpIPAM a new IP addresses for the control-plane nodes
module "ipam_control_plane" {
  source              = "./ipam"
  hostnames           = local.control_plane_fqdns
  ipam                = var.ipam
  ipam_token          = var.ipam_token
  machine_cidr        = var.machine_cidr
  static_ip_addresses = var.control_plane_ip_addresses
}

// Request from phpIPAM a new IP addresses for the compute nodes
module "ipam_compute" {
  source              = "./ipam"
  hostnames           = local.compute_fqdns
  ipam                = var.ipam
  ipam_token          = var.ipam_token
  machine_cidr        = var.machine_cidr
  static_ip_addresses = var.compute_ip_addresses
}

// Request from phpIPAM a new IP addresses for the load balancer nodes
module "ipam_lb" {
  source              = "./ipam"
  hostnames           = local.lb_fqdns
  ipam                = var.ipam
  ipam_token          = var.ipam_token
  machine_cidr        = var.machine_cidr
  static_ip_addresses = var.lb_ip_address == "" ? [] : [var.lb_ip_address]
}

module "lb" {
  source        = "./lb"
  lb_ip_address = module.ipam_lb.ip_addresses[0]

  api_backend_addresses = flatten([
    module.ipam_bootstrap.ip_addresses[0],
    module.ipam_control_plane.ip_addresses]
  )

  ingress_backend_addresses = module.ipam_compute.ip_addresses
  ssh_public_key_path       = var.ssh_public_key_path
}

module "dns_cluster_domain" {
  source         = "./cluster_domain"
  cluster_domain = var.cluster_domain
  base_domain    = var.base_domain
}

module "lb_a_records" {
  source  = "./host_a_record"
  zone_id = module.dns_cluster_domain.zone_id
  records = zipmap(
    local.api_lb_fqdns,
    [for name in local.api_lb_fqdns : module.ipam_lb.ip_addresses[0]]
  )
}

module "control_plane_a_records" {
  source  = "./host_a_record"
  zone_id = module.dns_cluster_domain.zone_id
  records = zipmap(local.control_plane_fqdns, module.ipam_control_plane.ip_addresses)
}

module "compute_a_records" {
  source  = "./host_a_record"
  zone_id = module.dns_cluster_domain.zone_id
  records = zipmap(local.compute_fqdns, module.ipam_compute.ip_addresses)
}

module "lb_vm" {
  source = "./vm"

  ignition               = module.lb.ignition
  hostnames_ip_addresses = zipmap(local.lb_fqdns, module.ipam_lb.ip_addresses)

  resource_pool_id      = vsphere_resource_pool.resource_pool.id
  datastore_id          = data.vsphere_datastore.datastore.id
  datacenter_id         = data.vsphere_datacenter.dc.id
  network_id            = data.vsphere_network.network.id
  folder_id             = vsphere_folder.folder.path
  guest_id              = data.vsphere_virtual_machine.template.guest_id
  template_uuid         = data.vsphere_virtual_machine.template.id
  disk_thin_provisioned = data.vsphere_virtual_machine.template.disks[0].thin_provisioned

  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr

  num_cpus      = 2
  memory        = 2096
  dns_addresses = var.vm_dns_addresses
}

module "bootstrap" {
  source = "./vm"

  ignition = file(var.bootstrap_ignition_path)

  hostnames_ip_addresses = zipmap(
    local.bootstrap_fqdns,
    module.ipam_bootstrap.ip_addresses
  )

  resource_pool_id      = vsphere_resource_pool.resource_pool.id
  datastore_id          = data.vsphere_datastore.datastore.id
  datacenter_id         = data.vsphere_datacenter.dc.id
  network_id            = data.vsphere_network.network.id
  folder_id             = vsphere_folder.folder.path
  guest_id              = data.vsphere_virtual_machine.template.guest_id
  template_uuid         = data.vsphere_virtual_machine.template.id
  disk_thin_provisioned = data.vsphere_virtual_machine.template.disks[0].thin_provisioned

  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr

  num_cpus      = 2
  memory        = 8192
  dns_addresses = var.vm_dns_addresses
}

module "control_plane_vm" {
  source = "./vm"

  // Using the output from control_plane_a_records
  // is on purpose. I want the A records to be created before
  // the virtual machines which gives additional time to
  // replicate the records.
  hostnames_ip_addresses = zipmap(
    module.control_plane_a_records.fqdns,
    module.ipam_control_plane.ip_addresses
  )

  ignition = file(var.control_plane_ignition_path)

  resource_pool_id      = vsphere_resource_pool.resource_pool.id
  datastore_id          = data.vsphere_datastore.datastore.id
  datacenter_id         = data.vsphere_datacenter.dc.id
  network_id            = data.vsphere_network.network.id
  folder_id             = vsphere_folder.folder.path
  guest_id              = data.vsphere_virtual_machine.template.guest_id
  template_uuid         = data.vsphere_virtual_machine.template.id
  disk_thin_provisioned = data.vsphere_virtual_machine.template.disks[0].thin_provisioned

  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr

  num_cpus      = var.control_plane_num_cpus
  memory        = var.control_plane_memory
  dns_addresses = var.vm_dns_addresses
}

module "compute_vm" {
  source = "./vm"

  hostnames_ip_addresses = zipmap(
    module.compute_a_records.fqdns,
    module.ipam_compute.ip_addresses
  )

  ignition = file(var.compute_ignition_path)

  resource_pool_id      = vsphere_resource_pool.resource_pool.id
  datastore_id          = data.vsphere_datastore.datastore.id
  datacenter_id         = data.vsphere_datacenter.dc.id
  network_id            = data.vsphere_network.network.id
  folder_id             = vsphere_folder.folder.path
  guest_id              = data.vsphere_virtual_machine.template.guest_id
  template_uuid         = data.vsphere_virtual_machine.template.id
  disk_thin_provisioned = data.vsphere_virtual_machine.template.disks[0].thin_provisioned

  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr

  num_cpus      = var.compute_num_cpus
  memory        = var.compute_memory
  dns_addresses = var.vm_dns_addresses
}
