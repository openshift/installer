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
    "secrets/ingress-tls.yaml",
    "secrets/pull.json",
    "security/priviledged-scc-tectonic.yaml",
    "updater/app-version-kind.yaml",
    "updater/app_versions/app-version-kube-core.yaml",
    "updater/app_versions/app-version-kube-addon.yaml",
    "updater/app_versions/app-version-tectonic-alm.yaml",
    "updater/app_versions/app-version-tectonic-cluster.yaml",
    "updater/app_versions/app-version-tectonic-ingress.yaml",
    "updater/app_versions/app-version-tectonic-utility.yaml",
    "updater/migration-status-kind.yaml",
    "updater/operators/kube-core-operator.yaml",
    "updater/operators/kube-addon-operator.yaml",
    "updater/operators/tectonic-alm-operator.yaml",
    "updater/operators/tectonic-channel-operator.yaml",
    "updater/operators/tectonic-ingress-controller-operator.yaml",
    "updater/operators/tectonic-utility-operator.yaml",
    "updater/operators/cluster-openshift-apiserver-operator.yaml",
    "updater/tectonic-channel-operator-config.yaml",
    "updater/tectonic-channel-operator-kind.yaml",
  ]
}

# Kubernetes Manifests, rendered
data "template_file" "manifest_file_list" {
  count    = "${length(var.manifest_names)}"
  template = "${file("${path.module}/resources/manifests/${var.manifest_names[count.index]}")}"

  vars {
    addon_resizer_image                        = "${var.container_images["addon_resizer"]}"
    kube_core_operator_image                   = "${var.container_images["kube_core_operator"]}"
    kube_addon_operator_image                  = "${var.container_images["kube_addon_operator"]}"
    tectonic_channel_operator_image            = "${var.container_images["tectonic_channel_operator"]}"
    tectonic_alm_operator_image                = "${var.container_images["tectonic_alm_operator"]}"
    tectonic_ingress_controller_operator_image = "${var.container_images["tectonic_ingress_controller_operator"]}"
    tectonic_utility_operator_image            = "${var.container_images["tectonic_utility_operator"]}"
    cluster_openshift_apiserver_operator_image = "${var.container_images["cluster_openshift_apiserver_operator"]}"

    config_reload_base_image      = "${var.container_base_images["config_reload"]}"
    addon_resizer_base_image      = "${var.container_base_images["addon_resizer"]}"
    kube_state_metrics_base_image = "${var.container_base_images["kube_state_metrics"]}"
    kube_rbac_proxy_base_image    = "${var.container_base_images["kube_rbac_proxy"]}"

    tectonic_version              = "${var.versions["tectonic"]}"
    tectonic_alm_operator_version = "${var.versions["alm"]}"

    pull_secret = "${base64encode(var.pull_secret)}"

    update_server  = "${var.update_server}"
    update_channel = "${var.update_channel}"
    update_app_id  = "${var.update_app_id}"

    base_address = "${var.base_address}"

    ingress_kind            = "${var.ingress_kind}"
    ingress_status_password = "${random_string.ingress_status_password.result}"
    ingress_ca_cert         = "${base64encode(var.ingress_ca_cert_pem)}"
    ingress_tls_cert        = "${base64encode(var.ingress_cert_pem)}"
    ingress_tls_key         = "${base64encode(var.ingress_key_pem)}"
    ingress_tls_bundle      = "${base64encode(var.ingress_bundle_pem)}"

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
