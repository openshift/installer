provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

data "aws_availability_zones" "azs" {}

data "template_file" "etcd_hostname_list" {
  count    = "${var.tectonic_etcd_count > 0 ? var.tectonic_etcd_count : length(data.aws_availability_zones.azs.names) == 5 ? 5 : 3}"
  template = "${var.tectonic_cluster_name}-etcd-${count.index}.${var.tectonic_base_domain}"
}

module "container_linux" {
  source = "../../modules/container_linux"

  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

module "vpc" {
  source = "../../modules/aws/vpc"

  base_domain              = "${var.tectonic_base_domain}"
  cidr_block               = "${var.tectonic_aws_vpc_cidr_block}"
  cluster_id               = "${local.cluster_id}"
  cluster_name             = "${var.tectonic_cluster_name}"
  custom_dns_name          = "${var.tectonic_dns_name}"
  enable_etcd_sg           = "${length(compact(var.tectonic_etcd_servers)) == 0 ? 1 : 0}"
  external_master_subnets  = "${compact(var.tectonic_aws_external_master_subnet_ids)}"
  external_vpc_id          = "${var.tectonic_aws_external_vpc_id}"
  external_worker_subnets  = "${compact(var.tectonic_aws_external_worker_subnet_ids)}"
  extra_tags               = "${var.tectonic_aws_extra_tags}"
  private_master_endpoints = "${var.tectonic_aws_private_endpoints}"
  public_master_endpoints  = "${var.tectonic_aws_public_endpoints}"

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

  base_domain                = "${var.tectonic_base_domain}"
  cluster_id                 = "${local.cluster_id}"
  cluster_name               = "${var.tectonic_cluster_name}"
  container_image            = "${var.tectonic_container_images["etcd"]}"
  container_linux_channel    = "${var.tectonic_container_linux_channel}"
  container_linux_version    = "${module.container_linux.version}"
  ec2_type                   = "${var.tectonic_aws_etcd_ec2_type}"
  external_endpoints         = "${compact(var.tectonic_etcd_servers)}"
  extra_tags                 = "${var.tectonic_aws_extra_tags}"
  ign_etcd_crt_id_list       = "${module.ignition_masters.etcd_crt_id_list}"
  ign_etcd_dropin_id_list    = "${module.ignition_masters.etcd_dropin_id_list}"
  ign_profile_env_id         = "${local.tectonic_http_proxy_enabled ? module.ignition_masters.profile_env_id : ""}"
  ign_systemd_default_env_id = "${local.tectonic_http_proxy_enabled ? module.ignition_masters.systemd_default_env_id : ""}"
  instance_count             = "${length(data.template_file.etcd_hostname_list.*.id)}"
  root_volume_iops           = "${var.tectonic_aws_etcd_root_volume_iops}"
  root_volume_size           = "${var.tectonic_aws_etcd_root_volume_size}"
  root_volume_type           = "${var.tectonic_aws_etcd_root_volume_type}"
  s3_bucket                  = "${local.s3_bucket}"
  sg_ids                     = "${concat(var.tectonic_aws_etcd_extra_sg_ids, list(module.vpc.etcd_sg_id))}"
  ssh_key                    = "${var.tectonic_aws_ssh_key}"
  subnets                    = "${module.vpc.worker_subnet_ids}"
  etcd_iam_role              = "${var.tectonic_aws_etcd_iam_role_name}"
  ec2_ami                    = "${var.tectonic_aws_ec2_ami_override}"
}

module "ignition_masters" {
  source = "../../modules/ignition"

