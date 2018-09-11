output "ca_cert_pem_list" {
  value = [
    "${var.root_ca_cert_pem}",
    "${var.ingress_ca_cert_pem}",
    "${var.etcd_ca_cert_pem}",
  ]
}

output "etcd_crt_id_list" {
  value = ["${flatten(list(
    data.ignition_file.root_ca.*.id,
    data.ignition_file.etcd_ca.*.id,
  ))}"]
}

output "ignition_file_id_list" {
  value = [
    "${data.ignition_file.root_ca_cert_pem.id}",
    "${data.ignition_file.ingress_ca_cert_pem.id}",
    "${data.ignition_file.etcd_ca_cert_pem.id}",
    "${data.ignition_file.registries_config.id}",
  ]
}

output "ignition_systemd_id_list" {
  value = [
    "${data.ignition_systemd_unit.kubelet.id}",
  ]
}
