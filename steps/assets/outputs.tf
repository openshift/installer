output "kubeconfig_kubelet_content" {
  value = "${module.bootkube.kubeconfig-kubelet}"
}

output "cluster_id" {
  value = "${module.tectonic.cluster_id}"
}
