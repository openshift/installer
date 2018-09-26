# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file
data "aws_region" "current" {
  current = true
}

// Fetch a list of available AZs
data "aws_availability_zones" "azs" {}

// Only reference data sources which are gauranteed to exist at any time (above) in this locals{} block
locals {
  // Define canonical source of truth for this
  external_vpc_mode = "${var.external_vpc_id != ""}"

  // List of possible AZs for each type of subnet
  new_worker_subnet_azs = ["${coalescelist(keys(var.new_worker_subnet_configs), data.aws_availability_zones.azs.names)}"]
  new_master_subnet_azs = ["${coalescelist(keys(var.new_master_subnet_configs), data.aws_availability_zones.azs.names)}"]

  // How many AZs to create worker and master subnets in (always zero if external_vpc_mode)
  new_worker_az_count = "${local.external_vpc_mode ? 0 : length(local.new_worker_subnet_azs)}"
  new_master_az_count = "${local.external_vpc_mode ? 0 : length(local.new_master_subnet_azs)}"

  // The base set of ids needs to build rest of vpc data sources
  // This is crux of dealing with existing vpc / new vpc incongruity
  vpc_id = "${local.external_vpc_mode ? var.external_vpc_id : element(concat(aws_vpc.new_vpc.*.id,list("")),0)}"

  // When referencing the _ids arrays or data source arrays via count = , always use the *_count variable rather than taking the length of the list
  worker_subnet_ids   = ["${coalescelist(aws_subnet.worker_subnet.*.id,var.external_worker_subnet_ids)}"]
  master_subnet_ids   = ["${coalescelist(aws_subnet.master_subnet.*.id,var.external_master_subnet_ids)}"]
  worker_subnet_count = "${local.external_vpc_mode ? length(var.external_worker_subnet_ids) : local.new_worker_az_count}"
  master_subnet_count = "${local.external_vpc_mode ? length(var.external_master_subnet_ids) : local.new_master_az_count}"
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "data.aws_subnet.external-worker" and "data.aws_subnet.worker". just "data.aws_subnet.worker" used everwhere for list of worker subnets for any valid input var state)

data "aws_vpc" "cluster_vpc" {
  id = "${local.vpc_id}"
}

data "aws_subnet" "worker" {
  count  = "${local.worker_subnet_count}"
  id     = "${local.worker_subnet_ids[count.index]}"
  vpc_id = "${local.vpc_id}"
}

data "aws_subnet" "master" {
  count  = "${local.master_subnet_count}"
  id     = "${local.master_subnet_ids[count.index]}"
  vpc_id = "${local.vpc_id}"
}

data "aws_route_table" "worker" {
  count = "${local.worker_subnet_count}"

  filter = {
    name   = "association.subnet-id"
    values = ["${list(local.worker_subnet_ids[count.index])}"]
  }
}

data "aws_route_table" "master" {
  count = "${local.master_subnet_count}"

  filter = {
    name   = "association.subnet-id"
    values = ["${list(local.master_subnet_ids[count.index])}"]
  }
}