  assets_location           = "${local.tectonic_bucket}/${local.tectonic_key}"
  base_domain               = "${var.tectonic_base_domain}"
  bootstrap_upgrade_cl      = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider            = "aws"
  cluster_name              = "${var.tectonic_cluster_name}"
  container_images          = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list   = "${var.tectonic_custom_ca_pem_list}"
  etcd_advertise_name_list  = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_ca_cert_pem          = "${local.etcd_ca_crt_pem}"
  etcd_client_crt_pem       = "${local.etcd_client_crt_pem}"
  etcd_client_key_pem       = "${local.etcd_client_key_pem}"
  etcd_count                = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_initial_cluster_list = "${data.template_file.etcd_hostname_list.*.rendered}"
  etcd_peer_crt_pem         = "${local.etcd_peer_crt_pem}"
  etcd_peer_key_pem         = "${local.etcd_peer_key_pem}"
  etcd_server_crt_pem       = "${local.etcd_server_crt_pem}"
  etcd_server_key_pem       = "${local.etcd_server_key_pem}"
  http_proxy                = "${var.tectonic_http_proxy_address}"
  https_proxy               = "${var.tectonic_https_proxy_address}"
  image_re                  = "${var.tectonic_image_re}"
  ingress_ca_cert_pem       = "${local.ingress_certs_ca_cert_pem}"
  iscsi_enabled             = "${var.tectonic_iscsi_enabled}"
  kube_ca_cert_pem          = "${local.kube_certs_ca_cert_pem}"
  kube_dns_service_ip       = "${local.kube_dns_service_ip}"
  kubeconfig_fetch_cmd      = "/opt/s3-puller.sh ${local.kubeconfig_bucket}/${local.kubeconfig_key} /etc/kubernetes/kubeconfig"
  kubelet_debug_config      = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label        = "node-role.kubernetes.io/master"
  kubelet_node_taints       = "node-role.kubernetes.io/master=:NoSchedule"
  no_proxy                  = "${var.tectonic_no_proxy}"
}

module "masters" {
  source = "../../modules/aws/master-asg"

  assets_s3_location                   = "${local.tectonic_bucket}/${local.tectonic_key}"
  autoscaling_group_extra_tags         = "${var.tectonic_autoscaling_group_extra_tags}"
  aws_lbs                              = "${module.vpc.aws_lbs}"
  base_domain                          = "${var.tectonic_base_domain}"
  cluster_id                           = "${local.cluster_id}"
  cluster_name                         = "${var.tectonic_cluster_name}"
  container_images                     = "${var.tectonic_container_images}"
  container_linux_channel              = "${var.tectonic_container_linux_channel}"
  container_linux_version              = "${module.container_linux.version}"
  ec2_type                             = "${var.tectonic_aws_master_ec2_type}"
  extra_tags                           = "${var.tectonic_aws_extra_tags}"
  ign_bootkube_path_unit_id            = "${data.ignition_systemd_unit.bootkube_path_unit.id}"
  ign_bootkube_service_id              = "${data.ignition_systemd_unit.bootkube_service.id}"
  ign_ca_cert_id_list                  = "${module.ignition_masters.ca_cert_id_list}"
  ign_docker_dropin_id                 = "${module.ignition_masters.docker_dropin_id}"
  ign_init_assets_service_id           = "${module.ignition_masters.init_assets_service_id}"
  ign_installer_kubelet_env_id         = "${module.ignition_masters.installer_kubelet_env_id}"
  ign_installer_runtime_mappings_id    = "${module.ignition_masters.installer_runtime_mappings_id}"
  ign_iscsi_service_id                 = "${module.ignition_masters.iscsi_service_id}"
  ign_k8s_node_bootstrap_service_id    = "${module.ignition_masters.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id               = "${module.ignition_masters.kubelet_service_id}"
  ign_locksmithd_service_id            = "${module.ignition_masters.locksmithd_service_id}"
  ign_max_user_watches_id              = "${module.ignition_masters.max_user_watches_id}"
  ign_profile_env_id                   = "${local.tectonic_http_proxy_enabled ? module.ignition_masters.profile_env_id : ""}"
  ign_rm_assets_path_unit_id           = "${module.ignition_masters.rm_assets_path_unit_id}"
  ign_rm_assets_service_id             = "${module.ignition_masters.rm_assets_service_id}"
  ign_s3_puller_id                     = "${module.ignition_masters.s3_puller_id}"
  ign_systemd_default_env_id           = "${local.tectonic_http_proxy_enabled ? module.ignition_masters.systemd_default_env_id : ""}"
  ign_tectonic_path_unit_id            = "${data.ignition_systemd_unit.tectonic_path_unit.id}"
  ign_tectonic_service_id              = "${data.ignition_systemd_unit.tectonic_service.id}"
  ign_update_ca_certificates_dropin_id = "${module.ignition_masters.update_ca_certificates_dropin_id}"
  instance_count                       = "1"
  master_iam_role                      = "${var.tectonic_aws_master_iam_role_name}"
  master_sg_ids                        = "${concat(var.tectonic_aws_master_extra_sg_ids, list(module.vpc.master_sg_id))}"
  private_endpoints                    = "${var.tectonic_aws_private_endpoints}"
  public_endpoints                     = "${var.tectonic_aws_public_endpoints}"
  root_volume_iops                     = "${var.tectonic_aws_master_root_volume_iops}"
  root_volume_size                     = "${var.tectonic_aws_master_root_volume_size}"
  root_volume_type                     = "${var.tectonic_aws_master_root_volume_type}"
  s3_bucket                            = "${local.tectonic_bucket}"
  ssh_key                              = "${var.tectonic_aws_ssh_key}"
  subnet_ids                           = "${module.vpc.master_subnet_ids}"
  ec2_ami                              = "${var.tectonic_aws_ec2_ami_override}"
}

module "ignition_workers" {
  source = "../../modules/ignition"

