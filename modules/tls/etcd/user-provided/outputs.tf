output "etcd_ca_crt_pem" {
  value = "${file(var.etcd_ca_crt_pem_path)}"
}

output "etcd_client_crt_pem" {
  value = "${file(var.etcd_client_crt_pem_path)}"
}

output "etcd_client_key_pem" {
  value = "${file(var.etcd_client_key_pem_path)}"
}

output "etcd_peer_crt_pem" {
  value = "${file(var.etcd_peer_crt_pem_path)}"
}

output "etcd_peer_key_pem" {
  value = "${file(var.etcd_peer_key_pem_path)}"
}

output "etcd_server_crt_pem" {
  value = "${file(var.etcd_server_crt_pem_path)}"
}

output "etcd_server_key_pem" {
  value = "${file(var.etcd_server_key_pem_path)}"
}

// The data.archive_file.etcd_tls_zip.id != "" assertion forces the etcd_tls_zip datasource to be run,
// hence ./.terraform/etcd_tls.zip must be generated and present.
output "etcd_tls_zip" {
  value = "${data.archive_file.etcd_tls_zip.id != "" ? file("./.terraform/etcd_tls.zip") : ""}"
}

output "id" {
  value = "${sha1("
  ${data.archive_file.etcd_tls_zip.id},
  ${join(" ",
    local_file.etcd_ca_crt.*.id,
    local_file.etcd_server_crt.*.id,
    local_file.etcd_server_key.*.id,
    local_file.etcd_client_crt.*.id,
    local_file.etcd_client_key.*.id,
    local_file.etcd_peer_crt.*.id,
    local_file.etcd_peer_key.*.id,
    )}
  ")}"
}
