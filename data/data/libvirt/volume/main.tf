resource "libvirt_volume" "coreos_base_initial" {
  name   = "${var.cluster_id}-base-initial"
  source = var.image
  pool   = var.pool
}

resource "libvirt_volume" "coreos_base" {
  name   = "${var.cluster_id}-base"
  base_volume_id = "${libvirt_volume.coreos_base_initial.id}"
  pool   = var.pool
  size   = 17179869184
}
