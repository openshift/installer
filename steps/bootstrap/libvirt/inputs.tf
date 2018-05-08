data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.cwd}/assets.tfstate"
  }
}

locals {
  libvirt_network_id     = "${data.terraform_remote_state.topology.libvirt_network_id}"
  libvirt_base_volume_id = "${data.terraform_remote_state.topology.libvirt_base_volume_id}"

  ignition_bootstrap = "${data.terraform_remote_state.assets.ignition_bootstrap}"
}
