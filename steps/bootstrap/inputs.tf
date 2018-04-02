data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

locals {
  subnet_ids_masters = "${data.terraform_remote_state.topology.subnet_ids_masters}"
  aws_lbs_masters    = "${data.terraform_remote_state.topology.aws_lbs_masters}"
  master_sg_id       = "${data.terraform_remote_state.topology.master_sg_id}"
}
