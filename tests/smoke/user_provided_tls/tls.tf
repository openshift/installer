module "etcd_certs" {
  source = "../../modules/tls/etcd/user-provided"

  etcd_ca_crt_pem_path     = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-ca.crt"
  etcd_server_crt_pem_path = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-server.crt"
  etcd_server_key_pem_path = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-server.key"
  etcd_peer_crt_pem_path   = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-peer.crt"
  etcd_peer_key_pem_path   = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-peer.key"
  etcd_client_crt_pem_path = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-client.crt"
  etcd_client_key_pem_path = "../../tests/smoke/user_provided_tls/certs/etcd/etcd-client.key"
}

module "identity_certs" {
  source = "../../modules/tls/identity/user-provided"

  client_key_pem_path  = "../../tests/smoke/user_provided_tls/certs/identity/identity-client.key"
  client_cert_pem_path = "../../tests/smoke/user_provided_tls/certs/identity/identity-client.crt"
  server_key_pem_path  = "../../tests/smoke/user_provided_tls/certs/identity/identity-server.key"
  server_cert_pem_path = "../../tests/smoke/user_provided_tls/certs/identity/identity-server.crt"
}

module "ingress_certs" {
  source = "../../modules/tls/ingress/user-provided"

  ca_cert_pem_path = "../../tests/smoke/user_provided_tls/certs/ingress/ca.crt"
  cert_pem_path    = "../../tests/smoke/user_provided_tls/certs/ingress/ingress.crt"
  key_pem_path     = "../../tests/smoke/user_provided_tls/certs/ingress/ingress.key"
}

module "kube_certs" {
  source = "../../modules/tls/kube/user-provided"

  aggregator_ca_cert_pem_path   = "../../tests/smoke/user_provided_tls/certs/kube/aggregator-ca.crt"
  ca_cert_pem_path              = "../../tests/smoke/user_provided_tls/certs/kube/ca.crt"
  ca_key_pem_path               = "../../tests/smoke/user_provided_tls/certs/kube/ca.key"
  admin_cert_pem_path           = "../../tests/smoke/user_provided_tls/certs/kube/admin.crt"
  admin_key_pem_path            = "../../tests/smoke/user_provided_tls/certs/kube/admin.key"
  apiserver_cert_pem_path       = "../../tests/smoke/user_provided_tls/certs/kube/apiserver.crt"
  apiserver_key_pem_path        = "../../tests/smoke/user_provided_tls/certs/kube/apiserver.key"
  apiserver_proxy_cert_pem_path = "../../tests/smoke/user_provided_tls/certs/kube/apiserver-proxy.crt"
  apiserver_proxy_key_pem_path  = "../../tests/smoke/user_provided_tls/certs/kube/apiserver-proxy.key"
}
