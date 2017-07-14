data "template_file" "calico-network-policy" {
  template = "${file("${path.module}/resources/manifests/kube-calico.yaml")}"

  vars {
    kube_apiserver_url = "${var.kube_apiserver_url}"
    cni_version        = "${var.cni_version}"
    log_level          = "${var.log_level}"
    calico_image       = "${var.calico_image}"
    calico_cni_image   = "${var.calico_cni_image}"
    cluster_cidr       = "${var.cluster_cidr}"
    host_cni_bin       = "/var/lib/cni/bin"

    bootkube_id = "${var.bootkube_id}"
  }
}

resource "local_file" "calico-network-policy" {
  count = "${ var.enabled ? 1 : 0 }"

  content  = "${data.template_file.calico-network-policy.rendered}"
  filename = "./generated/manifests/kube-calico.yaml"
}
