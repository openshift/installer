# Kubernetes Deployments and AppVersions TPR used in the payload
resource "template_dir" "payload_operators" {
  source_dir      = "../tectonic/resources/manifests/updater/operators"
  destination_dir = "./generated/operators"

  vars {
    kube_version_operator_image        = "${var.tectonic_container_images["kube_version_operator"]}"
    tectonic_channel_operator_image    = "${var.tectonic_container_images["tectonic_channel_operator"]}"
    tectonic_prometheus_operator_image = "${var.tectonic_container_images["tectonic_prometheus_operator"]}"
    tectonic_etcd_operator_image       = "${var.tectonic_container_images["tectonic_etcd_operator"]}"
    tectonic_cluo_operator_image       = "${var.tectonic_container_images["tectonic_cluo_operator"]}"
    kubernetes_addon_operator_image    = "${var.tectonic_container_images["kubernetes_addon_operator"]}"
    tectonic_alm_operator_image        = "${var.tectonic_container_images["tectonic_alm_operator"]}"
    tectonic_utility_operator_image    = "${var.tectonic_container_images["tectonic_utility_operator"]}"
  }
}

# TNO
data "template_file" "tectonic_network_operator" {
  template = "${file("${path.module}/../bootkube/resources/manifests/tectonic-network-operator.yaml")}"

  vars {
    tectonic_network_operator_image = "${var.tectonic_container_images["tectonic_network_operator"]}"
  }
}

resource "local_file" "tectonic_network_operator" {
  content  = "${data.template_file.tectonic_network_operator.rendered}"
  filename = "./generated/operators/tectonic-network-operator.yaml"

  depends_on = ["template_dir.payload_operators"]
}

resource "template_dir" "payload_appversions" {
  source_dir      = "../tectonic/resources/manifests/updater/app_versions"
  destination_dir = "./generated/app_versions"

  vars {
    kubernetes_version             = "${var.tectonic_versions["kubernetes"]}"
    monitoring_version             = "${var.tectonic_versions["monitoring"]}"
    tectonic_version               = "${var.tectonic_versions["tectonic"]}"
    tectonic_etcd_operator_version = "${var.tectonic_versions["tectonic-etcd"]}"
    tectonic_cluo_operator_version = "${var.tectonic_versions["cluo"]}"
    tectonic_alm_operator_version  = "${var.tectonic_versions["alm"]}"
  }
}

# TNO
resource "local_file" "appversion_tectonic_network" {
  content  = "${file("${path.module}/../bootkube/resources/manifests/app-version-tectonic-network.yaml")}"
  filename = "./generated/app_versions/app-version-tectonic-network.yaml"

  depends_on = ["template_dir.payload_appversions"]
}
