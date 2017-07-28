# This output is meant to be used to inject a dependency on the generated
# assets. As of TerraForm v0.9, it is difficult to make a module depend on
# another module (no depends_on, no triggers), or to make a data source
# depend on a module (no depends_on, no triggers, generally no dummy variable).
#
# For instance, using the 'archive_file' data source against the generated
# assets, which is a common use-case, is tricky. There is no mechanism for
# defining explicit dependencies and the only available variables are for the
# source, destination and archive type, leaving little opportunities for us to
# inject a dependency. Thanks to the property described below, this output can
# be used as part of the output filename, in order to enforce the creation of
# the archive after the assets have been properly generated.
#
# Both localfile and template_dir providers compute their IDs by hashing
# the content of the resources on disk. Because this output is computed from the
# combination of all the resources' IDs, it can't be guessed and can only be
# interpolated once the assets have all been created.
output "id" {
  value = "${sha1("
  ${data.archive_file.etcd_tls_zip.id}
  ${local_file.kubeconfig.id}
  ${local_file.bootkube_sh.id}
  ${template_dir.bootkube.id} ${template_dir.bootkube_bootstrap.id}
  ${join(" ",
    local_file.etcd_ca_crt.*.id,
    local_file.etcd_server_crt.*.id,
    local_file.etcd_server_key.*.id,
    local_file.etcd_client_crt.*.id,
    local_file.etcd_client_key.*.id,
    local_file.etcd_peer_crt.*.id,
    local_file.etcd_peer_key.*.id,
    template_dir.experimental.*.id,
    template_dir.bootstrap_experimental.*.id,
    template_dir.etcd_experimental.*.id,
    )}
  ")}"
}

output "etcd_tls_zip" {
  value = "${data.archive_file.etcd_tls_zip.id != "" ? file("./.terraform/etcd_tls.zip") : ""}"
}

output "kubeconfig" {
  value = "${data.template_file.kubeconfig.rendered}"
}

output "ca_cert" {
  value = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.cert_pem) : var.ca_cert}"
}

output "ca_key_alg" {
  value = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube_ca.*.key_algorithm) : var.ca_key_alg}"
}

output "ca_key" {
  value = "${var.ca_cert == "" ? join(" ", tls_private_key.kube_ca.*.private_key_pem) : var.ca_key}"
}

output "systemd_service" {
  value = "${data.template_file.bootkube_service.rendered}"
}

output "kube_dns_service_ip" {
  value = "${cidrhost(var.service_cidr, 10)}"
}

output "etcd_ca_crt_pem" {
  value = "${join("", tls_self_signed_cert.etcd_ca.*.cert_pem)}"
}

output "etcd_server_crt_pem" {
  value = "${join("", tls_locally_signed_cert.etcd_server.*.cert_pem)}"
}

output "etcd_server_key_pem" {
  value = "${join("", tls_private_key.etcd_server.*.private_key_pem)}"
}

output "etcd_client_crt_pem" {
  value = "${join("", tls_locally_signed_cert.etcd_client.*.cert_pem)}"
}

output "etcd_client_key_pem" {
  value = "${join("", tls_private_key.etcd_client.*.private_key_pem)}"
}

output "etcd_peer_crt_pem" {
  value = "${join("", tls_locally_signed_cert.etcd_peer.*.cert_pem)}"
}

output "etcd_peer_key_pem" {
  value = "${join("", tls_private_key.etcd_peer.*.private_key_pem)}"
}
