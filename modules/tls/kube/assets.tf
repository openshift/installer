resource "local_file" "apiserver_key" {
  content  = "${tls_private_key.apiserver.private_key_pem}"
  filename = "./generated/tls/apiserver.key"
}

data "template_file" "apiserver_cert" {
  template = "${join("", list(tls_locally_signed_cert.apiserver.cert_pem, var.kube_ca_cert_pem))}"
}

resource "local_file" "apiserver_cert" {
  content  = "${data.template_file.apiserver_cert.rendered}"
  filename = "./generated/tls/apiserver.crt"
}

resource "local_file" "openshift_apiserver_key" {
  content  = "${tls_private_key.openshift_apiserver.private_key_pem}"
  filename = "./generated/tls/openshift-apiserver.key"
}

data "template_file" "openshift_apiserver_cert" {
  template = "${join("", list(tls_locally_signed_cert.openshift_apiserver.cert_pem, var.aggregator_ca_cert_pem))}"
}

resource "local_file" "openshift_apiserver_cert" {
  content  = "${data.template_file.openshift_apiserver_cert.rendered}"
  filename = "./generated/tls/openshift-apiserver.crt"
}

resource "local_file" "apiserver_proxy_key" {
  content  = "${tls_private_key.apiserver_proxy.private_key_pem}"
  filename = "./generated/tls/apiserver-proxy.key"
}

resource "local_file" "apiserver_proxy_cert" {
  content  = "${tls_locally_signed_cert.apiserver_proxy.cert_pem}"
  filename = "./generated/tls/apiserver-proxy.crt"
}

resource "local_file" "admin_key" {
  content  = "${tls_private_key.admin.private_key_pem}"
  filename = "./generated/tls/admin.key"
}

resource "local_file" "admin_cert" {
  content  = "${tls_locally_signed_cert.admin.cert_pem}"
  filename = "./generated/tls/admin.crt"
}

resource "local_file" "kubelet_key" {
  content  = "${tls_private_key.kubelet.private_key_pem}"
  filename = "./generated/tls/kubelet.key"
}

resource "local_file" "kubelet_cert" {
  content  = "${tls_locally_signed_cert.kubelet.cert_pem}"
  filename = "./generated/tls/kubelet.crt"
}
