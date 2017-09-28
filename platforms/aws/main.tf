provider "aws" {
  region  = "${var.tectonic_aws_region}"
  version = "0.1.4"
}

data "aws_availability_zones" "azs" {}

module "vpc" {
  source = "../../modules/aws/vpc"

  cidr_block   = "${var.tectonic_aws_vpc_cidr_block}"
  base_domain  = "${var.tectonic_base_domain}"
  cluster_name = "${var.tectonic_cluster_name}"

  external_vpc_id         = "${var.tectonic_aws_external_vpc_id}"
  external_master_subnets = "${compact(var.tectonic_aws_external_master_subnet_ids)}"
  external_worker_subnets = "${compact(var.tectonic_aws_external_worker_subnet_ids)}"
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
  master_azs = "${ split("|", "${length(keys(var.tectonic_aws_master_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_master_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"
  worker_azs = "${ split("|", "${length(keys(var.tectonic_aws_worker_custom_subnets))}" > 0 ?
    join("|", keys(var.tectonic_aws_worker_custom_subnets)) :
    join("|", data.aws_availability_zones.azs.names)
  )}"
}

module "etcd" {
  source = "../../modules/aws/etcd"

  base_domain        = "${var.tectonic_base_domain}"
  cl_channel         = "${var.tectonic_cl_channel}"
  cluster_id         = "${module.tectonic.cluster_id}"
  cluster_name       = "${var.tectonic_cluster_name}"
  container_image    = "${var.tectonic_container_images["etcd"]}"
  dns_enabled        = "${!var.tectonic_experimental && length(compact(var.tectonic_etcd_servers)) == 0}"
  dns_zone_id        = "${var.tectonic_aws_private_endpoints ? data.null_data_source.zones.inputs["private"] : data.null_data_source.zones.inputs["public"]}"
  ec2_type           = "${var.tectonic_aws_etcd_ec2_type}"
  external_endpoints = "${compact(var.tectonic_etcd_servers)}"
  extra_tags         = "${var.tectonic_aws_extra_tags}"
  instance_count     = "${length(data.template_file.etcd_hostname_list.*.id)}"
  root_volume_iops   = "${var.tectonic_aws_etcd_root_volume_iops}"
  root_volume_size   = "${var.tectonic_aws_etcd_root_volume_size}"
  root_volume_type   = "${var.tectonic_aws_etcd_root_volume_type}"
  sg_ids             = "${concat(var.tectonic_aws_etcd_extra_sg_ids, list(module.vpc.etcd_sg_id))}"
  ssh_key            = "${var.tectonic_aws_ssh_key}"
  subnets            = "${module.vpc.worker_subnet_ids}"
  tls_enabled        = "${var.tectonic_etcd_tls_enabled}"
  tls_zip            = "${module.etcd_certs.etcd_tls_zip}"

  ign_etcd_dropin_id_list = "${module.ignition_masters.etcd_dropin_id_list}"
}

module "ignition_masters" {
  source = "../../modules/ignition"

  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider            = "aws"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_count                = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_tls_enabled          = "${var.tectonic_etcd_tls_enabled}"
  image_re                  = "${var.tectonic_image_re}"
  kube_dns_service_ip       = "${module.bootkube.kube_dns_service_ip}"
  kubeconfig_fetch_cmd      = "/opt/s3-puller.sh ${aws_s3_bucket_object.kubeconfig.bucket}/${aws_s3_bucket_object.kubeconfig.key} /etc/kubernetes/kubeconfig"
  kubelet_cni_bin_dir       = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  tectonic_vanilla_k8s      = "${var.tectonic_vanilla_k8s}"
}

module "masters" {
  source = "../../modules/aws/master-asg"

  api_sg_ids                   = ["${module.vpc.api_sg_id}"]
  assets_s3_location           = "${aws_s3_bucket_object.tectonic_assets.bucket}/${aws_s3_bucket_object.tectonic_assets.key}"
  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  base_domain                  = "${var.tectonic_base_domain}"
  cl_channel                   = "${var.tectonic_cl_channel}"
  cluster_id                   = "${module.tectonic.cluster_id}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  console_sg_ids               = ["${module.vpc.console_sg_id}"]
  container_images             = "${var.tectonic_container_images}"
  custom_dns_name              = "${var.tectonic_dns_name}"
  ec2_type                     = "${var.tectonic_aws_master_ec2_type}"
  external_zone_id             = "${data.null_data_source.zones.inputs["public"]}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  image_re                     = "${var.tectonic_image_re}"
  instance_count               = "${var.tectonic_master_count}"
  internal_zone_id             = "${data.null_data_source.zones.inputs["private"]}"
  master_iam_role              = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids                = "${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  private_endpoints            = "${var.tectonic_aws_private_endpoints}"
  public_endpoints             = "${var.tectonic_aws_public_endpoints}"
  root_volume_iops             = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size             = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type             = "${var.tectonic_aws_master_root_volume_type}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  subnet_ids                   = "${module.vpc.master_subnet_ids}"

