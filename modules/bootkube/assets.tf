# Self-hosted manifests (resources/generated/manifests/)
resource "template_dir" "bootkube" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "./generated/manifests"

  vars {
    tectonic_network_operator_image   = "${var.container_images["tectonic_network_operator"]}"
    tectonic_network_operator_version = "${var.versions["tno"]}"

    kco_config = "${indent(4, chomp(data.template_file.kco-config_yaml.rendered))}"

    calico_mtu            = "${var.calico_mtu}"
    cloud_provider_config = "${var.cloud_provider_config}"
    cluster_cidr          = "${var.cluster_cidr}"
    tectonic_networking   = "${var.tectonic_networking}"

    kube_ca_cert       = "${base64encode(var.kube_ca_cert_pem)}"
    apiserver_key      = "${base64encode(var.apiserver_key_pem)}"
    apiserver_cert     = "${base64encode(var.apiserver_cert_pem)}"
    oidc_ca_cert       = "${base64encode(var.oidc_ca_cert)}"
    pull_secret        = "${base64encode(file(var.pull_secret_path))}"
    serviceaccount_pub = "${base64encode(tls_private_key.service_account.public_key_pem)}"
    serviceaccount_key = "${base64encode(tls_private_key.service_account.private_key_pem)}"

    etcd_ca_cert     = "${base64encode(var.etcd_ca_cert_pem)}"
    etcd_client_cert = "${base64encode(var.etcd_client_cert_pem)}"
    etcd_client_key  = "${base64encode(var.etcd_client_key_pem)}"
  }
}

# kubeconfig (resources/generated/auth/kubeconfig)
data "template_file" "kubeconfig" {
  template = "${file("${path.module}/resources/kubeconfig")}"

  vars {
    kube_ca_cert = "${base64encode(var.kube_ca_cert_pem)}"
    kubelet_cert = "${base64encode(var.kubelet_cert_pem)}"
    kubelet_key  = "${base64encode(var.kubelet_key_pem)}"
    server       = "${var.kube_apiserver_url}"
    cluster_name = "${var.cluster_name}"
  }
}

resource "local_file" "kubeconfig" {
  content  = "${data.template_file.kubeconfig.rendered}"
  filename = "./generated/auth/kubeconfig"
}

# kvo-config.yaml (resources/generated/kco-config.yaml)
data "template_file" "kco-config_yaml" {
  template = "${file("${path.module}/resources/kco-config.yaml")}"

  vars {
    cloud_config_path      = "${var.cloud_config_path}"
    cloud_provider_profile = "${var.cloud_provider != "" ? "${var.cloud_provider}" : "metal"}"

    advertise_address = "${var.advertise_address}"
    cluster_cidr      = "${var.cluster_cidr}"

    etcd_servers = "${
      var.etcd_tls_enabled
          ? join(",", formatlist("https://%s:2379", var.etcd_endpoints))
          : join(",", formatlist("http://%s:2379", var.etcd_endpoints))
      }"

    service_cidr = "${var.service_cidr}"

    oidc_client_id      = "${var.oidc_client_id}"
    oidc_groups_claim   = "${var.oidc_groups_claim}"
    oidc_issuer_url     = "${var.oidc_issuer_url}"
    oidc_username_claim = "${var.oidc_username_claim}"

    master_count = "${var.master_count}"
  }
}

resource "local_file" "kco-config_yaml" {
  content  = "${data.template_file.kco-config_yaml.rendered}"
  filename = "./generated/kco-config.yaml"
}

# bootkube.sh (resources/generated/bootkube.sh)
data "template_file" "bootkube_sh" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image           = "${var.container_images["bootkube"]}"
    kube_core_renderer_image = "${var.container_images["kube_core_renderer"]}"
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
