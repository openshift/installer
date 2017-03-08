variable "tectonic_aws_external_vpc_id" {
  type = "string"
}

variable "tectonic_aws_vpc_cidr_block" {
  type = "string"
}

variable "tectonic_cluster_name" {
  type = "string"
}

resource "aws_vpc" "new_vpc" {
  count                = "${length(var.tectonic_aws_external_vpc_id) > 0 ? 0 : 1}"
  cidr_block           = "${var.tectonic_aws_vpc_cidr_block}"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags {
    Name = "${var.tectonic_cluster_name}"
    KubernetesCluster = "${var.tectonic_cluster_name}"
  }
}

output "vpc_id" {
  value = "${length(var.tectonic_aws_external_vpc_id) > 0 ? var.tectonic_aws_external_vpc_id : aws_vpc.new_vpc.id}"
}
