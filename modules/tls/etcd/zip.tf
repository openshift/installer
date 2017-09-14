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
