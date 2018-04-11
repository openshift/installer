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

# Self-hosted manifests (resources/generated/manifests/)
resource "template_dir" "bootkube" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "./generated/manifests"

  vars {
    tectonic_network_operator_image = "${var.container_images["tectonic_network_operator"]}"
    tnc_operator_image              = "${var.container_images["tnc_operator"]}"

    calico_mtu            = "${var.calico_mtu}"
    cloud_provider_config = "${var.cloud_provider_config}"
    cluster_cidr          = "${var.cluster_cidr}"
    tectonic_networking   = "${var.tectonic_networking}"

    root_ca_cert                   = "${base64encode(var.root_ca_cert_pem)}"
    aggregator_ca_cert             = "${base64encode(var.aggregator_ca_cert_pem)}"
    kube_ca_cert                   = "${base64encode(var.kube_ca_cert_pem)}"
    kube_ca_key                    = "${base64encode(var.kube_ca_key_pem)}"
    kubelet_bootstrap_token_id     = "${random_string.kubelet_bootstrap_token_id.result}"
    kubelet_bootstrap_token_secret = "${random_string.kubelet_bootstrap_token_secret.result}"
    apiserver_key                  = "${base64encode(var.apiserver_key_pem)}"
    apiserver_cert                 = "${base64encode(var.apiserver_cert_pem)}"
    apiserver_proxy_key            = "${base64encode(var.apiserver_proxy_key_pem)}"
    apiserver_proxy_cert           = "${base64encode(var.apiserver_proxy_cert_pem)}"
    oidc_ca_cert                   = "${base64encode(var.oidc_ca_cert)}"
    pull_secret                    = "${base64encode(file(var.pull_secret_path))}"
    serviceaccount_pub             = "${base64encode(tls_private_key.service_account.public_key_pem)}"
    serviceaccount_key             = "${base64encode(tls_private_key.service_account.private_key_pem)}"
    kube_dns_service_ip            = "${cidrhost(var.service_cidr, 10)}"

    etcd_ca_cert     = "${base64encode(var.etcd_ca_cert_pem)}"
    etcd_client_cert = "${base64encode(var.etcd_client_cert_pem)}"
    etcd_client_key  = "${base64encode(var.etcd_client_key_pem)}"
  }
}

# kubeconfig (resources/generated/auth/kubeconfig)
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

resource "local_file" "kubeconfig" {
  content  = "${data.template_file.kubeconfig.rendered}"
  filename = "./generated/auth/kubeconfig"
}

# kubeconfig-kubelet (resources/generated/auth/kubeconfig-kubelet)
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

resource "local_file" "kubeconfig-kubelet" {
  content  = "${data.template_file.kubeconfig-kubelet.rendered}"
  filename = "./generated/auth/kubeconfig-kubelet"
}

# bootkube.sh (resources/generated/bootkube.sh)
data "template_file" "bootkube_sh" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image           = "${var.container_images["bootkube"]}"
    kube_core_renderer_image = "${var.container_images["kube_core_renderer"]}"
    tnc_operator_image       = "${var.container_images["tnc_operator"]}"
  }
}

resource "local_file" "bootkube_sh" {
  content  = "${data.template_file.bootkube_sh.rendered}"
  filename = "./generated/bootkube.sh"
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

data "template_file" "initial_cluster" {
  count    = "${length(var.etcd_endpoints)}"
  template = "${var.etcd_endpoints[count.index]}=https://${var.etcd_endpoints[count.index]}:2380"
}
