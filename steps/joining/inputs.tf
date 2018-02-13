// This could be encapsulated as a data source
data "terraform_remote_state" "bootstrap" {
  backend = "local"

  config {
    path = "${path.module}/../../${var.tectonic_cluster_name}/bootstrap.tfstate"
  }
}

locals {
  aws_launch_configuration_masters = "${data.terraform_remote_state.bootstrap.aws_launch_configuration_masters}"
  subnet_ids_masters               = "${data.terraform_remote_state.bootstrap.subnet_ids_masters}"
  aws_lbs_masters                  = "${data.terraform_remote_state.bootstrap.aws_lbs_masters}"
  cluster_id                       = "${data.terraform_remote_state.bootstrap.cluster_id}"
  aws_launch_configuration_workers = "${data.terraform_remote_state.bootstrap.aws_launch_configuration_workers}"
  subnet_ids_workers               = "${data.terraform_remote_state.bootstrap.subnet_ids_workers}"
}
