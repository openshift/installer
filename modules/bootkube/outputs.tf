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
  value = "${sha1("${template_dir.bootkube-bootstrap.id} ${local_file.kubeconfig.id} ${local_file.bootkube.id} ${template_dir.bootkube.id} ${join(" ",local_file.etcd-operator.*.id,local_file.etcd-service.*.id,local_file.bootstrap-etcd.*.id)}")}"
}

output "kubeconfig" {
  value = "${data.template_file.kubeconfig.rendered}"
}

output "ca_cert" {
  value = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube-ca.*.cert_pem) : var.ca_cert}"
}

output "ca_key_alg" {
  value = "${var.ca_cert == "" ? join(" ", tls_self_signed_cert.kube-ca.*.key_algorithm) : var.ca_key_alg}"
}

output "ca_key" {
  value = "${var.ca_cert == "" ? join(" ", tls_private_key.kube-ca.*.private_key_pem) : var.ca_key}"
}

output "systemd_service" {
  value = "${data.template_file.bootkube_service.rendered}"
}

output "kube_dns_service_ip" {
  value = "${cidrhost(var.service_cidr, 10)}"
}
