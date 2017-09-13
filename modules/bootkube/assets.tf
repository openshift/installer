resource "template_dir" "experimental" {
  count           = "${var.experimental_enabled ? 1 : 0}"
  source_dir      = "${path.module}/resources/experimental/manifests"
  destination_dir = "./generated/experimental"

  vars {
    etcd_operator_image = "${var.container_images["etcd_operator"]}"
    etcd_service_ip     = "${cidrhost(var.service_cidr, 15)}"
    kenc_image          = "${var.container_images["kenc"]}"

    etcd_ca_cert = "${base64encode(data.template_file.etcd_ca_cert_pem.rendered)}"

    etcd_server_cert = "${base64encode(join("", tls_locally_signed_cert.etcd_server.*.cert_pem))}"
    etcd_server_key  = "${base64encode(join("", tls_private_key.etcd_server.*.private_key_pem))}"

    etcd_client_cert = "${base64encode(data.template_file.etcd_client_crt.rendered)}"
    etcd_client_key  = "${base64encode(data.template_file.etcd_client_key.rendered)}"

    etcd_peer_cert = "${base64encode(join("", tls_locally_signed_cert.etcd_peer.*.cert_pem))}"
    etcd_peer_key  = "${base64encode(join("", tls_private_key.etcd_peer.*.private_key_pem))}"
  }
}

resource "template_dir" "bootstrap_experimental" {
  count           = "${var.experimental_enabled ? 1 : 0}"
  source_dir      = "${path.module}/resources/experimental/bootstrap-manifests"
  destination_dir = "./generated/bootstrap-experimental"

  vars {
    etcd_image                = "${var.container_images["etcd"]}"
    etcd_version              = "${var.versions["etcd"]}"
    bootstrap_etcd_service_ip = "${cidrhost(var.service_cidr, 20)}"
  }
}

resource "template_dir" "etcd_experimental" {
  count           = "${var.experimental_enabled ? 1 : 0}"
  source_dir      = "${path.module}/resources/experimental/etcd"
  destination_dir = "./generated/etcd"

  vars {
    etcd_version              = "${var.versions["etcd"]}"
    bootstrap_etcd_service_ip = "${cidrhost(var.service_cidr, 20)}"
  }
}

