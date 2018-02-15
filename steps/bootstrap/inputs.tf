data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.module}/../../${var.tectonic_cluster_name}/assets.tfstate"
  }
}

locals {
  etcd_ca_crt_pem            = "${data.terraform_remote_state.assets.etcd_ca_crt_pem}"
  etcd_client_crt_pem        = "${data.terraform_remote_state.assets.etcd_client_crt_pem}"
  etcd_client_key_pem        = "${data.terraform_remote_state.assets.etcd_client_key_pem}"
  etcd_peer_crt_pem          = "${data.terraform_remote_state.assets.etcd_peer_crt_pem}"
  etcd_peer_key_pem          = "${data.terraform_remote_state.assets.etcd_peer_key_pem}"
  etcd_server_crt_pem        = "${data.terraform_remote_state.assets.etcd_server_crt_pem}"
  etcd_server_key_pem        = "${data.terraform_remote_state.assets.etcd_server_key_pem}"
  ingress_certs_ca_cert_pem  = "${data.terraform_remote_state.assets.ingress_certs_ca_cert_pem}"
  kube_certs_ca_cert_pem     = "${data.terraform_remote_state.assets.kube_certs_ca_cert_pem}"
  kube_dns_service_ip        = "${data.terraform_remote_state.assets.kube_dns_service_ip}"
  cluster_id                 = "${data.terraform_remote_state.assets.cluster_id}"
  kubeconfig_kubelet_content = "${data.terraform_remote_state.assets.kubeconfig_kubelet_content}"
}
