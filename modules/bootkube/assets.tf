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

    kco_config = "${indent(4, chomp(data.template_file.kco-config_yaml.rendered))}"

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

# kvo-config.yaml (resources/generated/kco-config.yaml)
data "template_file" "kco-config_yaml" {
  template = "${file("${path.module}/resources/kco-config.yaml")}"

  vars {
    kube_apiserver_url = "${var.kube_apiserver_url}"

    cloud_config_path      = "${var.cloud_config_path}"
    cloud_provider_profile = "${var.cloud_provider != "" ? "${var.cloud_provider}" : "metal"}"

    advertise_address = "${var.advertise_address}"
    cluster_cidr      = "${var.cluster_cidr}"

    etcd_servers = "${join(",", formatlist("https://%s:2379", var.etcd_endpoints))}"

    service_cidr = "${var.service_cidr}"

    oidc_client_id      = "${var.oidc_client_id}"
    oidc_groups_claim   = "${var.oidc_groups_claim}"
    oidc_issuer_url     = "${var.oidc_issuer_url}"
    oidc_username_claim = "${var.oidc_username_claim}"
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

# TNC
resource "local_file" "tnco_pod_config" {
  content  = "${data.template_file.tnco_config.rendered}"
  filename = "./generated/tnco-config.yaml"
}

data "template_file" "initial_cluster" {
  count    = "${length(var.etcd_endpoints)}"
  template = "${var.etcd_endpoints[count.index]}=https://${var.etcd_endpoints[count.index]}:2380"
}

data "template_file" "tnco_config" {
  template = "${file("${path.module}/resources/tnco-config.yaml")}"

  vars {
    cloud_provider_config     = "${var.cloud_provider_config}"
    kubeconfig_fetch_cmd      = "${var.kubeconfig_fetch_cmd != "" ? "ExecStartPre=${var.kubeconfig_fetch_cmd}" : ""}"
    cluster_dns_ip            = "${var.kube_dns_service_ip}"
    cloud_provider            = "${var.cloud_provider}"
    cluster_name              = "${var.cluster_name}"
    base_domain               = "${var.base_domain}"
    etcd_initial_cluster_list = "${length(var.etcd_endpoints) > 0 ? format("--initial-cluster=%s", join(",", data.template_file.initial_cluster.*.rendered)) : ""}"
  }
}
