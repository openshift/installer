provider "aws" {
  region = "${var.tectonic_aws_region}"
}

data "aws_availability_zones" "azs" {}

module "vpc" {
  source = "../../modules/aws/vpc"

  cidr_block   = "${var.tectonic_aws_vpc_cidr_block}"
  cluster_name = "${var.tectonic_cluster_name}"

  external_vpc_id         = "${var.tectonic_aws_external_vpc_id}"
  external_master_subnets = ["${compact(var.tectonic_aws_external_master_subnet_ids)}"]
  external_worker_subnets = ["${compact(var.tectonic_aws_external_worker_subnet_ids)}"]
  cluster_id              = "${module.tectonic.cluster_id}"
  extra_tags              = "${var.tectonic_aws_extra_tags}"
  enable_etcd_sg          = "${!var.tectonic_experimental && length(compact(var.tectonic_etcd_servers)) == 0 ? 1 : 0}"

  # VPC layout settings.
  #
  # The following parameters control the layout of the VPC accross availability zones.
  # Two modes are available:
  # A. Explicitly configure a list of AZs + associated subnet CIDRs
  # B. Let the module calculate subnets accross a set number of AZs
  #
  # To enable mode A, configure a set of AZs + CIDRs for masters and workers using the
  # "tectonic_aws_master_custom_subnets" and "tectonic_aws_worker_custom_subnets" variables.
  #
  # To enable mode B, make sure that "tectonic_aws_master_custom_subnets" and "tectonic_aws_worker_custom_subnets" 
  # ARE NOT SET.

