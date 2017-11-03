data "archive_file" "etcd_tls_zip" {
  type = "zip"

  output_path = "./.terraform/etcd_tls.zip"

  source {
    filename = "ca.crt"
    content  = "${file(var.etcd_ca_crt_pem_path)}"
  }

  source {
    filename = "server.crt"
    content  = "${file(var.etcd_server_crt_pem_path)}"
  }

  source {
    filename = "server.key"
    content  = "${file(var.etcd_server_key_pem_path)}"
  }

  source {
    filename = "peer.crt"
    content  = "${file(var.etcd_peer_crt_pem_path)}"
  }

  source {
    filename = "peer.key"
    content  = "${file(var.etcd_peer_key_pem_path)}"
  }

  source {
    filename = "client.crt"
    content  = "${file(var.etcd_client_crt_pem_path)}"
  }

  source {
    filename = "client.key"
    content  = "${file(var.etcd_client_key_pem_path)}"
  }
}
