data "terraform_remote_state" "topology" {
  backend = "local"

  config {
    path = "${path.cwd}/topology.tfstate"
  }
}

locals {
  private_zone_id           = "${var.tectonic_aws_external_private_zone != "" ? var.tectonic_aws_external_private_zone : data.terraform_remote_state.topology.private_zone_id}"
  tnc_s3_bucket_domain_name = "${data.terraform_remote_state.topology.tnc_s3_bucket_domain_name}"
  tnc_elb_dns_name          = "${data.terraform_remote_state.topology.tnc_elb_dns_name}"
  tnc_elb_zone_id           = "${data.terraform_remote_state.topology.tnc_elb_zone_id}"
}
