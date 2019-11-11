resource "libvirt_volume" "coreos_orig" {
  name   = "${var.cluster_id}-orig"
  source = var.image
  pool   = var.pool
}

resource "libvirt_volume" "coreos_base" {
  name           = "${var.cluster_id}-base"
  base_volume_id = libvirt_volume.coreos_orig.id
  pool           = var.pool
  # If you change this, you probably want to change the bootstrap too
  size = 16000000000
}
