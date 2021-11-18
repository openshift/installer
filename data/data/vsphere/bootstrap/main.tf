locals {
  description = "Created By OpenShift Installer"
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_bootstrap" {
  name             = "${var.cluster_id}-bootstrap"
  resource_pool_id = var.resource_pool
  datastore_id     = var.datastore
  num_cpus         = 4
  memory           = 16384
  guest_id         = var.guest_id
  folder           = var.folder
  enable_disk_uuid = "true"
  annotation       = local.description

  wait_for_guest_net_timeout  = 0
  wait_for_guest_net_routable = false

  network_interface {
    network_id = var.vsphere_network
  }

  disk {
    label            = "disk0"
    size             = 120
    eagerly_scrub    = var.scrub_disk
    thin_provisioned = var.thin_disk
  }

  clone {
    template_uuid = var.template
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
  }
  tags = var.tags
}

