locals {
  auth_dir                = "./generated/auth"
  kubeconfig_kubelet_path = "${local.auth_dir}/kubeconfig-kubelet"
  kubeconfig_path         = "${local.auth_dir}/kubeconfig"
}

# some files we want rendered to disk for other use
resource "local_file" "kubeconfig-kubelet" {
  content  = "${module.bootkube.kubeconfig-kubelet_rendered}"
  filename = "${local.kubeconfig_kubelet_path}"
}

resource "local_file" "kubeconfig" {
  content  = "${module.bootkube.kubeconfig_rendered}"
  filename = "${local.kubeconfig_path}"
}
