resource "vsphere_virtual_machine" "vm" {
  count = var.instance_count

  name                 = "${var.cluster_id}-${var.name}-${count.index}"
  resource_pool_id     = var.resource_pool
  datastore_id         = var.datastore
  num_cpus             = var.num_cpus
  num_cores_per_socket = var.cores_per_socket
  memory               = var.memory
  guest_id             = var.guest_id
  folder               = var.folder
  enable_disk_uuid     = "true"

  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"

  network_interface {
    network_id = var.network
  }

  disk {
    label            = "disk0"
    size             = var.disk_size
    eagerly_scrub    = var.scrub_disk
    thin_provisioned = var.thin_disk
  }

  clone {
    template_uuid = var.template
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-${var.name}-${count.index}"
  }

  tags = var.tags
}

