// CoreOS Install Profile
resource "matchbox_profile" "coreos-install" {
  name = "coreos-install"
  kernel = "/assets/coreos/${var.tectonic_coreos_version}/coreos_production_pxe.vmlinuz"
  initrd = [
    "/assets/coreos/${var.tectonic_coreos_version}/coreos_production_pxe_image.cpio.gz"
  ]
  args = [
    "coreos.config.url=${var.tectonic_matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
    "coreos.first_boot=yes",
    "console=tty0",
    "console=ttyS0"
  ]
  container_linux_config = "${file("${path.module}/cl/coreos-install.yaml.tmpl")}"
}

// Self-hosted Kubernetes Controller profile
resource "matchbox_profile" "tectonic-controller" {
  name = "bootkube-controller"
  container_linux_config = "${file("${path.module}/cl/bootkube-controller.yaml.tmpl")}"
}

// Self-hosted Kubernetes Worker profile
resource "matchbox_profile" "tectonic-worker" {
  name = "bootkube-worker"
  container_linux_config = "${file("${path.module}/cl/bootkube-worker.yaml.tmpl")}"
}