  # These counts could be deducted by length(keys(var.tectonic_aws_master_custom_subnets)) 
  # but there is a restriction on passing computed values as counts. This approach works around that.
  master_az_count = "${length(keys(var.tectonic_aws_master_custom_subnets)) > 0 ? "${length(keys(var.tectonic_aws_master_custom_subnets))}" : "${length(data.aws_availability_zones.azs.names)}"}"
  worker_az_count = "${length(keys(var.tectonic_aws_worker_custom_subnets)) > 0 ? "${length(keys(var.tectonic_aws_worker_custom_subnets))}" : "${length(data.aws_availability_zones.azs.names)}"}"
  # The appending of the "padding" element is required as workaround since the function
  # element() won't work on empty lists. See https://github.com/hashicorp/terraform/issues/11210
  master_subnets = "${concat(values(var.tectonic_aws_master_custom_subnets),list("padding"))}"
  worker_subnets = "${concat(values(var.tectonic_aws_worker_custom_subnets),list("padding"))}"
  # The split() / join() trick works around the limitation of ternary operator expressions 
  # only being able to return strings.
  master_azs = ["${ split("|", "${length(keys(var.tectonic_aws_master_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_master_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"]
  worker_azs = ["${ split("|", "${length(keys(var.tectonic_aws_worker_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_worker_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"]
}

module "etcd" {
  source = "../../modules/aws/etcd"

  instance_count = "${var.tectonic_experimental ? 0 : var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : length(data.aws_availability_zones.azs.names) == 5 ? 5 : 3}"
  ec2_type       = "${var.tectonic_aws_etcd_ec2_type}"
  sg_ids         = "${concat(var.tectonic_aws_etcd_extra_sg_ids, list(module.vpc.etcd_sg_id))}"

  ssh_key         = "${var.tectonic_aws_ssh_key}"
  cl_channel      = "${var.tectonic_cl_channel}"
  container_image = "${var.tectonic_container_images["etcd"]}"

  subnets = ["${module.vpc.worker_subnet_ids}"]

  dns_zone_id  = "${var.tectonic_aws_external_private_zone == "" ? join("", aws_route53_zone.tectonic-int.*.zone_id) : var.tectonic_aws_external_private_zone}"
  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  external_endpoints = ["${compact(var.tectonic_etcd_servers)}"]
  cluster_id         = "${module.tectonic.cluster_id}"
  extra_tags         = "${var.tectonic_aws_extra_tags}"

  root_volume_type = "${var.tectonic_aws_etcd_root_volume_type}"
  root_volume_size = "${var.tectonic_aws_etcd_root_volume_size}"
  root_volume_iops = "${var.tectonic_aws_etcd_root_volume_iops}"

  dns_enabled = "${!var.tectonic_experimental && length(compact(var.tectonic_etcd_servers)) == 0}"
  tls_enabled = "${var.tectonic_etcd_tls_enabled}"

  tls_ca_crt_pem     = "${module.bootkube.etcd_ca_crt_pem}"
  tls_client_crt_pem = "${module.bootkube.etcd_client_crt_pem}"
  tls_client_key_pem = "${module.bootkube.etcd_client_key_pem}"
  tls_peer_crt_pem   = "${module.bootkube.etcd_peer_crt_pem}"
  tls_peer_key_pem   = "${module.bootkube.etcd_peer_key_pem}"
}

module "ignition-masters" {
  source = "../../modules/aws/ignition"

  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  kubeconfig_s3_location    = "${aws_s3_bucket_object.kubeconfig.bucket}/${aws_s3_bucket_object.kubeconfig.key}"
  assets_s3_location        = "${aws_s3_bucket_object.tectonic-assets.bucket}/${aws_s3_bucket_object.tectonic-assets.key}"
  container_images          = "${var.tectonic_container_images}"
  bootkube_service          = "${module.bootkube.systemd_service}"
  tectonic_service          = "${module.tectonic.systemd_service}"
  tectonic_service_disabled = "${var.tectonic_vanilla_k8s}"
  cluster_name              = "${var.tectonic_cluster_name}"
  image_re                  = "${var.tectonic_image_re}"
}

module "masters" {
  source = "../../modules/aws/master-asg"

  instance_count  = "${var.tectonic_master_count}"
  ec2_type        = "${var.tectonic_aws_master_ec2_type}"
  cluster_name    = "${var.tectonic_cluster_name}"
  master_iam_role = "${var.tectonic_aws_master_iam_role_name}"

  subnet_ids = ["${module.vpc.master_subnet_ids}"]

  master_sg_ids  = "${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  api_sg_ids     = ["${module.vpc.api_sg_id}"]
  console_sg_ids = ["${module.vpc.console_sg_id}"]

  ssh_key    = "${var.tectonic_aws_ssh_key}"
  cl_channel = "${var.tectonic_cl_channel}"
  user_data  = "${module.ignition-masters.ignition}"

  internal_zone_id             = "${var.tectonic_aws_external_private_zone == "" ? join("", aws_route53_zone.tectonic-int.*.zone_id) : var.tectonic_aws_external_private_zone}"
  external_zone_id             = "${join("", data.aws_route53_zone.tectonic-ext.*.zone_id)}"
  base_domain                  = "${var.tectonic_base_domain}"
  public_vpc                   = "${var.tectonic_aws_external_vpc_public}"
  cluster_id                   = "${module.tectonic.cluster_id}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  custom_dns_name              = "${var.tectonic_dns_name}"

  root_volume_type = "${var.tectonic_aws_master_root_volume_type}"
  root_volume_size = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_iops = "${var.tectonic_aws_master_root_volume_iops}"
}

module "ignition-workers" {
  source = "../../modules/aws/ignition"

  kubelet_node_label     = "node-role.kubernetes.io/node"
  kubelet_node_taints    = ""
  kube_dns_service_ip    = "${module.bootkube.kube_dns_service_ip}"
  kubeconfig_s3_location = "${aws_s3_bucket_object.kubeconfig.bucket}/${aws_s3_bucket_object.kubeconfig.key}"
  assets_s3_location     = ""
  container_images       = "${var.tectonic_container_images}"
  bootkube_service       = ""
  tectonic_service       = ""
  cluster_name           = ""
  image_re               = "${var.tectonic_image_re}"
}

module "workers" {
  source = "../../modules/aws/worker-asg"

  instance_count  = "${var.tectonic_worker_count}"
  ec2_type        = "${var.tectonic_aws_worker_ec2_type}"
  cluster_name    = "${var.tectonic_cluster_name}"
  worker_iam_role = "${var.tectonic_aws_worker_iam_role_name}"

  vpc_id     = "${module.vpc.vpc_id}"
  subnet_ids = ["${module.vpc.worker_subnet_ids}"]
  sg_ids     = "${concat(var.tectonic_aws_worker_extra_sg_ids, list(module.vpc.worker_sg_id))}"

  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  cl_channel                   = "${var.tectonic_cl_channel}"
  user_data                    = "${module.ignition-workers.ignition}"
  cluster_id                   = "${module.tectonic.cluster_id}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"

  root_volume_type = "${var.tectonic_aws_worker_root_volume_type}"
  root_volume_size = "${var.tectonic_aws_worker_root_volume_size}"
  root_volume_iops = "${var.tectonic_aws_worker_root_volume_iops}"
}
