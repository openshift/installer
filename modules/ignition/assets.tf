data "template_file" "kubelet" {
  template = "${file("${path.module}/resources/services/kubelet.service")}"

  vars {
    cloud_provider        = "${var.cloud_provider}"
    cloud_provider_config = "${var.cloud_provider_config != "" ? "--cloud-config=/etc/kubernetes/cloud/config" : ""}"
    cluster_dns_ip        = "${var.kube_dns_service_ip}"
    debug_config          = "${var.kubelet_debug_config}"
    node_label            = "${var.kubelet_node_label}"
    node_taints_param     = "${var.kubelet_node_taints != "" ? "--register-with-taints=${var.kubelet_node_taints}" : ""}"
  }
}

data "ignition_systemd_unit" "kubelet" {
  name    = "kubelet.service"
  enabled = true
  content = "${data.template_file.kubelet.rendered}"
}

data "ignition_file" "registries_config" {
  filesystem = "root"
  mode       = "0644"
  path       = "/etc/containers/registries.conf"

  content {
    content = "${file("${path.module}/resources/files/registries.conf")}"
  }
}
