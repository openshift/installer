data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

locals {
  master_subnet_ids = "${data.terraform_remote_state.topology.subnet_ids_masters}"
  worker_subnet_ids = "${data.terraform_remote_state.topology.subnet_ids_workers}"
  master_sg_id      = "${data.terraform_remote_state.topology.master_sg_id}"
  worker_sg_id      = "${data.terraform_remote_state.topology.worker_sg_id}"
  aws_lbs           = "${data.terraform_remote_state.topology.aws_lbs}"

  private_zone_id = "${var.tectonic_aws_external_private_zone != "" ? var.tectonic_aws_external_private_zone : data.terraform_remote_state.topology.private_zone_id}"
}