  bootstrap_upgrade_cl    = "${var.tectonic_bootstrap_upgrade_cl}"
  cloud_provider          = "aws"
  container_images        = "${var.tectonic_container_images}"
  custom_ca_cert_pem_list = "${var.tectonic_custom_ca_pem_list}"
  etcd_ca_cert_pem        = "${local.etcd_ca_crt_pem}"
  http_proxy              = "${var.tectonic_http_proxy_address}"
  https_proxy             = "${var.tectonic_https_proxy_address}"
  image_re                = "${var.tectonic_image_re}"
  ingress_ca_cert_pem     = "${local.ingress_certs_ca_cert_pem}"
  iscsi_enabled           = "${var.tectonic_iscsi_enabled}"
  kube_ca_cert_pem        = "${local.kube_certs_ca_cert_pem}"
  kube_dns_service_ip     = "${local.kube_dns_service_ip}"
  kubeconfig_fetch_cmd    = "/opt/s3-puller.sh ${local.kubeconfig_bucket}/${local.kubeconfig_key} /etc/kubernetes/kubeconfig"
  kubelet_debug_config    = "${var.tectonic_kubelet_debug_config}"
  kubelet_node_label      = "node-role.kubernetes.io/node"
  kubelet_node_taints     = ""
  no_proxy                = "${var.tectonic_no_proxy}"
}

module "workers" {
  source = "../../modules/aws/worker-asg"

