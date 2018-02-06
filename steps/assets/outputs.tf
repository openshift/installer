output "kube_dns_service_ip" {
  value = "${module.bootkube.kube_dns_service_ip}"
}

output "bootkube_service" {
  value = "${module.bootkube.systemd_service_rendered}"
}

output "bootkube_path_unit" {
  value = "${module.bootkube.systemd_path_unit_rendered}"
}

output "tectonic_service" {
  value = "${module.tectonic.systemd_service_rendered}"
}

output "tectonic_path_unit" {
  value = "${module.tectonic.systemd_path_unit_rendered}"
}

output "tectonic_bucket" {
  value = "${aws_s3_bucket_object.tectonic_assets.bucket}"
}

output "tectonic_key" {
  value = "${aws_s3_bucket_object.tectonic_assets.key}"
}

output "kubeconfig_bucket" {
  value = "${aws_s3_bucket_object.kubeconfig.bucket}"
}

output "kubeconfig_key" {
  value = "${aws_s3_bucket_object.kubeconfig.key}"
}

output "kubeconfig_content" {
  value = "${module.bootkube.kubeconfig}"
}

output "s3_bucket" {
  value = "${aws_s3_bucket.tectonic.bucket}"
}

output "s3_bucket_domain_name" {
  value = "${aws_s3_bucket.tectonic.bucket_domain_name}"
}

output "cluster_id" {
  value = "${module.tectonic.cluster_id}"
}

// TLS
output "etcd_ca_crt_pem" {
  value = "${module.etcd_certs.etcd_ca_crt_pem}"
}

output "etcd_client_crt_pem" {
  value = "${module.etcd_certs.etcd_client_crt_pem}"
}

output "etcd_client_key_pem" {
  value = "${module.etcd_certs.etcd_client_key_pem}"
}

output "etcd_peer_crt_pem" {
  value = "${module.etcd_certs.etcd_peer_crt_pem}"
}

output "etcd_peer_key_pem" {
  value = "${module.etcd_certs.etcd_peer_key_pem}"
}

output "etcd_server_crt_pem" {
  value = "${module.etcd_certs.etcd_server_crt_pem}"
}

output "etcd_server_key_pem" {
  value = "${module.etcd_certs.etcd_server_key_pem}"
}

output "ingress_certs_ca_cert_pem" {
  value = "${module.ingress_certs.ca_cert_pem}"
}

output "kube_certs_ca_cert_pem" {
  value = "${module.kube_certs.ca_cert_pem}"
}
