#List of all tectonic manifests
# To update this list, issue this command from the repo root:
# find modules/tectonic/resources/manifests -type f -printf '    "%P",\n' | sort
variable "manifest_names" {
  default = [
    "ingress/cluster-config.yaml",
    "ingress/pull.json",
    "ingress/README.md",
    "ingress/svc-account.yaml",
    "rbac/binding-admin.yaml",
    "rbac/binding-discovery.yaml",
    "rbac/role-admin.yaml",
    "rbac/role-user.yaml",
    "secrets/ca-cert.yaml",
    "secrets/identity-grpc-client.yaml",
    "secrets/identity-grpc-server.yaml",
    "secrets/ingress-tls.yaml",
    "secrets/license.json",
    "secrets/pull.json",
    "updater/app-version-kind.yaml",
    "updater/app_versions/app-version-kube-core.yaml",
    "updater/app_versions/app-version-kubernetes-addon.yaml",
    "updater/app_versions/app-version-tectonic-alm.yaml",
    "updater/app_versions/app-version-tectonic-cluster.yaml",
    "updater/app_versions/app-version-tectonic-ingress.yaml",
    "updater/app_versions/app-version-tectonic-monitoring.yaml",
    "updater/app_versions/app-version-tectonic-utility.yaml",
    "updater/migration-status-kind.yaml",
    "updater/operators/kube-core-operator.yaml",
    "updater/operators/kubernetes-addon-operator.yaml",
    "updater/operators/tectonic-alm-operator.yaml",
    "updater/operators/tectonic-channel-operator.yaml",
    "updater/operators/tectonic-ingress-controller-operator.yaml",
    "updater/operators/tectonic-prometheus-operator.yaml",
    "updater/operators/tectonic-utility-operator.yaml",
    "updater/tectonic-channel-operator-config.yaml",
    "updater/tectonic-channel-operator-kind.yaml",
    "updater/tectonic-monitoring-config.yaml",
  ]
}

# Kubernetes Manifests, rendered
data "template_file" "manifest_file_list" {
  count    = "${length(var.manifest_names)}"
  template = "${file("${path.module}/resources/manifests/${var.manifest_names[count.index]}")}"

  vars {
    addon_resizer_image                        = "${var.container_images["addon_resizer"]}"
    kube_core_operator_image                   = "${var.container_images["kube_core_operator"]}"
    kubernetes_addon_operator_image            = "${var.container_images["kubernetes_addon_operator"]}"
    tectonic_channel_operator_image            = "${var.container_images["tectonic_channel_operator"]}"
    tectonic_prometheus_operator_image         = "${var.container_images["tectonic_prometheus_operator"]}"
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

    monitoring_version            = "${var.versions["monitoring"]}"
    tectonic_version              = "${var.versions["tectonic"]}"
    tectonic_alm_operator_version = "${var.versions["alm"]}"

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

# Ignition entry for every tectonic manifest
# Drops them in /opt/tectonic/tectonic/<path>
data "ignition_file" "tectonic_manifest_list" {
  count      = "${length(var.manifest_names)}"
  filesystem = "root"
  mode       = "0644"

  path = "/opt/tectonic/tectonic/${var.manifest_names[count.index]}"

  content {
    content = "${data.template_file.manifest_file_list.*.rendered[count.index]}"
  }
}

# Log the generated manifest files to disk for debugging and user visibility
# Dest: ./generated/tectonic/<path>
resource "local_file" "manifest_files" {
  count    = "${length(var.manifest_names)}"
  filename = "./generated/tectonic/${var.manifest_names[count.index]}"
  content  = "${data.template_file.manifest_file_list.*.rendered[count.index]}"
}
