# Canonical internal state definitions for this module.
# read only: only locals and data source definitions allowed. No resources or module blocks in this file
data "azure_region" "current" {}

// Fetch a list of available AZs
data "aws_availability_zones" "azs" {}

// Only reference data sources which are gauranteed to exist at any time (above) in this locals{} block
locals {
  // List of possible AZs for each type of subnet
  new_subnet_azs = "${data.aws_availability_zones.azs.names}"

  // How many AZs to create subnets in
  new_az_count = "${length(local.new_subnet_azs)}"

  // The VPC ID to use to build the rest of the vpc data sources
  vnet_id = "${azurerm_virtual_network.new_vnet.id}"

  // When referencing the _ids arrays or data source arrays via count = , always use the *_count variable rather than taking the length of the list
  private_subnet_ids   = "${azurerm_virtual_network..private_subnet.*.id}"
  public_subnet_ids    = "${aws_subnet.public_subnet.*.id}"
  private_subnet_count = "${local.new_az_count}"
  public_subnet_count  = "${local.new_az_count}"
}

# all data sources should be input variable-agnostic and used as canonical source for querying "state of resources" and building outputs
# (ie: we don't want "aws.new_vpc" and "data.aws_vpc.cluster_vpc", just "data.aws_vpc.cluster_vpc" used everwhere).

data "azure_vnet" "cluster_vnet" {
  id = "${local.vnet_id}"
}
