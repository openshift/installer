# Self-hosted manifests (resources/generated/manifests/)
resource "template_folder" "bootkube" {
  input_path = "${path.module}/resources/manifests"
  output_path = "${path.cwd}/generated/manifests"

  vars {
    hyperkube_image = "${var.container_images["hyperkube"]}"
    pod_checkpointer_image = "${var.container_images["pod_checkpointer"]}"
    kubedns_image         = "${var.container_images["kubedns"]}"
    kubednsmasq_image     = "${var.container_images["kubednsmasq"]}"
    dnsmasq_metrics_image = "${var.container_images["dnsmasq_metrics"]}"
    exechealthz_image     = "${var.container_images["exechealthz"]}"
    flannel_image         = "${var.container_images["flannel"]}"

    etcd_servers = "${join(",", var.etcd_servers)}"
    cloud_provider = "${var.cloud_provider}"

    cluster_cidr = "${var.cluster_cidr}"
    service_cidr = "${var.service_cidr}"
    kube_dns_service_ip = "${var.kube_dns_service_ip}"
    advertise_address = "${var.advertise_address}"

    anonymous_auth = "${var.anonymous_auth}"
    oidc_issuer_url = "${var.oidc_issuer_url}"
    oidc_client_id = "${var.oidc_client_id}"
    oidc_username_claim = "${var.oidc_username_claim}"
    oidc_groups_claim = "${var.oidc_groups_claim}"

    ca_cert = "${base64encode(tls_self_signed_cert.kube-ca.cert_pem)}"
    apiserver_key = "${base64encode(tls_private_key.apiserver.private_key_pem)}"
    apiserver_cert = "${base64encode(tls_locally_signed_cert.apiserver.cert_pem)}"
    serviceaccount_pub = "${base64encode(tls_private_key.service-account.public_key_pem)}"
    serviceaccount_key = "${base64encode(tls_private_key.service-account.private_key_pem)}"
  }
}

# kubeconfig (resources/generated/kubeconfig)
data "template_file" "kubeconfig" {
  template = "${file("${path.module}/resources/kubeconfig")}"

  vars {
    ca_cert = "${base64encode(tls_self_signed_cert.kube-ca.cert_pem)}"
    kubelet_cert = "${base64encode(tls_locally_signed_cert.kubelet.cert_pem)}"
    kubelet_key = "${base64encode(tls_private_key.kubelet.private_key_pem)}"
    server = "${var.kube_apiserver_url}"
  }
}

resource "localfile_file" "kubeconfig" {
  content = "${data.template_file.kubeconfig.rendered}"
  destination = "${path.cwd}/generated/kubeconfig"
}

# bootkube.sh (resources/generated/bootkube.sh)
data "template_file" "bootkube" {
  template = "${file("${path.module}/resources/bootkube.sh")}"

  vars {
    bootkube_image = "${var.container_images["bootkube"]}"
    etcd_server = "${element(var.etcd_servers, 0)}"
  }
}

resource "localfile_file" "bootkube" {
  content = "${data.template_file.bootkube.rendered}"
  destination = "${path.cwd}/generated/bootkube.sh"
}