resource "local_file" "client_key" {
  content  = "${tls_private_key.identity_client.private_key_pem}"
  filename = "./generated/tls/identity-client.key"
}

resource "local_file" "client_cert" {
  content  = "${tls_locally_signed_cert.identity_client.cert_pem}"
  filename = "./generated/tls/identity-client.crt"
}

resource "local_file" "server_key" {
  content  = "${tls_private_key.identity_server.private_key_pem}"
  filename = "./generated/tls/identity-server.key"
}

resource "local_file" "server_cert" {
  content  = "${data.template_file.identity_server_chained.rendered}"
  filename = "./generated/tls/identity-server.crt"
}
