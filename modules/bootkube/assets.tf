data "null_data_source" "etcd" {
  inputs = {
    ca_flag   = "${var.etcd_ca_cert != "" ? "- --etcd-cafile=/etc/kubernetes/secrets/etcd-ca.crt" : "# no etcd-ca.crt given" }"
    cert_flag = "${var.etcd_client_cert != "" ? "- --etcd-certfile=/etc/kubernetes/secrets/etcd-client.crt" : "# no etcd-client.crt given" }"
    key_flag  = "${var.etcd_client_key != "" ? "- --etcd-keyfile=/etc/kubernetes/secrets/etcd-client.key" : "# no etcd-client.key given" }"

    # The file() interpolation function expects an existing file to be present, even if used inside a ternary operator branch.
    ca_path   = "${var.etcd_ca_cert != "" ? var.etcd_ca_cert : "/dev/null" }"
    cert_path = "${var.etcd_client_cert != "" ? var.etcd_client_cert : "/dev/null" }"
    key_path  = "${var.etcd_client_key != "" ? var.etcd_client_key : "/dev/null" }"

    no_certs = "${var.etcd_ca_cert == "" && var.etcd_client_cert == "" && var.etcd_client_key == "" ? 1 : 0}"
  }
}

# Self-hosted manifests (resources/generated/manifests/)
resource "template_dir" "bootkube" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "${path.cwd}/generated/manifests"

  vars {
    hyperkube_image        = "${var.container_images["hyperkube"]}"
    pod_checkpointer_image = "${var.container_images["pod_checkpointer"]}"
    kubedns_image          = "${var.container_images["kubedns"]}"
    kubednsmasq_image      = "${var.container_images["kubednsmasq"]}"
    kubedns_sidecar_image  = "${var.container_images["kubedns_sidecar"]}"
    flannel_image          = "${var.container_images["flannel"]}"
    etcd_operator_image    = "${var.container_images["etcd_operator"]}"
    kenc_image             = "${var.container_images["kenc"]}"

    # Choose the etcd endpoints to use.
    # 1. If experimental mode is enabled (self-hosted etcd), then use
    # var.etcd_service_ip.
    # 2. Else if no etcd TLS certificates are provided, i.e. we bootstrap etcd
    # nodes ourselves (using http), then use insecure http var.etcd_endpoints.
    # 3. Else (if etcd TLS certific are provided), then use the secure https
    # var.etcd_endpoints.
    etcd_servers = "${
      var.experimental_enabled 
        ? format("http://%s:2379", cidrhost(var.service_cidr, 15))
        : data.null_data_source.etcd.outputs.no_certs
          ? join(",", formatlist("http://%s:2379", var.etcd_endpoints))
          : join(",", formatlist("https://%s:2379", var.etcd_endpoints))
      }"

    etcd_ca_flag   = "${data.null_data_source.etcd.outputs.ca_flag}"
    etcd_cert_flag = "${data.null_data_source.etcd.outputs.cert_flag}"
    etcd_key_flag  = "${data.null_data_source.etcd.outputs.key_flag}"

    etcd_service_ip = "${cidrhost(var.service_cidr, 15)}"

    cloud_provider = "${var.cloud_provider}"

    cluster_cidr        = "${var.cluster_cidr}"
    service_cidr        = "${var.service_cidr}"
    kube_dns_service_ip = "${cidrhost(var.service_cidr, 10)}"
    advertise_address   = "${var.advertise_address}"

    anonymous_auth      = "${var.anonymous_auth}"
    oidc_issuer_url     = "${var.oidc_issuer_url}"
    oidc_client_id      = "${var.oidc_client_id}"
    oidc_username_claim = "${var.oidc_username_claim}"
    oidc_groups_claim   = "${var.oidc_groups_claim}"

    ca_cert            = "${base64encode(var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube-ca.*.cert_pem) : var.ca_cert)}"
    apiserver_key      = "${base64encode(tls_private_key.apiserver.private_key_pem)}"
    apiserver_cert     = "${base64encode(tls_locally_signed_cert.apiserver.cert_pem)}"
    serviceaccount_pub = "${base64encode(tls_private_key.service-account.public_key_pem)}"
    serviceaccount_key = "${base64encode(tls_private_key.service-account.private_key_pem)}"

    etcd_ca_cert     = "${base64encode(file(data.null_data_source.etcd.outputs.ca_path))}"
    etcd_client_cert = "${base64encode(file(data.null_data_source.etcd.outputs.cert_path))}"
    etcd_client_key  = "${base64encode(file(data.null_data_source.etcd.outputs.key_path))}"
  }
}

