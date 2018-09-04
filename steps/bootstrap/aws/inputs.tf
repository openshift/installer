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
  subnet_ids = "${data.terraform_remote_state.infra.subnet_ids_masters}"
  aws_lbs    = "${data.terraform_remote_state.infra.aws_lbs}"
  sg_id      = "${data.terraform_remote_state.infra.master_sg_id}"
  s3_bucket  = "${data.terraform_remote_state.infra.s3_bucket}"

  ignition_bootstrap = "${data.terraform_remote_state.assets.ignition_bootstrap}"
}
