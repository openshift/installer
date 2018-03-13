data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.cwd}/assets.tfstate"
  }
}

locals {
  kubeconfig_kubelet_content = "${data.terraform_remote_state.assets.kubeconfig_kubelet_content}"
  ignition_bootstrap         = "${data.terraform_remote_state.assets.ignition_bootstrap}"
  ignition_etcd              = "${data.terraform_remote_state.assets.ignition_etcd}"
}
