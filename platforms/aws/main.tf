data "aws_availability_zones" "azs" {}

module "vpc" {
  source                       = "../../modules/aws/vpc"
  tectonic_aws_external_vpc_id = "${var.tectonic_aws_external_vpc_id}"
  tectonic_aws_vpc_cidr_block  = "${var.tectonic_aws_vpc_cidr_block}"
  tectonic_cluster_name        = "${var.tectonic_cluster_name}"
  tectonic_aws_az_count        = "${var.tectonic_aws_az_count}"
}

module "etcd" {
  source = "../../modules/aws/etcd"

  vpc_id                = "${module.vpc.vpc_id}"
  node_count            = "${var.tectonic_aws_az_count == 5 ? 5 : 3}"
  ssh_key               = "${var.tectonic_aws_ssh_key}"
  dns_zone              = "${module.dns.int_zone_id}"
  etcd_subnets          = ["${module.vpc.worker_subnet_ids}"]
  tectonic_base_domain  = "${var.tectonic_base_domain}"
  tectonic_cluster_name = "${var.tectonic_cluster_name}"
  tectonic_cl_channel   = "${var.tectonic_cl_channel}"
  external_endpoints    = ["${var.tectonic_etcd_servers}"]
}

module "masters" {
  source = "../../modules/aws/master-asg"

  vpc_id                       = "${module.vpc.vpc_id}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  tectonic_base_domain         = "${var.tectonic_base_domain}"
  tectonic_cluster_name        = "${var.tectonic_cluster_name}"
  tectonic_cl_channel          = "${var.tectonic_cl_channel}"
  tectonic_master_count        = "${var.tectonic_master_count}"
  etcd_endpoints               = ["${module.etcd.endpoints}"]
  tectonic_aws_master_ec2_type = "${var.tectonic_aws_master_ec2_type}"
  extra_sg_ids                 = ["${module.vpc.cluster_default_sg}"]
  kube_image_url               = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
  kube_image_tag               = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
  master_subnet_ids            = ["${module.vpc.master_subnet_ids}"]
}

module "workers" {
  source = "../../modules/aws/worker-asg"

  vpc_id                       = "${module.vpc.vpc_id}"
  tectonic_worker_count        = "${var.tectonic_worker_count}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  etcd_endpoints               = "${module.etcd.endpoints}"
  tectonic_base_domain         = "${var.tectonic_base_domain}"
  tectonic_cluster_name        = "${var.tectonic_cluster_name}"
  worker_subnet_ids            = ["${module.vpc.worker_subnet_ids}"]
  tectonic_cl_channel          = "${var.tectonic_cl_channel}"
  tectonic_aws_worker_ec2_type = "${var.tectonic_aws_worker_ec2_type}"
  kube_image_url               = "${element(split(":", var.tectonic_container_images["hyperkube"]), 0)}"
  kube_image_tag               = "${element(split(":", var.tectonic_container_images["hyperkube"]), 1)}"
  extra_sg_ids                 = ["${module.vpc.cluster_default_sg}"]
}

module "dns" {
  source = "../../modules/aws/dns"

  vpc_id               = "${module.vpc.vpc_id}"
  tectonic_base_domain = "${var.tectonic_base_domain}"
  tectonic_dns_name    = "${var.tectonic_dns_name}"
  console-elb          = "${module.masters.console-elb}"
  api-internal-elb     = "${module.masters.api-internal-elb}"
  api-external-elb     = "${module.masters.api-external-elb}"
}