  autoscaling_group_extra_tags         = "${var.tectonic_autoscaling_group_extra_tags}"
  cluster_id                           = "${local.cluster_id}"
  cluster_name                         = "${var.tectonic_cluster_name}"
  container_linux_channel              = "${var.tectonic_container_linux_channel}"
  container_linux_version              = "${module.container_linux.version}"
  ec2_type                             = "${var.tectonic_aws_worker_ec2_type}"
  extra_tags                           = "${var.tectonic_aws_extra_tags}"
  ign_ca_cert_id_list                  = "${module.ignition_masters.ca_cert_id_list}"
  ign_docker_dropin_id                 = "${module.ignition_workers.docker_dropin_id}"
  ign_installer_kubelet_env_id         = "${module.ignition_workers.installer_kubelet_env_id}"
  ign_installer_runtime_mappings_id    = "${module.ignition_workers.installer_runtime_mappings_id}"
  ign_iscsi_service_id                 = "${module.ignition_workers.iscsi_service_id}"
  ign_k8s_node_bootstrap_service_id    = "${module.ignition_workers.k8s_node_bootstrap_service_id}"
  ign_kubelet_service_id               = "${module.ignition_workers.kubelet_service_id}"
  ign_locksmithd_service_id            = "${module.ignition_workers.locksmithd_service_id}"
  ign_max_user_watches_id              = "${module.ignition_workers.max_user_watches_id}"
  ign_profile_env_id                   = "${local.tectonic_http_proxy_enabled ? module.ignition_workers.profile_env_id : ""}"
  ign_s3_puller_id                     = "${module.ignition_workers.s3_puller_id}"
  ign_systemd_default_env_id           = "${local.tectonic_http_proxy_enabled ? module.ignition_workers.systemd_default_env_id : ""}"
  ign_update_ca_certificates_dropin_id = "${module.ignition_workers.update_ca_certificates_dropin_id}"
  instance_count                       = "0"
  load_balancers                       = "${var.tectonic_aws_worker_load_balancers}"
  root_volume_iops                     = "${var.tectonic_aws_worker_root_volume_iops}"
  root_volume_size                     = "${var.tectonic_aws_worker_root_volume_size}"
  root_volume_type                     = "${var.tectonic_aws_worker_root_volume_type}"
  s3_bucket                            = "${local.s3_bucket}"
  sg_ids                               = "${concat(var.tectonic_aws_worker_extra_sg_ids, list(module.vpc.worker_sg_id))}"
  ssh_key                              = "${var.tectonic_aws_ssh_key}"
  subnet_ids                           = "${module.vpc.worker_subnet_ids}"
  vpc_id                               = "${module.vpc.vpc_id}"
  worker_iam_role                      = "${var.tectonic_aws_worker_iam_role_name}"
  ec2_ami                              = "${var.tectonic_aws_ec2_ami_override}"
}

module "dns" {
  source = "../../modules/dns/route53"

  api_external_elb_dns_name      = "${module.vpc.aws_api_external_dns_name}"
  api_external_elb_zone_id       = "${module.vpc.aws_elb_api_external_zone_id}"
  api_internal_elb_dns_name      = "${module.vpc.aws_api_internal_dns_name}"
  api_internal_elb_zone_id       = "${module.vpc.aws_elb_api_internal_zone_id}"
  api_ip_addresses               = "${module.vpc.aws_lbs}"
  base_domain                    = "${var.tectonic_base_domain}"
  cluster_id                     = "${local.cluster_id}"
  cluster_name                   = "${var.tectonic_cluster_name}"
  console_elb_dns_name           = "${module.vpc.aws_console_dns_name}"
  console_elb_zone_id            = "${module.vpc.aws_elb_console_zone_id}"
  custom_dns_name                = "${var.tectonic_dns_name}"
  elb_alias_enabled              = true
  etcd_count                     = "${length(data.template_file.etcd_hostname_list.*.id)}"
  etcd_ip_addresses              = "${module.etcd.ip_addresses}"
  external_endpoints             = ["${compact(var.tectonic_etcd_servers)}"]
  master_count                   = "${var.tectonic_master_count}"
  tectonic_external_private_zone = "${var.tectonic_aws_external_private_zone}"
  tectonic_external_vpc_id       = "${module.vpc.vpc_id}"
  tectonic_extra_tags            = "${var.tectonic_aws_extra_tags}"
  tectonic_private_endpoints     = "${var.tectonic_aws_private_endpoints}"
  tectonic_public_endpoints      = "${var.tectonic_aws_public_endpoints}"
}

data "ignition_systemd_unit" "tectonic_service" {
  name    = "tectonic.service"
  enabled = false
  content = "${local.tectonic_service}"
}

data "ignition_systemd_unit" "tectonic_path_unit" {
  name    = "tectonic.path"
  enabled = true
  content = "${local.tectonic_path_unit}"
}

data "ignition_systemd_unit" "bootkube_service" {
  name    = "bootkube.service"
  enabled = false
  content = "${local.bootkube_service}"
}

data "ignition_systemd_unit" "bootkube_path_unit" {
  name    = "bootkube.path"
  enabled = true
  content = "${local.bootkube_path_unit}"
}
