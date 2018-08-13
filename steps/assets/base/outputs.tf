output "kubeconfig_kubelet_content" {
  value = "${module.bootkube.kubeconfig-kubelet}"
}

# right now, etcd is currently the same per-platform.
output "ignition_etcd" {
  value = "${data.ignition_config.etcd.*.rendered}"
}

# TODO(cdc) clean this up, get rid of ignition_etcd
output "ignition_etcd_files" {
  value = ["${module.ignition_bootstrap.etcd_crt_id_list}"]
}

output "ignition_bootstrap" {
  value = "${data.ignition_config.bootstrap.rendered}"
}
