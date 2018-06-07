variable "manifest_names" {
  default = [
    "01-tectonic-namespace.yaml",
    "02-ingress-namespace.yaml",
    "03-openshift-web-console-namespace.yaml",
    "app-version-kind.yaml",
    "app-version-tectonic-network.yaml",
    "app-version-tnc.yaml",
    "kube-apiserver-secret.yaml",
    "kube-cloud-config.yaml",
    "kube-controller-manager-secret.yaml",
    "node-config-kind.yaml",
    "openshift-apiserver-secret.yaml",
    "pull.json",
    "tectonic-network-operator.yaml",
    "tectonic-node-controller-operator.yaml",
  ]
}

# Self-hosted manifests (resources/generated/manifests/)
data "template_file" "manifest_file_list" {
  count    = "${length(var.manifest_names)}"
  template = "${file("${path.module}/resources/manifests/${var.manifest_names[count.index]}")}"

  vars {
    tectonic_network_operator_image = "${var.container_images["tectonic_network_operator"]}"
    tnc_operator_image              = "${var.container_images["tnc_operator"]}"

    cloud_provider_config = "${var.cloud_provider_config}"

    root_ca_cert             = "${base64encode(var.root_ca_cert_pem)}"
    aggregator_ca_cert       = "${base64encode(var.aggregator_ca_cert_pem)}"
    aggregator_ca_key        = "${base64encode(var.aggregator_ca_key_pem)}"
    kube_ca_cert             = "${base64encode(var.kube_ca_cert_pem)}"
    kube_ca_key              = "${base64encode(var.kube_ca_key_pem)}"
    service_serving_ca_cert  = "${base64encode(var.service_serving_ca_cert_pem)}"
    service_serving_ca_key   = "${base64encode(var.service_serving_ca_key_pem)}"
    apiserver_key            = "${base64encode(var.apiserver_key_pem)}"
    apiserver_cert           = "${base64encode(var.apiserver_cert_pem)}"
    openshift_apiserver_key  = "${base64encode(var.openshift_apiserver_key_pem)}"
    openshift_apiserver_cert = "${base64encode(var.openshift_apiserver_cert_pem)}"
    apiserver_proxy_key      = "${base64encode(var.apiserver_proxy_key_pem)}"
    apiserver_proxy_cert     = "${base64encode(var.apiserver_proxy_cert_pem)}"
    oidc_ca_cert             = "${base64encode(var.oidc_ca_cert)}"
    pull_secret              = "${base64encode(file(var.pull_secret_path))}"
    serviceaccount_pub       = "${base64encode(tls_private_key.service_account.public_key_pem)}"
    serviceaccount_key       = "${base64encode(tls_private_key.service_account.private_key_pem)}"
    kube_dns_service_ip      = "${cidrhost(var.service_cidr, 10)}"

    openshift_loopback_kubeconfig = "${base64encode(data.template_file.kubeconfig.rendered)}"

    etcd_ca_cert     = "${base64encode(var.etcd_ca_cert_pem)}"
    etcd_client_cert = "${base64encode(var.etcd_client_cert_pem)}"
    etcd_client_key  = "${base64encode(var.etcd_client_key_pem)}"
  }
}

# Ignition entry for every bootkube manifest
# Drops them in /opt/tectonic/manifests/<path>
data "ignition_file" "manifest_file_list" {
  count      = "${length(var.manifest_names)}"
  filesystem = "root"
  mode       = "0644"

  path = "/opt/tectonic/manifests/${var.manifest_names[count.index]}"

  content {
    content = "${data.template_file.manifest_file_list.*.rendered[count.index]}"
  }
}

# Log the generated manifest files to disk for debugging and user visibility
# Dest: ./generated/manifests/<path>
resource "local_file" "manifest_files" {
  count    = "${length(var.manifest_names)}"
  filename = "./generated/manifests/${var.manifest_names[count.index]}"
  content  = "${data.template_file.manifest_file_list.*.rendered[count.index]}"
}
