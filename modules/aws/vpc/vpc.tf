data "aws_availability_zones" "azs" {}

resource "aws_vpc" "new_vpc" {
  count                = "${length(var.tectonic_aws_external_vpc_id) > 0 ? 0 : 1}"
  cidr_block           = "${var.tectonic_aws_vpc_cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Name              = "${var.tectonic_cluster_name}"
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }
}

data "aws_vpc" "cluster_vpc" {
  id = "${var.tectonic_aws_external_vpc_id == "" ? aws_vpc.new_vpc.id : var.tectonic_aws_external_vpc_id }"
}
