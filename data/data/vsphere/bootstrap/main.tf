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

resource "vsphere_virtual_machine" "vm_bootstrap" {
  name                        = "${var.cluster_id}-bootstrap"
  resource_pool_id            = var.resource_pool[0].id
  datastore_id                = var.datastore[0].id
  num_cpus                    = var.vsphere_control_planes[0].numCPUs
  num_cores_per_socket        = var.vsphere_control_planes[0].numCoresPerSocket
  memory                      = var.vsphere_control_planes[0].memoryMiB
  guest_id                    = var.template[0].guest_id
  folder                      = trimprefix(var.vsphere_control_planes[0].workspace.folder, "/${var.vsphere_control_planes[0].workspace.datacenter}/vm")
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags
  firmware                    = "efi"

  network_interface {
    network_id = var.template[0].network_interfaces.0.network_id
  }

  disk {
    label = "disk0"
    size  = var.vsphere_control_planes[0].diskGiB

    eagerly_scrub    = var.template[0].disks.0.eagerly_scrub
    thin_provisioned = var.template[0].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[0].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
    "guestinfo.domain"                        = "${var.cluster_domain}"
    "stealclock.enable"                       = "TRUE"
  }

  // Potential issues on destroy if disk type changes
  // underneath terraform.
  lifecycle {
    ignore_changes = [
      disk[0],
    ]
  }
}
