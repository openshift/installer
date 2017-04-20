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

    etcd_servers   = "${join(",", var.etcd_servers)}"
    cloud_provider = "${var.cloud_provider}"

    cluster_cidr        = "${var.cluster_cidr}"
    service_cidr        = "${var.service_cidr}"
    kube_dns_service_ip = "${var.kube_dns_service_ip}"
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
  }
}

# Self-hosted bootstrapping manifests (resources/generated/manifests-bootstrap/)
resource "template_dir" "bootkube-bootstrap" {
  source_dir      = "${path.module}/resources/bootstrap-manifests"
  destination_dir = "${path.cwd}/generated/bootstrap-manifests"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"

    etcd_servers = "${join(",", var.etcd_servers)}"

    advertise_address = "${var.advertise_address}"
    cluster_cidr      = "${var.cluster_cidr}"
    service_cidr      = "${var.service_cidr}"
  }
}

# kubeconfig (resources/generated/kubeconfig)
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
  filename = "${path.cwd}/generated/kubeconfig"
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
