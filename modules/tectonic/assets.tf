# Unique Cluster ID (uuid)
resource "random_id" "cluster_id" {
  byte_length = 16
}

# Kubernetes Manifests (resources/generated/manifests/)
resource "template_dir" "tectonic" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "./generated/tectonic"

  vars {
    addon_resizer_image                   = "${var.container_images["addon_resizer"]}"
    config_reload_image                   = "${var.container_images["config_reload"]}"
    console_image                         = "${var.container_images["console"]}"
    error_server_image                    = "${var.container_images["error_server"]}"
    heapster_image                        = "${var.container_images["heapster"]}"
    identity_image                        = "${var.container_images["identity"]}"
    ingress_controller_image              = "${var.container_images["ingress_controller"]}"
    container_linux_update_operator_image = "${var.container_images["container_linux_update_operator"]}"
    kube_version_operator_image           = "${var.container_images["kube_version_operator"]}"
    node_agent_image                      = "${var.container_images["node_agent"]}"
    node_exporter_image                   = "${var.container_images["node_exporter"]}"
    kube_state_metrics_image              = "${var.container_images["kube_state_metrics"]}"
    prometheus_operator_image             = "${var.container_images["prometheus_operator"]}"
    stats_emitter_image                   = "${var.container_images["stats_emitter"]}"
    stats_extender_image                  = "${var.container_images["stats_extender"]}"
    tectonic_channel_operator_image       = "${var.container_images["tectonic_channel_operator"]}"
    tectonic_prometheus_operator_image    = "${var.container_images["tectonic_prometheus_operator"]}"

    kubernetes_version = "${var.versions["kubernetes"]}"
    monitoring_version = "${var.versions["monitoring"]}"
    prometheus_version = "${var.versions["prometheus"]}"
    tectonic_version   = "${var.versions["tectonic"]}"
    etcd_version       = "${var.versions["etcd"]}"

    etcd_cluster_size = "${var.master_count > 2 ? 3 : 1}"

    license     = "${base64encode(file(var.license_path))}"
    pull_secret = "${base64encode(file(var.pull_secret_path))}"
    ca_cert     = "${base64encode(var.ca_cert)}"

    update_server  = "${var.update_server}"
    update_channel = "${var.update_channel}"
    update_app_id  = "${var.update_app_id}"

    admin_user_id       = "${random_id.admin_user_id.b64}"
    admin_email         = "${var.admin_email}"
    admin_password_hash = "${var.admin_password_hash}"

    console_base_address = "${var.base_address}"
    console_client_id    = "${var.console_client_id}"
    console_secret       = "${random_id.console_secret.b64}"
    console_callback     = "https://${var.base_address}/auth/callback"

    ingress_kind     = "${var.ingress_kind}"
    ingress_tls_cert = "${base64encode(tls_locally_signed_cert.ingress.cert_pem)}"
    ingress_tls_key  = "${base64encode(tls_private_key.ingress.private_key_pem)}"

    identity_server_tls_cert = "${base64encode(tls_locally_signed_cert.identity-server.cert_pem)}"
    identity_server_tls_key  = "${base64encode(tls_private_key.identity-server.private_key_pem)}"
    identity_client_tls_cert = "${base64encode(tls_locally_signed_cert.identity-client.cert_pem)}"
    identity_client_tls_key  = "${base64encode(tls_private_key.identity-client.private_key_pem)}"

    kubectl_client_id = "${var.kubectl_client_id}"
    kubectl_secret    = "${random_id.kubectl_secret.b64}"

    kube_apiserver_url = "${var.kube_apiserver_url}"
    oidc_issuer_url    = "https://${var.base_address}/identity"

    # TODO: We could also patch https://www.terraform.io/docs/providers/random/ to add an UUID resource.
    cluster_id = "${format("%s-%s-%s-%s-%s", substr(random_id.cluster_id.hex, 0, 8), substr(random_id.cluster_id.hex, 8, 4), substr(random_id.cluster_id.hex, 12, 4), substr(random_id.cluster_id.hex, 16, 4), substr(random_id.cluster_id.hex, 20, 12))}"

    platform                 = "${var.platform}"
    certificates_strategy    = "${var.ca_generated == "true" ? "installerGeneratedCA" : "userProvidedCA"}"
    identity_api_service     = "${var.identity_api_service}"
    tectonic_updater_enabled = "${var.experimental ? "true" : "false"}"
  }
}

# tectonic.sh (resources/generated/tectonic.sh)
data "template_file" "tectonic" {
  template = "${file("${path.module}/resources/tectonic.sh")}"

  vars {
    ingress_kind = "${var.ingress_kind}"
  }
}

resource "local_file" "tectonic" {
  content  = "${data.template_file.tectonic.rendered}"
  filename = "./generated/tectonic.sh"
}

# tectonic.sh (resources/generated/tectonic-rkt.sh)
data "template_file" "tectonic-rkt" {
  template = "${file("${path.module}/resources/tectonic-rkt.sh")}"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
    experimental    = "${var.experimental ? "true" : "false"}"
  }
}

resource "local_file" "tectonic-rkt" {
  content  = "${data.template_file.tectonic-rkt.rendered}"
  filename = "./generated/tectonic-rkt.sh"
}

# tectonic.service (available as output variable)
data "template_file" "tectonic_service" {
  template = "${file("${path.module}/resources/tectonic.service")}"
}
