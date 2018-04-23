resource "local_file" "etcd_client_cert" {
  content  = "${tls_locally_signed_cert.etcd_client.cert_pem}"
  filename = "./generated/tls/etcd-client.crt"
}

data "ignition_file" "etcd_client_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${tls_locally_signed_cert.etcd_client.cert_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client.crt"
}

resource "local_file" "etcd_client_key" {
  content  = "${tls_private_key.etcd_client.private_key_pem}"
  filename = "./generated/tls/etcd-client.key"
}

data "ignition_file" "etcd_client_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${tls_private_key.etcd_client.private_key_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client.key"
}
