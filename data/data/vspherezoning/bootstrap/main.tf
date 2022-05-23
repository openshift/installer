locals {
  description   = "Created By OpenShift Installer"
  vcenter_key   = keys(var.vsphere_vcenters)[0]
  bootstrap_key = keys(var.vsphere_control_planes_zone)[0]
}

provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_bootstrap" {
  name = "${var.cluster_id}-boostrap"

  resource_pool_id            = var.resource_pool[local.bootstrap_key].id
  datastore_id                = var.datastore[local.bootstrap_key].id
  num_cores_per_socket        = var.vsphere_control_planes_zone[local.bootstrap_key].numCoresPerSocket
  num_cpus                    = var.vsphere_control_planes_zone[local.bootstrap_key].numCPUs
  memory                      = var.vsphere_control_planes_zone[local.bootstrap_key].memoryMiB
  guest_id                    = var.template[local.bootstrap_key].guest_id
  folder                      = var.vsphere_control_planes_zone[local.bootstrap_key].workspace.folder
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags

  network_interface {
    network_id = var.template[local.bootstrap_key].network_interfaces.0.network_id
  }

  disk {
    label            = "disk0"
    size             = 120
    eagerly_scrub    = var.template[local.bootstrap_key].disks.0.eagerly_scrub
    thin_provisioned = var.template[local.bootstrap_key].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[local.bootstrap_key].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
    "guestinfo.domain"                        = "${var.cluster_domain}"
  }

  lifecycle {
    ignore_changes = [
      disk[0].eagerly_scrub,
    ]
  }
}
