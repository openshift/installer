// This could be encapsulated as a data source
data "terraform_remote_state" "bootstrap" {
  backend = "local"

  config {
    path = "${path.cwd}/bootstrap.tfstate"
  }
}

data "terraform_remote_state" "assets" {
  backend = "local"

  config {
    path = "${path.cwd}/assets.tfstate"
  }
}

locals {
  container_linux_version = "${data.terraform_remote_state.bootstrap.container_linux_version}"
  instance_count          = "${data.terraform_remote_state.bootstrap.etcd_instance_count}"
  ignition_etcd           = "${data.terraform_remote_state.assets.ignition_etcd}"
  sg_id                   = "${data.terraform_remote_state.bootstrap.etcd_sg_id}"
  subnet_ids_workers      = "${data.terraform_remote_state.bootstrap.subnet_ids_workers}"
  s3_bucket               = "${data.terraform_remote_state.bootstrap.s3_bucket}"
  private_zone_id         = "${data.terraform_remote_state.bootstrap.private_zone_id}"
  tnc_elb_dns_name        = "${data.terraform_remote_state.bootstrap.tnc_elb_dns_name}"
  tnc_elb_zone_id         = "${data.terraform_remote_state.bootstrap.tnc_elb_zone_id}"
}