# Self-hosted manifests (resources/generated/manifests/)
resource "template_dir" "bootkube" {
  source_dir      = "${path.module}/resources/manifests"
  destination_dir = "./generated/manifests"

  vars {
    hyperkube_image        = "${var.container_images["hyperkube"]}"
    pod_checkpointer_image = "${var.container_images["pod_checkpointer"]}"
    kubedns_image          = "${var.container_images["kubedns"]}"
    kubednsmasq_image      = "${var.container_images["kubednsmasq"]}"
    kubedns_sidecar_image  = "${var.container_images["kubedns_sidecar"]}"

    # Choose the etcd endpoints to use.
    # 1. If experimental mode is enabled (self-hosted etcd), then use
    # var.etcd_service_ip.
    # 2. Else if no etcd TLS certificates are provided, i.e. we bootstrap etcd
    # nodes ourselves (using http), then use insecure http var.etcd_endpoints.
    # 3. Else (if etcd TLS certific are provided), then use the secure https
    # var.etcd_endpoints.
    etcd_servers = "${
      var.experimental_enabled
        ? format("https://%s:2379", cidrhost(var.service_cidr, 15))
        : data.template_file.etcd_ca_cert_pem.rendered == ""
          ? join(",", formatlist("http://%s:2379", var.etcd_endpoints))
          : join(",", formatlist("https://%s:2379", var.etcd_endpoints))
      }"

    etcd_service_ip           = "${cidrhost(var.service_cidr, 15)}"
    bootstrap_etcd_service_ip = "${cidrhost(var.service_cidr, 20)}"

    cloud_provider             = "${var.cloud_provider}"
    cloud_provider_config      = "${var.cloud_provider_config}"
    cloud_provider_config_flag = "${var.cloud_provider_config != "" ? "- --cloud-config=/etc/kubernetes/cloud/config" : "# no cloud provider config given"}"

    cluster_cidr        = "${var.cluster_cidr}"
    service_cidr        = "${var.service_cidr}"
    kube_dns_service_ip = "${cidrhost(var.service_cidr, 10)}"
    advertise_address   = "${var.advertise_address}"

    anonymous_auth      = "${var.anonymous_auth}"
    oidc_issuer_url     = "${var.oidc_issuer_url}"
    oidc_client_id      = "${var.oidc_client_id}"
    oidc_username_claim = "${var.oidc_username_claim}"
    oidc_groups_claim   = "${var.oidc_groups_claim}"

    ca_cert            = "${base64encode(var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem) : var.ca_cert)}"
    apiserver_key      = "${base64encode(tls_private_key.apiserver.private_key_pem)}"
    apiserver_cert     = "${base64encode(tls_locally_signed_cert.apiserver.cert_pem)}"
    serviceaccount_pub = "${base64encode(tls_private_key.service_account.public_key_pem)}"
    serviceaccount_key = "${base64encode(tls_private_key.service_account.private_key_pem)}"

    etcd_ca_flag   = "${data.template_file.etcd_ca_cert_pem.rendered != "" ? "- --etcd-cafile=/etc/kubernetes/secrets/etcd-client-ca.crt" : "# no etcd-client-ca.crt given" }"
    etcd_cert_flag = "${data.template_file.etcd_client_crt.rendered != "" ? "- --etcd-certfile=/etc/kubernetes/secrets/etcd-client.crt" : "# no etcd-client.crt given" }"
    etcd_key_flag  = "${data.template_file.etcd_client_key.rendered != "" ? "- --etcd-keyfile=/etc/kubernetes/secrets/etcd-client.key" : "# no etcd-client.key given" }"

    etcd_ca_cert     = "${base64encode(data.template_file.etcd_ca_cert_pem.rendered)}"
    etcd_client_cert = "${base64encode(data.template_file.etcd_client_crt.rendered)}"
    etcd_client_key  = "${base64encode(data.template_file.etcd_client_key.rendered)}"

    kubernetes_version = "${replace(var.versions["kubernetes"], "+", "-")}"

    master_count              = "${var.master_count}"
    node_monitor_grace_period = "${var.node_monitor_grace_period}"
    pod_eviction_timeout      = "${var.pod_eviction_timeout}"
  }
}

# Self-hosted bootstrapping manifests (resources/generated/manifests-bootstrap/)
resource "template_dir" "bootkube_bootstrap" {
  source_dir      = "${path.module}/resources/bootstrap-manifests"
  destination_dir = "./generated/bootstrap-manifests"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
    etcd_image      = "${var.container_images["etcd"]}"

    etcd_servers = "${
      var.experimental_enabled
        ? format("https://%s:2379,https://127.0.0.1:12379", cidrhost(var.service_cidr, 15))
        : data.template_file.etcd_ca_cert_pem.rendered == ""
          ? join(",", formatlist("http://%s:2379", var.etcd_endpoints))
          : join(",", formatlist("https://%s:2379", var.etcd_endpoints))
      }"

    etcd_ca_flag   = "${data.template_file.etcd_ca_cert_pem.rendered != "" ? "- --etcd-cafile=/etc/kubernetes/secrets/etcd-client-ca.crt" : "# no etcd-client-ca.crt given" }"
    etcd_cert_flag = "${data.template_file.etcd_client_crt.rendered != "" ? "- --etcd-certfile=/etc/kubernetes/secrets/etcd-client.crt" : "# no etcd-client.crt given" }"
    etcd_key_flag  = "${data.template_file.etcd_client_key.rendered != "" ? "- --etcd-keyfile=/etc/kubernetes/secrets/etcd-client.key" : "# no etcd-client.key given" }"

    cloud_provider             = "${var.cloud_provider}"
    cloud_provider_config      = "${var.cloud_provider_config}"
    cloud_provider_config_flag = "${var.cloud_provider_config != "" ? "- --cloud-config=/etc/kubernetes/cloud/config" : "# no cloud provider config given"}"

    advertise_address = "${var.advertise_address}"
    cluster_cidr      = "${var.cluster_cidr}"
    service_cidr      = "${var.service_cidr}"
  }
}

