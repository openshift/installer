// This could be encapsulated as a data source
data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

locals {
  etcd_sg_id = "${data.terraform_remote_state.topology.etcd_sg_id}"
}
