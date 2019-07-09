resource "libvirt_volume" "coreos_base" {
  name   = "${var.cluster_id}-base"
  source = var.image
  pool   = var.pool
}