# Self-hosted bootstrapping manifests (resources/generated/manifests-bootstrap/)
resource "template_dir" "bootkube-bootstrap" {
  source_dir      = "${path.module}/resources/bootstrap-manifests"
  destination_dir = "${path.cwd}/generated/bootstrap-manifests"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
    etcd_image      = "${var.container_images["etcd"]}"

    etcd_servers = "${
      var.experimental_enabled 
        ? format("http://%s:2379,http://127.0.0.1:12379", cidrhost(var.service_cidr, 15))
        : data.null_data_source.etcd.outputs.no_certs
          ? join(",", formatlist("http://%s:2379", var.etcd_endpoints))
          : join(",", formatlist("https://%s:2379", var.etcd_endpoints))
      }"

    etcd_ca_flag   = "${data.null_data_source.etcd.outputs.ca_flag}"
    etcd_cert_flag = "${data.null_data_source.etcd.outputs.cert_flag}"
    etcd_key_flag  = "${data.null_data_source.etcd.outputs.key_flag}"

    advertise_address = "${var.advertise_address}"
    cloud_provider    = "${var.cloud_provider}"
    cluster_cidr      = "${var.cluster_cidr}"
    service_cidr      = "${var.service_cidr}"
  }
}

# Self-hosted experimental etcd
data "template_file" "etcd-operator" {
  template = "${file("${path.module}/resources/experimental/manifests/etcd-operator.yaml")}"

  vars {
    etcd_operator_image = "${var.container_images["etcd_operator"]}"
  }
}

resource "local_file" "etcd-operator" {
  count      = "${var.experimental_enabled ? 1 : 0}"
  depends_on = ["template_dir.bootkube"]

  content  = "${data.template_file.etcd-operator.rendered}"
  filename = "${path.cwd}/generated/experimental/etcd-operator.yaml"
}

data "template_file" "etcd-service" {
  template = "${file("${path.module}/resources/experimental/manifests/etcd-service.yaml")}"

  vars {
    etcd_service_ip = "${cidrhost(var.service_cidr, 15)}"
  }
}

resource "local_file" "etcd-service" {
  count      = "${var.experimental_enabled ? 1 : 0}"
  depends_on = ["template_dir.bootkube"]

  content  = "${data.template_file.etcd-service.rendered}"
  filename = "${path.cwd}/generated/experimental/etcd-service.yaml"
}

data "template_file" "bootstrap-etcd" {
  template = "${file("${path.module}/resources/experimental/bootstrap-manifests/bootstrap-etcd.yaml")}"

  vars {
    etcd_image = "${var.container_images["etcd"]}"
  }
}

resource "local_file" "bootstrap-etcd" {
  count      = "${var.experimental_enabled ? 1 : 0}"
  depends_on = ["template_dir.bootkube-bootstrap"]

  content  = "${data.template_file.bootstrap-etcd.rendered}"
  filename = "${path.cwd}/generated/bootstrap-experimental/bootstrap-etcd.yaml"
}

# etcd certs
resource "local_file" "etcd_ca_crt" {
  count    = "${var.etcd_ca_cert == "" ? 0 : 1}"
  content  = "${file(var.etcd_ca_cert)}"
  filename = "${path.cwd}/generated/tls/etcd-ca.crt"
}

resource "local_file" "etcd_client_crt" {
  count    = "${var.etcd_client_cert == "" ? 0 : 1}"
  content  = "${file(var.etcd_client_cert)}"
  filename = "${path.cwd}/generated/tls/etcd-client.crt"
}

resource "local_file" "etcd_client_key" {
  count    = "${var.etcd_client_key == "" ? 0 : 1}"
  content  = "${file(var.etcd_client_key)}"
  filename = "${path.cwd}/generated/tls/etcd-client.key"
}

# kubeconfig (resources/generated/auth/kubeconfig)
data "template_file" "kubeconfig" {
  template = "${file("${path.module}/resources/kubeconfig")}"

  vars {
    ca_cert      = "${base64encode(var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube-ca.*.cert_pem) : var.ca_cert)}"
    kubelet_cert = "${base64encode(tls_locally_signed_cert.kubelet.cert_pem)}"
    kubelet_key  = "${base64encode(tls_private_key.kubelet.private_key_pem)}"
    server       = "${var.kube_apiserver_url}"
  }
}

resource "local_file" "kubeconfig" {
  content  = "${data.template_file.kubeconfig.rendered}"
  filename = "${path.cwd}/generated/auth/kubeconfig"
}

# bootkube.sh (resources/generated/bootkube.sh)
data "template_file" "bootkube" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image = "${var.container_images["bootkube"]}"
  }
}

resource "local_file" "bootkube" {
  content  = "${data.template_file.bootkube.rendered}"
  filename = "${path.cwd}/generated/bootkube.sh"
}

# bootkube.service (available as output variable)
data "template_file" "bootkube_service" {
  template = "${file("${path.module}/resources/bootkube.service")}"
}
