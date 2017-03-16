# These subnet data-sources import external subnets from their user-supplied subnet IDs
# whenever an external VPC is specified
#
data "aws_subnet" "external_worker" {
  count = "${var.tectonic_aws_external_vpc_id == "" ? 0 : var.tectonic_aws_az_count}"
  id    = "${var.tectonic_aws_external_vpc_worker_subnets[count.index]}"
}

data "aws_subnet" "external_master" {
  count = "${var.tectonic_aws_external_vpc_id == "" ? 0 : var.tectonic_aws_az_count}"
  id    = "${var.tectonic_aws_external_vpc_master_subnets[count.index]}"
}
