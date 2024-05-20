locals {
  description = "Created By OpenShift Installer"
  vcenter_key = keys(var.vsphere_vcenters)[0]
}

provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}


data "vsphere_datacenter" "datacenter" {
  count = 1
  name  = var.vsphere_control_planes[count.index].workspace.datacenter
}

data "vsphere_resource_pool" "resource_pool" {
  count = 1
  name  = var.vsphere_control_planes[count.index].workspace.resourcePool
}

data "vsphere_datastore" "datastore" {
  count         = 1
  name          = var.vsphere_control_planes[count.index].workspace.datastore
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
}

data "vsphere_virtual_machine" "template" {
  count         = 1
  name          = var.vsphere_control_planes[count.index].template
  datacenter_id = data.vsphere_datacenter.datacenter[count.index].id
}

resource "vsphere_virtual_machine" "vm_bootstrap" {
  count = 1

  name                        = "${var.cluster_id}-bootstrap"
  resource_pool_id            = data.vsphere_resource_pool.resource_pool[count.index].id
  datastore_id                = data.vsphere_datastore.datastore[count.index].id
  num_cpus                    = var.vsphere_control_planes[0].numCPUs
  num_cores_per_socket        = var.vsphere_control_planes[0].numCoresPerSocket
  memory                      = var.vsphere_control_planes[0].memoryMiB
  guest_id                    = data.vsphere_virtual_machine.template[count.index].guest_id
  folder                      = trimprefix(var.vsphere_control_planes[0].workspace.folder, "/${var.vsphere_control_planes[0].workspace.datacenter}/vm")
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags
  firmware                    = "efi"

  network_interface {
    network_id = var.vsphere_control_planes[0].network.devices[0].networkName
  }

  disk {
    label = "disk0"
    size  = var.vsphere_control_planes[0].diskGiB

    eagerly_scrub    = data.vsphere_virtual_machine.template[count.index].disks.0.eagerly_scrub
    thin_provisioned = data.vsphere_virtual_machine.template[count.index].disks.0.thin_provisioned
  }

  clone {
    template_uuid = data.vsphere_virtual_machine.template[count.index].uuid
  }

  extra_config = merge(
    {
      "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
      "guestinfo.ignition.config.data.encoding" = "base64"
      "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
      "guestinfo.domain"                        = "${var.cluster_domain}"
      "stealclock.enable"                       = "TRUE"
    },
    length(var.vsphere_bootstrap_network_kargs) > 0 ?
    {
      "guestinfo.afterburn.initrd.network-kargs" = "${var.vsphere_bootstrap_network_kargs}"
    } : {}
  )

  // Potential issues on destroy if disk type changes
  // underneath terraform.
  lifecycle {
    ignore_changes = [
      disk[0],
    ]
  }
}
