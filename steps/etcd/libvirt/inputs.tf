// This could be encapsulated as a data source
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
  ignition               = "${data.terraform_remote_state.assets.ignition_etcd}"
  libvirt_network_id     = "${data.terraform_remote_state.topology.libvirt_network_id}"
  libvirt_base_volume_id = "${data.terraform_remote_state.topology.libvirt_base_volume_id}"
}
