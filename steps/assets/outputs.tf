output "kubeconfig_kubelet_content" {
  value = "${module.bootkube.kubeconfig-kubelet}"
}

output "ignition_bootstrap" {
  value = "${data.ignition_config.bootstrap.rendered}"
}

output "ignition_etcd" {
  value = "${data.ignition_config.etcd.*.rendered}"
}
