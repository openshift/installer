data "template_file" "userdata-worker" {
  template = "${file("${path.module}/userdata-worker.yml")}"

  vars {
    kube_config      = "${base64encode(file("${path.root}/../assets/auth/kubeconfig"))}"
    tectonic_version = "${var.tectonic_version}"
  }
}
