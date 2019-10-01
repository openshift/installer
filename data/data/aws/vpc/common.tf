# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file

// Only reference data sources which are guaranteed to exist at any time (above) in this locals{} block
locals {
  // How many AZs to create subnets in
  new_az_count = length(var.availability_zones)

  // The VPC ID to use to build the rest of the vpc data sources
  vpc_id = aws_vpc.new_vpc.id

  private_subnet_ids = aws_subnet.private_subnet.*.id
  public_subnet_ids  = aws_subnet.public_subnet.*.id
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "aws.new_vpc" and "data.aws_vpc.cluster_vpc", just "data.aws_vpc.cluster_vpc" used everwhere).

data "aws_vpc" "cluster_vpc" {
  id = local.vpc_id
}

