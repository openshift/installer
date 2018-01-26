output "aggregator_ca_cert_pem" {
  value = "${file(var.aggregator_ca_cert_pem_path)}"
}

output "ca_cert_pem" {
  value = "${file(var.ca_cert_pem_path)}"
}

output "kubelet_cert_pem" {
  value = "${file(var.kubelet_cert_pem_path)}"
}

output "kubelet_key_pem" {
  value = "${file(var.kubelet_key_pem_path)}"
}

output "apiserver_cert_pem" {
  value = "${file(var.apiserver_cert_pem_path)}"
}

output "apiserver_key_pem" {
  value = "${file(var.apiserver_key_pem_path)}"
}

output "apiserver_proxy_cert_pem" {
  value = "${file(var.apiserver_proxy_cert_pem_path)}"
}

output "apiserver_proxy_key_pem" {
  value = "${file(var.apiserver_proxy_key_pem_path)}"
}

output "id" {
  value = "${sha1("
  ${join(" ",
    list(local_file.aggregator_ca_crt.id,
    local_file.apiserver_key.id,
    local_file.apiserver_crt.id,
    local_file.apiserver_proxy_key.id,
    local_file.apiserver_proxy_crt.id,
    local_file.kube_ca_crt.id,
    local_file.kubelet_key.id,
    local_file.kubelet_crt.id,)
    )}
  ")}"
}
