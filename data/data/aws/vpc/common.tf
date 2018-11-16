# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file
data "aws_region" "current" {}

// Fetch a list of available AZs
data "aws_availability_zones" "azs" {}

// Only reference data sources which are gauranteed to exist at any time (above) in this locals{} block
locals {
  // List of possible AZs for each type of subnet
  new_worker_subnet_azs = ["${coalescelist(keys(var.new_worker_subnet_configs), data.aws_availability_zones.azs.names)}"]
  new_master_subnet_azs = ["${coalescelist(keys(var.new_master_subnet_configs), data.aws_availability_zones.azs.names)}"]

  // How many AZs to create worker and master subnets in
  new_worker_az_count = "${length(local.new_worker_subnet_azs)}"
  new_master_az_count = "${length(local.new_master_subnet_azs)}"

  // The VPC ID to use to build the rest of the vpc data sources
  vpc_id = "${aws_vpc.new_vpc.id}"

  // When referencing the _ids arrays or data source arrays via count = , always use the *_count variable rather than taking the length of the list
  worker_subnet_ids   = "${aws_subnet.worker_subnet.*.id}"
  master_subnet_ids   = "${aws_subnet.master_subnet.*.id}"
  worker_subnet_count = "${local.new_worker_az_count}"
  master_subnet_count = "${local.new_master_az_count}"
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "aws.new_vpc" and "data.aws_vpc.cluster_vpc", just "data.aws_vpc.cluster_vpc" used everwhere).

data "aws_vpc" "cluster_vpc" {
  id = "${local.vpc_id}"
}
