data "template_file" "flannel" {
  template = "${file("${path.module}/resources/manifests/kube-flannel.yaml")}"

  vars {
    flannel_image     = "${var.flannel_image}"
    flannel_cni_image = "${var.flannel_cni_image}"
    cluster_cidr      = "${var.cluster_cidr}"
    host_cni_bin      = "/var/lib/cni/bin"
  }
}

resource "local_file" "flannel" {
  count = "${ var.enabled ? 1 : 0 }"

  content  = "${data.template_file.flannel.rendered}"
  filename = "./generated/net-manifests/kube-flannel.yaml"
}
