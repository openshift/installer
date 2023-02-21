locals {
  description  = "Created By OpenShift Installer"
  vcenter_key  = keys(var.vsphere_vcenters)[0]
  template_map = { for t in var.template : t.name => t }
}

provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}

data "vsphere_datacenter" "datacenter" {
  count = var.master_count
  name  = var.vsphere_control_planes[count.index].workspace.datacenter
}

data "vsphere_resource_pool" "resource_pool" {
  count = var.master_count
  name  = var.vsphere_control_planes[count.index].workspace.resourcePool
}

data "vsphere_datastore" "datastore" {
  count         = var.master_count
  name          = var.vsphere_control_planes[count.index].workspace.datastore
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
}

resource "vsphere_virtual_machine" "vm_master" {
  count = var.master_count

  name                 = "${var.cluster_id}-master-${count.index}"
  resource_pool_id     = data.vsphere_resource_pool.resource_pool[count.index].id
  datastore_id         = data.vsphere_datastore.datastore[count.index].id
  num_cpus             = var.vsphere_control_planes[0].numCPUs
  num_cores_per_socket = var.vsphere_control_planes[0].numCoresPerSocket
  memory               = var.vsphere_control_planes[0].memoryMiB
  folder               = trimprefix(var.vsphere_control_planes[count.index].workspace.folder, "/${var.vsphere_control_planes[count.index].workspace.datacenter}/vm")

  guest_id = local.template_map[var.vsphere_control_planes[count.index].template].guest_id

  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags

  network_interface {
    network_id = local.template_map[var.vsphere_control_planes[count.index].template].network_interfaces.0.network_id
  }

  disk {
    label            = "disk0"
    size             = var.vsphere_control_planes[0].diskGiB
    eagerly_scrub    = local.template_map[var.vsphere_control_planes[count.index].template].disks.0.eagerly_scrub
    thin_provisioned = local.template_map[var.vsphere_control_planes[count.index].template].disks.0.thin_provisioned
  }

  clone {
    template_uuid = local.template_map[var.vsphere_control_planes[count.index].template].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-master-${count.index}"
    "stealclock.enable"                       = "TRUE"
  }
}
