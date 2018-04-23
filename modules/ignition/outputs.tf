output "etcd_dropin_id_list" {
  value = "${data.ignition_systemd_unit.etcd.*.id}"
}

output "etcd_dropin_rendered_list" {
  value = "${data.template_file.etcd.*.rendered}"
}

output "update_ca_certificates_dropin_rendered" {
  value = "${data.template_file.update_ca_certificates_dropin.rendered}"
}

output "ca_cert_pem_list" {
  value = [
    "${var.root_ca_cert_pem}",
    "${var.ingress_ca_cert_pem}",
    "${var.etcd_ca_cert_pem}",
    "${var.custom_ca_cert_pem_list}",
  ]
}

output "etcd_crt_id_list" {
  value = ["${flatten(list(
    data.ignition_file.etcd_ca.*.id,
    data.ignition_file.etcd_client_key.*.id,
    data.ignition_file.etcd_client_crt.*.id,
    data.ignition_file.etcd_server_key.*.id,
    data.ignition_file.etcd_server_crt.*.id,
    data.ignition_file.etcd_peer_key.*.id,
    data.ignition_file.etcd_peer_crt.*.id,
  ))}"]
}

output "ignition_file_id_list" {
  value = ["${flatten(list(
    data.ignition_file.custom_ca_cert_pem.*.id,
    list(
      data.ignition_file.root_ca_cert_pem.id,
      data.ignition_file.ingress_ca_cert_pem.id,
      data.ignition_file.etcd_ca_cert_pem.id,
      data.ignition_file.installer_kubelet_env.id,
      data.ignition_file.installer_runtime_mappings.id,
      data.ignition_file.max_user_watches.id,
      data.ignition_file.profile_env.id,
      data.ignition_file.systemd_default_env.id,
    ),
  ))}"]
}

output "ignition_systemd_id_list" {
  value = [
    "${data.ignition_systemd_unit.docker_dropin.id}",
    "${data.ignition_systemd_unit.kubelet.id}",
    "${data.ignition_systemd_unit.locksmithd.id}",
    "${data.ignition_systemd_unit.k8s_node_bootstrap.id}",
    "${data.ignition_systemd_unit.update_ca_certificates_dropin.id}",
    "${data.ignition_systemd_unit.iscsi.id}",
    "${data.ignition_systemd_unit.rm_assets.id}",
  ]
}
