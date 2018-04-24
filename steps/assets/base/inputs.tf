data "terraform_remote_state" "tls" {
  backend = "local"

  config {
    path = "${path.cwd}/tls.tfstate"
  }
}

locals {
  admin_cert_pem           = "${data.terraform_remote_state.tls.admin_cert_pem}"
  admin_key_pem            = "${data.terraform_remote_state.tls.admin_key_pem}"
  aggregator_ca_cert_pem   = "${data.terraform_remote_state.tls.aggregator_ca_cert_pem}"
  aggregator_ca_key_pem    = "${data.terraform_remote_state.tls.aggregator_ca_key_pem}"
  apiserver_cert_pem       = "${data.terraform_remote_state.tls.apiserver_cert_pem}"
  apiserver_key_pem        = "${data.terraform_remote_state.tls.apiserver_key_pem}"
  apiserver_proxy_cert_pem = "${data.terraform_remote_state.tls.apiserver_proxy_cert_pem}"
  apiserver_proxy_key_pem  = "${data.terraform_remote_state.tls.apiserver_proxy_key_pem}"
  etcd_ca_cert_pem         = "${data.terraform_remote_state.tls.etcd_ca_cert_pem}"
  etcd_ca_key_pem          = "${data.terraform_remote_state.tls.etcd_ca_key_pem}"
  etcd_client_cert_pem     = "${data.terraform_remote_state.tls.etcd_client_cert_pem}"
  etcd_client_key_pem      = "${data.terraform_remote_state.tls.etcd_client_key_pem}"
  identity_client_ca_cert  = "${data.terraform_remote_state.tls.identity_client_ca_cert}"
  identity_client_cert_pem = "${data.terraform_remote_state.tls.identity_client_cert_pem}"
  identity_client_key_pem  = "${data.terraform_remote_state.tls.identity_client_key_pem}"
  identity_server_ca_cert  = "${data.terraform_remote_state.tls.identity_server_ca_cert}"
  identity_server_cert_pem = "${data.terraform_remote_state.tls.identity_server_cert_pem}"
  identity_server_key_pem  = "${data.terraform_remote_state.tls.identity_server_key_pem}"
  ingress_ca_cert_pem      = "${data.terraform_remote_state.tls.ingress_ca_cert_pem}"
  ingress_cert_pem         = "${data.terraform_remote_state.tls.ingress_cert_pem}"
  ingress_key_pem          = "${data.terraform_remote_state.tls.ingress_key_pem}"
  kube_ca_cert_pem         = "${data.terraform_remote_state.tls.kube_ca_cert_pem}"
  kube_ca_key_pem          = "${data.terraform_remote_state.tls.kube_ca_key_pem}"
  kubelet_cert_pem         = "${data.terraform_remote_state.tls.kubelet_cert_pem}"
  kubelet_key_pem          = "${data.terraform_remote_state.tls.kubelet_key_pem}"
  oidc_ca_cert             = "${data.terraform_remote_state.tls.oidc_ca_cert}"
  root_ca_cert_pem         = "${data.terraform_remote_state.tls.root_ca_cert_pem}"
}
