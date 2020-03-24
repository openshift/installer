resource "vsphere_virtual_machine" "vm" {
  name             = "${var.cluster_id}-bootstrap"
  resource_pool_id = var.resource_pool
  datastore_id     = var.datastore
  num_cpus         = 4
  memory           = 16384
  guest_id         = var.guest_id
  folder           = var.folder
  enable_disk_uuid = "true"

  wait_for_guest_net_timeout  = 0
  wait_for_guest_net_routable = false

  network_interface {
    network_id = var.network
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
    "guestinfo.ignition.config.data"          = base64encode(var.ignition)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
  }
  tags = var.tags
}

