resource "matchbox_profile" "bootstrap" {
  name   = "${var.cluster_id}-bootstrap"
  kernel = var.pxe_kernel

  initrd = [
    var.pxe_initrd,
  ]

  args = concat(
    var.pxe_kernel_args,
    ["coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?cluster_id=${var.cluster_id}&role=bootstrap"],
  )

  raw_ignition = var.igntion_config_content
}

resource "matchbox_group" "bootstrap" {
  name    = "${var.cluster_id}-bootstrap"
  profile = matchbox_profile.bootstrap.name

  selector = {
    cluster_id = var.cluster_id
    role       = "bootstrap"
  }
}
