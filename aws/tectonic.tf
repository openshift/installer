data "aws_availability_zones" "azs" {}

module "vpc" {
  source          = "./vpc"
  external_vpc_id = "${var.external_vpc_id}"
  az_count        = "${var.az_count}"
}

module "etcd" {
  source      = "./etcd"
  vpc_id      = "${module.vpc.cluster_vpc_id}"
  node_count  = "${var.az_count}"
  etcd_domain = "${var.etcd_domain}"
}

output "etcd_endpoints" {
  value = "${module.etcd.endpoints}"
}
