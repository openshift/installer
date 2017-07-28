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
    tectonic_monitoring_auth_image        = "${var.container_images["tectonic_monitoring_auth"]}"
    prometheus_image                      = "${var.container_images["prometheus"]}"
    prometheus_config_reload_image        = "${var.container_images["prometheus_config_reload"]}"
    alertmanager_image                    = "${var.container_images["alertmanager"]}"
    stats_emitter_image                   = "${var.container_images["stats_emitter"]}"
    stats_extender_image                  = "${var.container_images["stats_extender"]}"
    tectonic_channel_operator_image       = "${var.container_images["tectonic_channel_operator"]}"
    tectonic_prometheus_operator_image    = "${var.container_images["tectonic_prometheus_operator"]}"
    tectonic_etcd_operator_image          = "${var.container_images["tectonic_etcd_operator"]}"

    kubernetes_version             = "${var.versions["kubernetes"]}"
    monitoring_version             = "${var.versions["monitoring"]}"
    prometheus_version             = "${var.versions["prometheus"]}"
    alertmanager_version           = "${var.versions["alertmanager"]}"
    tectonic_version               = "${var.versions["tectonic"]}"
    etcd_version                   = "${var.versions["etcd"]}"
    tectonic_etcd_operator_version = "${var.versions["tectonic-etcd"]}"

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

    tectonic_monitoring_auth_cookie_secret = "${base64encode(random_id.tectonic_monitoring_auth_cookie_secret.b64)}"

    alertmanager_external_url = "https://${var.base_address}/alertmanager"
    alertmanager_callback     = "https://${var.base_address}/alertmanager/auth/callback"
    prometheus_external_url   = "https://${var.base_address}/prometheus"
    prometheus_callback       = "https://${var.base_address}/prometheus/auth/callback"

    ingress_kind     = "${var.ingress_kind}"
    ingress_tls_cert = "${base64encode(tls_locally_signed_cert.ingress.cert_pem)}"
    ingress_tls_key  = "${base64encode(tls_private_key.ingress.private_key_pem)}"

    identity_server_tls_cert = "${base64encode(tls_locally_signed_cert.identity_server.cert_pem)}"
    identity_server_tls_key  = "${base64encode(tls_private_key.identity_server.private_key_pem)}"
    identity_client_tls_cert = "${base64encode(tls_locally_signed_cert.identity_client.cert_pem)}"
    identity_client_tls_key  = "${base64encode(tls_private_key.identity_client.private_key_pem)}"

    kubectl_client_id = "${var.kubectl_client_id}"
    kubectl_secret    = "${random_id.kubectl_secret.b64}"

    kube_apiserver_url = "${var.kube_apiserver_url}"
    oidc_issuer_url    = "https://${var.base_address}/identity"
    stats_url          = "${var.stats_url}"

    # TODO: We could also patch https://www.terraform.io/docs/providers/random/ to add an UUID resource.
    cluster_id   = "${format("%s-%s-%s-%s-%s", substr(random_id.cluster_id.hex, 0, 8), substr(random_id.cluster_id.hex, 8, 4), substr(random_id.cluster_id.hex, 12, 4), substr(random_id.cluster_id.hex, 16, 4), substr(random_id.cluster_id.hex, 20, 12))}"
    cluster_name = "${var.cluster_name}"

    platform                 = "${var.platform}"
    certificates_strategy    = "${var.ca_generated == "true" ? "installerGeneratedCA" : "userProvidedCA"}"
    identity_api_service     = "${var.identity_api_service}"
    tectonic_updater_enabled = "${var.experimental ? "true" : "false"}"

    image_re = "${var.image_re}"
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
data "template_file" "tectonic_rkt" {
  template = "${file("${path.module}/resources/tectonic-rkt.sh")}"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
    experimental    = "${var.experimental ? "true" : "false"}"
  }
}

resource "local_file" "tectonic_rkt" {
  content  = "${data.template_file.tectonic_rkt.rendered}"
  filename = "./generated/tectonic-rkt.sh"
}

# tectonic.service (available as output variable)
data "template_file" "tectonic_service" {
  template = "${file("${path.module}/resources/tectonic.service")}"
}
