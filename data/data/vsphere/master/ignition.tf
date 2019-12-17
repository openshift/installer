locals {
  ignition_encoded = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition)}"
}

data "ignition_file" "hostname" {
  count = var.instance_count

  filesystem = "root"
  path       = "/etc/hostname"
  mode       = "420"

  content {
    content = "${var.name}-${count.index}"
  }
}

data "ignition_config" "ign" {
  count = var.instance_count

  append {
    source = local.ignition_encoded
  }

  files = [
    data.ignition_file.hostname[count.index].rendered
  ]
}

