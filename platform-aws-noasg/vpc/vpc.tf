variable "external_vpc_id" {
  type = "string"
}

variable "vpc_cid_block" {
  type = "string"
}

variable "cluster_name" {
  type = "string"
}

resource "aws_vpc" "new_vpc" {
  count                = "${length(var.external_vpc_id) > 0 ? 0 : 1}"
  cidr_block           = "${var.vpc_cid_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Name = "${var.cluster_name}"
  }
}

output "vpc_id" {
  value = "${length(var.external_vpc_id) > 0 ? var.external_vpc_id : aws_vpc.new_vpc.id}"
}
