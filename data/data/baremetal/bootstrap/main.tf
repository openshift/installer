resource "libvirt_pool" "bootstrap" {
  name = "${var.cluster_id}-bootstrap"
  type = "dir"
  path = "/var/lib/libvirt/openshift-images/${var.cluster_id}-bootstrap"
}

resource "libvirt_volume" "bootstrap" {
  name   = "${var.cluster_id}-bootstrap"
  pool   = libvirt_pool.bootstrap.name
  source = var.image
}

resource "libvirt_ignition" "bootstrap" {
  name    = "${var.cluster_id}-bootstrap.ign"
  content = var.ignition
}

resource "libvirt_domain" "bootstrap" {
  name = "${var.cluster_id}-bootstrap"

  memory = "6144"

  vcpu = "4"

  coreos_ignition = libvirt_ignition.bootstrap.id

  disk {
    volume_id = libvirt_volume.bootstrap.id
  }

  console {
    type        = "pty"
    target_port = 0
  }

  cpu = {
    mode = "host-passthrough"
  }

  dynamic "network_interface" {
    for_each = var.bridges
    content {
      bridge = network_interface.value
    }
  }
}

