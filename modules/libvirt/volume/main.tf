# Create a QOW volume from the downloaded path
resource "libvirt_volume" "coreos_base" {
  name   = "coreos_base"
  source = "file://${var.coreos_qow_path}"
}
