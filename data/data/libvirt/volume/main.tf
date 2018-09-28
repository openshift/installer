resource "libvirt_volume" "coreos_base" {
  name   = "coreos_base"
  source = "${var.image}"
}
