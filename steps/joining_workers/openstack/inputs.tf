// This could be encapsulated as a data source
data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

locals {
  worker_sg_id = "${data.terraform_remote_state.topology.worker_sg_id}"
}
