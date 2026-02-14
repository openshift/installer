
locals {
  failure_domains     = length(var.failure_domains) == 0 ? [{
        datacenter = var.vsphere_datacenter
        cluster = var.vsphere_cluster
        datastore = var.vsphere_datastore
        network = var.vm_network
        distributed_virtual_switch_uuid = ""
  }] : var.failure_domains

  failure_domain_count = length(local.failure_domains)
  bootstrap_fqdns     = ["bootstrap-0.${var.cluster_domain}"]
  lb_fqdns            = ["lb-0.${var.cluster_domain}"]
  api_lb_fqdns        = formatlist("%s.%s", ["api", "api-int", "*.apps"], var.cluster_domain)
  control_plane_fqdns = [for idx in range(var.control_plane_count) : "${var.cluster_id}-control-plane-${idx}.${var.cluster_domain}"]
  compute_fqdns       = [for idx in range(var.compute_count) : "${var.cluster_id}-compute-${idx}.${var.cluster_domain}"]
  datastores          = [for idx in range(length(local.failure_domains)) : local.failure_domains[idx]["datastore"]]
  datacenters         = [for idx in range(length(local.failure_domains)) : local.failure_domains[idx]["datacenter"]]
  datacenters_distinct = distinct([for idx in range(length(local.failure_domains)) : local.failure_domains[idx]["datacenter"]])
  clusters            = [for idx in range(length(local.failure_domains)) : local.failure_domains[idx]["cluster"]]
  networks            = [for idx in range(length(local.failure_domains)) : local.failure_domains[idx]["cluster"]]
  folders             = [for idx in range(length(local.datacenters)) : "/${local.datacenters[idx]}/vm/${var.cluster_id}"]
}

provider "vsphere" {
  user                 = var.vsphere_user
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_server
  allow_unverified_ssl = true
}

data "vsphere_datacenter" "dc" {
   count = length(local.datacenters_distinct)
   name = local.datacenters_distinct[count.index]
}

data "vsphere_compute_cluster" "compute_cluster" {
   count = length(local.failure_domains)
   name = local.clusters[count.index]
   datacenter_id = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.datacenters[count.index])].id
}
#
data "vsphere_datastore" "datastore" {
   count = length(local.failure_domains)
   name = local.datastores[count.index]
   datacenter_id = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.datacenters[count.index])].id
}

#
data "vsphere_network" "network" {
  count = length(local.failure_domains)
  name          = local.failure_domains[count.index]["network"]
  datacenter_id = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index]["datacenter"])].id
  distributed_virtual_switch_uuid = local.failure_domains[count.index]["distributed_virtual_switch_uuid"]
}

data "vsphere_virtual_machine" "template" {
  count = length(local.datacenters_distinct)
  name          = var.vm_template
  datacenter_id = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.datacenters_distinct[count.index])].id
}

resource "vsphere_resource_pool" "resource_pool" {
  count                   = length(data.vsphere_compute_cluster.compute_cluster)
  name                    = var.cluster_id
  parent_resource_pool_id = data.vsphere_compute_cluster.compute_cluster[count.index].resource_pool_id
}

resource "vsphere_folder" "folder" {
  count = length(local.datacenters_distinct)
  path          = var.cluster_id
  type          = "vm"
  datacenter_id = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.datacenters_distinct[count.index])].id
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
   vmname                = element(split(".", local.lb_fqdns[0]), 0)
   ipaddress             = module.ipam_lb.ip_addresses[0]
   ignition               = module.lb.ignition
   resource_pool_id      = vsphere_resource_pool.resource_pool[0].id
   datastore_id          = data.vsphere_datastore.datastore[0].id
   datacenter_id         = data.vsphere_datacenter.dc[0].id
   network_id            = data.vsphere_network.network[0].id
   folder_id             = vsphere_folder.folder[0].path
   guest_id              = data.vsphere_virtual_machine.template[0].guest_id
   template_uuid         = data.vsphere_virtual_machine.template[0].id
   disk_thin_provisioned = data.vsphere_virtual_machine.template[0].disks[0].thin_provisioned
   cluster_domain = var.cluster_domain
   machine_cidr   = var.machine_cidr
   num_cpus      = 2
   memory        = 2096
   dns_addresses = var.vm_dns_addresses
 }

module "bootstrap" {
  source = "./vm"

  ignition = file(var.bootstrap_ignition_path)

  vmname                = element(split(".", local.bootstrap_fqdns[0]), 0)
  ipaddress             = module.ipam_bootstrap.ip_addresses[0]
  resource_pool_id      = vsphere_resource_pool.resource_pool[0].id
  datastore_id          = data.vsphere_datastore.datastore[0].id
  datacenter_id         = data.vsphere_datacenter.dc[0].id
  network_id            = data.vsphere_network.network[0].id
  folder_id             = vsphere_folder.folder[0].path
  guest_id              = data.vsphere_virtual_machine.template[0].guest_id
  template_uuid         = data.vsphere_virtual_machine.template[0].id
  disk_thin_provisioned = data.vsphere_virtual_machine.template[0].disks[0].thin_provisioned

  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr

  num_cpus      = 2
  memory        = 8192
  dns_addresses = var.vm_dns_addresses
}

module "control_plane_vm" {
  count = length(module.control_plane_a_records.fqdns)
  source = "./vm"
  // Using the output from control_plane_a_records
  // is on purpose. I want the A records to be created before
  // the virtual machines which gives additional time to
  // replicate the records.


  vmname = element(split(".", module.control_plane_a_records.fqdns[count.index]), 0)
  ipaddress = module.ipam_control_plane.ip_addresses[count.index]
  ignition = file(var.control_plane_ignition_path)
  resource_pool_id      = vsphere_resource_pool.resource_pool[count.index % local.failure_domain_count].id
  datastore_id          = data.vsphere_datastore.datastore[count.index % local.failure_domain_count].id
  datacenter_id         = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].id
  network_id            = data.vsphere_network.network[count.index % local.failure_domain_count].id
  folder_id             = vsphere_folder.folder[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].path
  guest_id              = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].guest_id
  template_uuid         = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count ]["datacenter"])].id
  disk_thin_provisioned = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].disks[0].thin_provisioned
  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr
  num_cpus      = var.control_plane_num_cpus
  memory        = var.control_plane_memory
  dns_addresses = var.vm_dns_addresses
}
module "compute_vm" {
  count = length(module.compute_a_records.fqdns)
  source = "./vm"
  ignition = file(var.compute_ignition_path)
  vmname = element(split(".", module.compute_a_records.fqdns[count.index]), 0)
  ipaddress = module.ipam_compute.ip_addresses[count.index]

  resource_pool_id      = vsphere_resource_pool.resource_pool[count.index % local.failure_domain_count].id
  datastore_id          = data.vsphere_datastore.datastore[count.index % local.failure_domain_count].id
  datacenter_id         = data.vsphere_datacenter.dc[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].id
  network_id            = data.vsphere_network.network[count.index % local.failure_domain_count].id
  folder_id             = vsphere_folder.folder[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].path
  guest_id              = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].guest_id
  template_uuid         = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].id
  disk_thin_provisioned = data.vsphere_virtual_machine.template[index(data.vsphere_datacenter.dc.*.name, local.failure_domains[count.index % local.failure_domain_count]["datacenter"])].disks[0].thin_provisioned
  cluster_domain = var.cluster_domain
  machine_cidr   = var.machine_cidr
  num_cpus      = var.compute_num_cpus
  memory        = var.compute_memory
  dns_addresses = var.vm_dns_addresses
}
