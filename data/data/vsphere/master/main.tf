locals {
  description = "Created By OpenShift Installer"
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_master" {
  count = var.master_count

  name                 = "${var.cluster_id}-master-${count.index}"
  resource_pool_id     = var.resource_pool
  datastore_id         = var.datastore
  num_cpus             = var.vsphere_control_plane_num_cpus
  num_cores_per_socket = var.vsphere_control_plane_cores_per_socket
  memory               = var.vsphere_control_plane_memory_mib
  guest_id             = var.guest_id
  folder               = var.folder
  enable_disk_uuid     = "true"
  annotation           = local.description

  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"

  network_interface {
    network_id = var.vsphere_network
  }

  disk {
    label            = "disk0"
    size             = var.vsphere_control_plane_disk_gib
    eagerly_scrub    = var.scrub_disk
    thin_provisioned = var.thin_disk
  }

  clone {
    template_uuid = var.template
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-master-${count.index}"
  }

  tags = var.tags
}

