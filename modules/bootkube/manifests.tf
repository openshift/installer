variable "manifest_names" {
  default = [
    "01-tectonic-namespace.yaml",
    "02-ingress-namespace.yaml",
    "03-openshift-web-console-namespace.yaml",
    "04-openshift-machine-config-operator.yaml",        # https://github.com/openshift/machine-config-operator/tree/master/install/00_namespace.yaml
    "05-openshift-cluster-api-namespace.yaml",
    "app-version-kind.yaml",
    "app-version-tectonic-network.yaml",
    "kube-apiserver-secret.yaml",
    "kube-cloud-config.yaml",
    "kube-controller-manager-secret.yaml",
    "machine-config-operator-00-config-crd.yaml",       # https://github.com/openshift/machine-config-operator/tree/master/install/01_mcoconfig.crd.yaml
    "machine-config-operator-01-images-configmap.yaml", # https://github.com/openshift/machine-config-operator/tree/master/install/02_images.configmap.yaml
    "machine-config-operator-02-rbac.yaml",             # https://github.com/openshift/machine-config-operator/tree/master/install/03_rbac.yaml
    "machine-config-operator-03-deployment.yaml",       # https://github.com/openshift/machine-config-operator/tree/master/install/04_deployment.yaml
    "machine-config-server-tls-secret.yaml",
    "openshift-apiserver-secret.yaml",
    "cluster-apiserver-certs.yaml",
    "pull.json",
    "tectonic-network-operator.yaml",
    "operatorstatus-crd.yaml",
    "app-version-mao.yaml",
    "machine-api-operator.yaml",
    "ign-config.yaml",
  ]
}

# Self-hosted manifests (resources/generated/manifests/)
data "template_file" "manifest_file_list" {
  count    = "${length(var.manifest_names)}"
  template = "${file("${path.module}/resources/manifests/${var.manifest_names[count.index]}")}"

  vars {
    tectonic_network_operator_image = "${var.container_images["tectonic_network_operator"]}"
    machine_config_operator_image   = "${var.container_images["machine_config_operator"]}"

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
    clusterapi_ca_cert       = "${base64encode(var.clusterapi_ca_cert_pem)}"
    clusterapi_ca_key        = "${base64encode(var.clusterapi_ca_key_pem)}"
    oidc_ca_cert             = "${base64encode(var.oidc_ca_cert)}"
    pull_secret              = "${base64encode(var.pull_secret)}"
    serviceaccount_pub       = "${base64encode(var.service_account_public_key_pem)}"
    serviceaccount_key       = "${base64encode(var.service_account_private_key_pem)}"
    kube_dns_service_ip      = "${cidrhost(var.service_cidr, 10)}"

    openshift_loopback_kubeconfig = "${base64encode(data.template_file.kubeconfig.rendered)}"

    etcd_ca_cert     = "${base64encode(var.etcd_ca_cert_pem)}"
    etcd_client_cert = "${base64encode(var.etcd_client_cert_pem)}"
    etcd_client_key  = "${base64encode(var.etcd_client_key_pem)}"

    mcs_tls_cert = "${base64encode(var.mcs_cert_pem)}"
    mcs_tls_key  = "${base64encode(var.mcs_key_pem)}"

    worker_ign_config = "${base64encode(var.worker_ign_config)}"
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

# Conditionally include the libvirt-certs secret.
data "template_file" "libvirt_certs" {
  template = "${file("${path.module}/resources/manifests/libvirt-certs-secret.yaml")}"

  vars {
    ca_cert     = "${base64encode(var.libvirt_tls_ca_pem)}"
    client_cert = "${base64encode(var.libvirt_tls_cert_pem)}"
    client_key  = "${base64encode(var.libvirt_tls_key_pem)}"
  }
}

data "ignition_file" "libvirt_certs_secret" {
  count      = "${var.libvirt_tls_cert_pem != "" ? 1 : 0}"
  filesystem = "root"
  mode       = "0644"

  path = "/opt/tectonic/manifests/libvirt-certs-secret.yaml"

  content {
    content = "${data.template_file.libvirt_certs.rendered}"
  }
}

resource "local_file" "libvirt_certs" {
  count    = "${var.libvirt_tls_cert_pem != "" ? 1 : 0}"
  filename = "./generated/manifests/libvirt-certs-secret.yaml"
  content  = "${data.template_file.libvirt_certs.rendered}"
}
