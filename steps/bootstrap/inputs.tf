data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.module}/../../${var.tectonic_cluster_name}/assets.tfstate"
  }
}

locals {
  cluster_id                 = "${data.terraform_remote_state.assets.cluster_id}"
  kubeconfig_kubelet_content = "${data.terraform_remote_state.assets.kubeconfig_kubelet_content}"
}
