output "kubeconfig_kubelet_content" {
  value = "${module.bootkube.kubeconfig-kubelet}"
}

# For the bootstrap nodes, let them build their ignition configuration
# as they see fit.
output "ignition_bootstrap_files" {
  value = ["${compact(flatten(list(
    list(
      data.ignition_file.kube-system_cluster_config.id,
      data.ignition_file.tectonic_cluster_config.id,
      data.ignition_file.tnco_config.id,
      data.ignition_file.kco_config.id,
      data.ignition_file.bootstrap_kubeconfig.id,
    ),
    module.ignition_bootstrap.ignition_file_id_list,
    module.bootkube.ignition_file_id_list,
    module.tectonic.ignition_file_id_list,
    module.ca_certs.ignition_file_id_list,
    module.kube_certs.ignition_file_id_list,
    module.etcd_certs.ignition_file_id_list,
   )))}"]
}

output "ignition_bootstrap_systemd" {
  value = ["${compact(flatten(list(
    list(
      module.bootkube.systemd_service_id,
      module.bootkube.systemd_path_unit_id,
      module.tectonic.systemd_service_id,
      module.tectonic.systemd_path_unit_id,
    ),
    module.ignition_bootstrap.ignition_systemd_id_list,
   )))}"]
}

# right now, etcd is currently the same per-platform.
output "ignition_etcd" {
  value = "${data.ignition_config.etcd.*.rendered}"
}
