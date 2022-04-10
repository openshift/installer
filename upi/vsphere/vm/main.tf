resource "vsphere_virtual_machine" "vm" {
  for_each = var.hostnames_ip_addresses

  name = element(split(".", each.key), 0)

  resource_pool_id = var.resource_pool_id
  datastore_id     = var.datastore_id
  num_cpus         = var.num_cpus
  memory           = var.memory
  guest_id         = var.guest_id
  folder           = var.folder_id
  enable_disk_uuid = "true"

  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"

  network_interface {
    network_id = var.network_id
  }

  disk {
    label            = "disk0"
    size             = 60
    thin_provisioned = var.disk_thin_provisioned
  }

  clone {
    template_uuid = var.template_uuid
  }
  cdrom {
    datastore_id = var.datastore_id
    path         = "ISO/rhcos-410.84.202201251210-0-live.x86_64_last1.iso"
  }


  extra_config = {
    "guestinfo.ignition.config.data"           = base64encode(var.ignition)
    "guestinfo.ignition.config.data.encoding"  = "base64"
    "guestinfo.afterburn.initrd.network-kargs" = "ip=${each.value}::${cidrhost(var.machine_cidr, 1)}:${cidrnetmask(var.machine_cidr)}:${element(split(".", each.key), 0)}:ens192:none:${join(":", var.dns_addresses)}"
  }
}

