provider "libvirt" {
  uri = var.libvirt_uri
}

provider "ironic" {
  url                = var.ironic_uri
  inspector          = var.inspector_uri
  microversion       = "1.56"
  timeout            = 900
  auth_strategy      = "http_basic"
  ironic_username    = var.ironic_username
  ironic_password    = var.ironic_password
  inspector_username = var.ironic_username
  inspector_password = var.ironic_password
}

resource "libvirt_pool" "bootstrap" {
  name = "${var.cluster_id}-bootstrap"
  type = "dir"
  path = "/var/lib/libvirt/openshift-images/${var.cluster_id}-bootstrap"
}

resource "libvirt_volume" "bootstrap-base" {
  name   = "${var.cluster_id}-bootstrap-base"
  pool   = libvirt_pool.bootstrap.name
  source = var.bootstrap_os_image
}

resource "libvirt_volume" "bootstrap" {
  name           = "${var.cluster_id}-bootstrap"
  pool           = libvirt_pool.bootstrap.name
  base_volume_id = libvirt_volume.bootstrap-base.id
  # Keep this in sync with the main libvirt size in
  # data/data/libvirt/bootstrap/main.tf
  size = "34359738368"
}

resource "libvirt_ignition" "bootstrap" {
  name    = "${var.cluster_id}-bootstrap.ign"
  pool    = libvirt_pool.bootstrap.name
  content = var.ignition_bootstrap
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

  cpu {
    mode = "host-passthrough"
  }

  dynamic "network_interface" {
    for_each = var.bridges
    content {
      bridge = network_interface.value["name"]
      mac    = network_interface.value["mac"]
    }
  }

  graphics {
    type        = "vnc"
    listen_type = "address"
  }
}
