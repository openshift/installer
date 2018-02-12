resource "local_file" "aggregator_ca_crt" {
  content  = "${file(var.aggregator_ca_cert_pem_path)}"
  filename = "./generated/tls/aggregator-ca.crt"
}

resource "local_file" "apiserver_key" {
  content  = "${file(var.apiserver_key_pem_path)}"
  filename = "./generated/tls/apiserver.key"
}

resource "local_file" "apiserver_crt" {
  content  = "${file(var.apiserver_cert_pem_path)}"
  filename = "./generated/tls/apiserver.crt"
}

resource "local_file" "apiserver_proxy_key" {
  content  = "${file(var.apiserver_proxy_key_pem_path)}"
  filename = "./generated/tls/apiserver-proxy.key"
}

resource "local_file" "apiserver_proxy_crt" {
  content  = "${file(var.apiserver_proxy_cert_pem_path)}"
  filename = "./generated/tls/apiserver-proxy.crt"
}

resource "local_file" "kube_ca_crt" {
  content  = "${file(var.ca_cert_pem_path)}"
  filename = "./generated/tls/ca.crt"
}

resource "local_file" "kube_ca_key" {
  content  = "${file(var.ca_key_pem_path)}"
  filename = "./generated/tls/ca.key"
}

resource "local_file" "admin_key" {
  content  = "${file(var.admin_key_pem_path)}"
  filename = "./generated/tls/admin.key"
}

resource "local_file" "admin_crt" {
  content  = "${file(var.admin_cert_pem_path)}"
  filename = "./generated/tls/admin.crt"
}
