provider "aws" {
  region  = "${var.region}"
  profile = "${var.profile}"
  version = "1.8.0"

  assume_role {
    role_arn = "${var.role_arn}"
  }
}

data "aws_availability_zones" "zones" {}

locals {
  zone_count     = "${length(data.aws_availability_zones.zones.names)}"
  zone_count_odd = "${local.zone_count % 2 == 0 ? local.zone_count - 1 : local.zone_count}"
}
