data "ignition_config" "main" {
  files = [
    "${var.ign_installer_kubelet_env_id}",
    "${var.ign_installer_runtime_mappings_id}",
    "${var.ign_max_user_watches_id}",
    "${var.ign_s3_puller_id}",
    "${var.ign_ca_cert_id_list}",
  ]

  systemd = [
    "${var.ign_docker_dropin_id}",
    "${var.ign_k8s_node_bootstrap_service_id}",
    "${var.ign_kubelet_service_id}",
    "${var.ign_locksmithd_service_id}",
    "${var.ign_update_ca_certificates_dropin_id}",
  ]
}
