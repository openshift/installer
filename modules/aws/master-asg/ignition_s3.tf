data "ignition_config" "s3" {
  append {
    source = "http://${var.cluster_name}-ncg.${var.base_domain}/ignition?profile=master"
  }

  files = ["${data.ignition_file.kubeconfig.id}"]
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/etc/kubernetes/kubeconfig"
  mode       = 0644

  content {
    content = "${var.kubeconfig_content}"
  }
}
