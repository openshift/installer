# Kubelet tls bootstraping id and secret
resource "random_string" "kubelet_bootstrap_token_id" {
  length  = 6
  special = false
  upper   = false
}

resource "random_string" "kubelet_bootstrap_token_secret" {
  length  = 16
  special = false
  upper   = false
}

# kubeconfig (/auth/kubeconfig)
data "template_file" "kubeconfig" {
  template = "${file("${path.module}/resources/kubeconfig")}"

  vars {
    root_ca_cert = "${base64encode(var.root_ca_cert_pem)}"
    admin_cert   = "${base64encode(var.admin_cert_pem)}"
    admin_key    = "${base64encode(var.admin_key_pem)}"
    server       = "${var.kube_apiserver_url}"
    cluster_name = "${var.cluster_name}"
  }
}

data "ignition_file" "kubeconfig" {
  filesystem = "root"
  path       = "/opt/tectonic/auth/kubeconfig"
  mode       = "0600"

  content {
    content = "${data.template_file.kubeconfig.rendered}"
  }
}

# kubeconfig-kubelet 
data "template_file" "kubeconfig-kubelet" {
  template = "${file("${path.module}/resources/kubeconfig-kubelet")}"

  vars {
    root_ca_cert                   = "${base64encode(var.root_ca_cert_pem)}"
    kubelet_bootstrap_token_id     = "${random_string.kubelet_bootstrap_token_id.result}"
    kubelet_bootstrap_token_secret = "${random_string.kubelet_bootstrap_token_secret.result}"
    server                         = "${var.kube_apiserver_url}"
    cluster_name                   = "${var.cluster_name}"
  }
}

data "ignition_file" "kubeconfig-kubelet" {
  filesystem = "root"
  path       = "/opt/tectonic/auth/kubeconfig-kubelet"
  mode       = "0600"

  content {
    content = "${data.template_file.kubeconfig-kubelet.rendered}"
  }
}

# bootkube.sh 
data "template_file" "bootkube_sh" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image           = "${var.container_images["bootkube"]}"
    kube_core_renderer_image = "${var.container_images["kube_core_renderer"]}"
    tnc_operator_image       = "${var.container_images["tnc_operator"]}"
  }
}

data "ignition_file" "bootkube_sh" {
  filesystem = "root"
  path       = "/opt/tectonic/bootkube.sh"
  mode       = "0755"

  content {
    content = "${data.template_file.bootkube_sh.rendered}"
  }
}

# bootkube.service (available as output variable)
data "template_file" "bootkube_service" {
  template = "${file("${path.module}/resources/bootkube.service")}"
}

data "ignition_systemd_unit" "bootkube_service" {
  name    = "bootkube.service"
  enabled = false
  content = "${data.template_file.bootkube_service.rendered}"
}

# bootkube.path (available as output variable)
data "template_file" "bootkube_path_unit" {
  template = "${file("${path.module}/resources/bootkube.path")}"
}

data "ignition_systemd_unit" "bootkube_path_unit" {
  name    = "bootkube.path"
  enabled = true
  content = "${data.template_file.bootkube_path_unit.rendered}"
}
