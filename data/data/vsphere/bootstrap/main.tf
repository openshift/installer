locals {
  description = "Created By OpenShift Installer"
  bootstrap_key = keys(var.vcenter_region_zone_map)[0]
  region = var.vcenter_region_zone_map[local.bootstrap_key].region
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_bootstrap" {
  name             = "${var.cluster_id}-bootstrap"
  resource_pool_id = var.cluster[local.bootstrap_key].resource_pool_id
  datastore_id     = var.datastore[local.bootstrap_key].id
  num_cpus         = 4
  memory           = 16384
  guest_id         = var.template[local.bootstrap_key].guest_id


  folder           = var.folder[local.region]
  enable_disk_uuid = "true"
  annotation       = local.description

  wait_for_guest_net_timeout  = 0
  wait_for_guest_net_routable = false

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
    template_uuid = var.template[local.bootstrap_key].id
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
  }
  tags = var.tags
}