  ign_bootkube_path_unit_id         = "${module.bootkube.systemd_path_unit_id}"
  ign_bootkube_service_id           = "${module.bootkube.systemd_service_id}"
  ign_docker_dropin_id              = "${module.ignition_masters.docker_dropin_id}"
  ign_installer_kubelet_env_id      = "${module.ignition_masters.installer_kubelet_env_id}"
  ign_k8s_node_bootstrap_service_id = "${module.ignition_masters.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id            = "${module.ignition_masters.kubelet_service_id}"
  ign_locksmithd_service_id         = "${module.ignition_masters.locksmithd_service_id}"
  ign_max_user_watches_id           = "${module.ignition_masters.max_user_watches_id}"
  ign_s3_puller_id                  = "${module.ignition_masters.s3_puller_id}"
  ign_tectonic_path_unit_id         = "${var.tectonic_vanilla_k8s ? "" : module.tectonic.systemd_path_unit_id}"
  ign_tectonic_service_id           = "${module.tectonic.systemd_service_id}"
}

module "ignition_workers" {
  source = "../../modules/ignition"

  bootstrap_upgrade_cl = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider       = "aws"
  container_images     = "${var.tectonic_container_images}"
  image_re             = "${var.tectonic_image_re}"
  kube_dns_service_ip  = "${module.bootkube.kube_dns_service_ip}"
  kubeconfig_fetch_cmd = "/opt/s3-puller.sh ${aws_s3_bucket_object.kubeconfig.bucket}/${aws_s3_bucket_object.kubeconfig.key} /etc/kubernetes/kubeconfig"
  kubelet_cni_bin_dir  = "${var.tectonic_calico_network_policy ? "/var/lib/cni/bin" : "" }"
  kubelet_node_label   = "node-role.kubernetes.io/node"
  kubelet_node_taints  = ""
  tectonic_vanilla_k8s = "${var.tectonic_vanilla_k8s}"
}

module "workers" {
  source = "../../modules/aws/worker-asg"

  autoscaling_group_extra_tags = "${var.tectonic_autoscaling_group_extra_tags}"
  cl_channel                   = "${var.tectonic_cl_channel}"
  cluster_id                   = "${module.tectonic.cluster_id}"
  cluster_name                 = "${var.tectonic_cluster_name}"
  ec2_type                     = "${var.tectonic_aws_worker_ec2_type}"
  extra_tags                   = "${var.tectonic_aws_extra_tags}"
  instance_count               = "${var.tectonic_worker_count}"
  root_volume_iops             = "${var.tectonic_aws_worker_root_volume_iops}"
  root_volume_size             = "${var.tectonic_aws_worker_root_volume_size}"
  root_volume_type             = "${var.tectonic_aws_worker_root_volume_type}"
  sg_ids                       = "${concat(var.tectonic_aws_worker_extra_sg_ids, list(module.vpc.worker_sg_id))}"
  ssh_key                      = "${var.tectonic_aws_ssh_key}"
  subnet_ids                   = "${module.vpc.worker_subnet_ids}"
  vpc_id                       = "${module.vpc.vpc_id}"
  worker_iam_role              = "${var.tectonic_aws_worker_iam_role_name}"
  load_balancers               = ["${var.tectonic_aws_worker_load_balancers}"]

  ign_docker_dropin_id              = "${module.ignition_workers.docker_dropin_id}"
  ign_installer_kubelet_env_id      = "${module.ignition_workers.installer_kubelet_env_id}"
  ign_k8s_node_bootstrap_service_id = "${module.ignition_workers.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id            = "${module.ignition_workers.kubelet_service_id}"
  ign_locksmithd_service_id         = "${module.ignition_masters.locksmithd_service_id}"
  ign_max_user_watches_id           = "${module.ignition_workers.max_user_watches_id}"
  ign_s3_puller_id                  = "${module.ignition_workers.s3_puller_id}"
}
