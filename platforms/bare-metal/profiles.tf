// Configure the matchbox provider
provider "matchbox" {
  endpoint = "${var.matchbox_rpc_endpoint}"
  client_cert = "${var.matchbox_client_cert}"
  client_key = "${var.matchbox_client_key}"
  ca         = "${var.matchbox_ca}"
}

// CoreOS Install Profile
resource "matchbox_profile" "coreos-install" {
  name = "coreos-install"
  kernel = "/assets/coreos/${var.coreos_version}/coreos_production_pxe.vmlinuz"
  initrd = [
    "/assets/coreos/${var.coreos_version}/coreos_production_pxe_image.cpio.gz"
  ]
  args = [
    "coreos.config.url=${var.matchbox_http_endpoint}/ignition?uuid=$${uuid}&mac=$${mac:hexhyp}",
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