# kubeconfig (resources/generated/auth/kubeconfig)
data "template_file" "kubeconfig" {
  template = "${file("${path.module}/resources/kubeconfig")}"

  vars {
    ca_cert      = "${base64encode(var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem) : var.ca_cert)}"
    kubelet_cert = "${base64encode(tls_locally_signed_cert.kubelet.cert_pem)}"
    kubelet_key  = "${base64encode(tls_private_key.kubelet.private_key_pem)}"
    server       = "${var.kube_apiserver_url}"
    cluster_name = "${var.cluster_name}"
  }
}

resource "local_file" "kubeconfig" {
  content  = "${data.template_file.kubeconfig.rendered}"
  filename = "./generated/auth/kubeconfig"
}

# bootkube.sh (resources/generated/bootkube.sh)
data "template_file" "bootkube_sh" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image = "${var.container_images["bootkube"]}"
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

# etcd assets
data "template_file" "etcd_ca_cert_pem" {
  template = "${var.experimental_enabled || var.etcd_tls_enabled
    ? join("", tls_self_signed_cert.etcd_ca.*.cert_pem)
    : file(var.etcd_ca_cert)
  }"
}

data "template_file" "etcd_client_crt" {
  template = "${var.experimental_enabled || var.etcd_tls_enabled
    ? join("", tls_locally_signed_cert.etcd_client.*.cert_pem)
    : file(var.etcd_client_cert)
  }"
}

data "template_file" "etcd_client_key" {
  template = "${var.experimental_enabled || var.etcd_tls_enabled
    ? join("", tls_private_key.etcd_client.*.private_key_pem)
    : file(var.etcd_client_key)
  }"
}

resource "local_file" "etcd_ca_crt" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled || var.etcd_ca_cert != "/dev/null" ? 1 : 0}"
  content  = "${data.template_file.etcd_ca_cert_pem.rendered}"
  filename = "./generated/tls/etcd-client-ca.crt"
}

resource "local_file" "etcd_client_crt" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled || var.etcd_client_cert != "/dev/null" ? 1 : 0}"
  content  = "${data.template_file.etcd_client_crt.rendered}"
  filename = "./generated/tls/etcd-client.crt"
}

resource "local_file" "etcd_client_key" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled || var.etcd_client_key != "/dev/null" ? 1 : 0}"
  content  = "${data.template_file.etcd_client_key.rendered}"
  filename = "./generated/tls/etcd-client.key"
}

resource "local_file" "etcd_server_crt" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled ? 1 : 0}"
  content  = "${join("", tls_locally_signed_cert.etcd_server.*.cert_pem)}"
  filename = "./generated/tls/etcd/server.crt"
}

resource "local_file" "etcd_server_key" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled ? 1 : 0}"
  content  = "${join("", tls_private_key.etcd_server.*.private_key_pem)}"
  filename = "./generated/tls/etcd/server.key"
}

resource "local_file" "etcd_peer_crt" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled ? 1 : 0}"
  content  = "${join("", tls_locally_signed_cert.etcd_peer.*.cert_pem)}"
  filename = "./generated/tls/etcd/peer.crt"
}

resource "local_file" "etcd_peer_key" {
  count    = "${var.experimental_enabled || var.etcd_tls_enabled ? 1 : 0}"
  content  = "${join("", tls_private_key.etcd_peer.*.private_key_pem)}"
  filename = "./generated/tls/etcd/peer.key"
}

data "archive_file" "etcd_tls_zip" {
  type = "zip"

  output_path = "./.terraform/etcd_tls.zip"

  source {
    filename = "ca.crt"
    content  = "${data.template_file.etcd_ca_cert_pem.rendered}"
  }

  source {
    filename = "server.crt"
    content  = "${join("", tls_locally_signed_cert.etcd_server.*.cert_pem)}"
  }

  source {
    filename = "server.key"
    content  = "${join("", tls_private_key.etcd_server.*.private_key_pem)}"
  }

  source {
    filename = "peer.crt"
    content  = "${join("", tls_locally_signed_cert.etcd_peer.*.cert_pem)}"
  }

  source {
    filename = "peer.key"
    content  = "${join("", tls_private_key.etcd_peer.*.private_key_pem)}"
  }

  source {
    filename = "client.crt"
    content  = "${data.template_file.etcd_client_crt.rendered}"
  }

  source {
    filename = "client.key"
    content  = "${data.template_file.etcd_client_key.rendered}"
  }
}
