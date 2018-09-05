data "terraform_remote_state" "infra" {
  backend = "local"

  config {
    path = "${path.cwd}/infra.tfstate"
  }
}

data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.cwd}/assets.tfstate"
  }
}

locals {
  libvirt_network_id     = "${data.terraform_remote_state.infra.libvirt_network_id}"
  libvirt_base_volume_id = "${data.terraform_remote_state.infra.libvirt_base_volume_id}"

  ignition_bootstrap = "${data.terraform_remote_state.assets.ignition_bootstrap}"
}
