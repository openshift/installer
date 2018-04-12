# Unique Cluster ID (uuid)
resource "random_id" "cluster_id" {
  byte_length = 16
}

# Kubernetes Manifests (resources/generated/manifests/)
resource "template_dir" "tectonic" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "./generated/tectonic"

  vars {
    addon_resizer_image                        = "${var.container_images["addon_resizer"]}"
    kube_core_operator_image                   = "${var.container_images["kube_core_operator"]}"
    kubernetes_addon_operator_image            = "${var.container_images["kubernetes_addon_operator"]}"
    tectonic_channel_operator_image            = "${var.container_images["tectonic_channel_operator"]}"
    tectonic_prometheus_operator_image         = "${var.container_images["tectonic_prometheus_operator"]}"
    tectonic_cluo_operator_image               = "${var.container_images["tectonic_cluo_operator"]}"
    tectonic_alm_operator_image                = "${var.container_images["tectonic_alm_operator"]}"
    tectonic_ingress_controller_operator_image = "${var.container_images["tectonic_ingress_controller_operator"]}"
    tectonic_utility_operator_image            = "${var.container_images["tectonic_utility_operator"]}"

    tectonic_monitoring_auth_base_image = "${var.container_base_images["tectonic_monitoring_auth"]}"
    config_reload_base_image            = "${var.container_base_images["config_reload"]}"
    addon_resizer_base_image            = "${var.container_base_images["addon_resizer"]}"
    kube_state_metrics_base_image       = "${var.container_base_images["kube_state_metrics"]}"
    prometheus_operator_base_image      = "${var.container_base_images["prometheus_operator"]}"
    prometheus_config_reload_base_image = "${var.container_base_images["prometheus_config_reload"]}"
    prometheus_base_image               = "${var.container_base_images["prometheus"]}"
    alertmanager_base_image             = "${var.container_base_images["alertmanager"]}"
    node_exporter_base_image            = "${var.container_base_images["node_exporter"]}"
    grafana_base_image                  = "${var.container_base_images["grafana"]}"
    grafana_watcher_base_image          = "${var.container_base_images["grafana_watcher"]}"
    kube_rbac_proxy_base_image          = "${var.container_base_images["kube_rbac_proxy"]}"

    monitoring_version             = "${var.versions["monitoring"]}"
    tectonic_version               = "${var.versions["tectonic"]}"
    tectonic_cluo_operator_version = "${var.versions["cluo"]}"
    tectonic_alm_operator_version  = "${var.versions["alm"]}"

    license     = "${base64encode(file(var.license_path))}"
    pull_secret = "${base64encode(file(var.pull_secret_path))}"

    update_server  = "${var.update_server}"
    update_channel = "${var.update_channel}"
    update_app_id  = "${var.update_app_id}"

    admin_email = "${lower(var.admin_email)}"

    base_address = "${var.base_address}"

    ingress_ca_cert  = "${base64encode(var.ingress_ca_cert_pem)}"
    ingress_tls_cert = "${base64encode(var.ingress_cert_pem)}"
    ingress_tls_key  = "${base64encode(var.ingress_key_pem)}"

    identity_server_tls_cert = "${base64encode(var.identity_server_cert_pem)}"
    identity_server_tls_key  = "${base64encode(var.identity_server_key_pem)}"
    identity_server_ca_cert  = "${base64encode(var.identity_server_ca_cert)}"
    identity_client_tls_cert = "${base64encode(var.identity_client_cert_pem)}"
    identity_client_tls_key  = "${base64encode(var.identity_client_key_pem)}"
    identity_client_ca_cert  = "${base64encode(var.identity_client_ca_cert)}"

    platform = "${var.platform}"
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

# tectonic.sh (resources/generated/tectonic-wrapper.sh)
data "template_file" "tectonic_wrapper" {
  template = "${file("${path.module}/resources/tectonic-wrapper.sh")}"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
  }
}

resource "local_file" "tectonic_wrapper" {
  content  = "${data.template_file.tectonic_wrapper.rendered}"
  filename = "./generated/tectonic-wrapper.sh"
}

# tectonic.service (available as output variable)
data "template_file" "tectonic_service" {
  template = "${file("${path.module}/resources/tectonic.service")}"
}

data "ignition_systemd_unit" "tectonic_service" {
  name    = "tectonic.service"
  enabled = false
  content = "${data.template_file.tectonic_service.rendered}"
}

# tectonic.path (available as output variable)
data "template_file" "tectonic_path" {
  template = "${file("${path.module}/resources/tectonic.path")}"
}

data "ignition_systemd_unit" "tectonic_path" {
  name    = "tectonic.path"
  enabled = true
  content = "${data.template_file.tectonic_path.rendered}"
}
