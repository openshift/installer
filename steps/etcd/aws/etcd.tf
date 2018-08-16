provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"

  assume_role {
    role_arn     = "${var.tectonic_aws_installer_role == "" ? "" : "${var.tectonic_aws_installer_role}"}"
    session_name = "TECTONIC_INSTALLER_${var.tectonic_cluster_name}"
  }
}

module "defaults" {
  source = "../../../modules/aws/target-defaults"

  region     = "${var.tectonic_aws_region}"
  profile    = "${var.tectonic_aws_profile}"
  role_arn   = "${var.tectonic_aws_installer_role}"
  etcd_count = "${var.tectonic_etcd_count}"
}

module "container_linux" {
  source = "../../../modules/container_linux"

  release_channel = "${var.tectonic_container_linux_channel}"
  release_version = "${var.tectonic_container_linux_version}"
}

resource "aws_s3_bucket_object" "ignition_etcd" {
  count   = "${module.defaults.etcd_count}"
  bucket  = "${local.s3_bucket}"
  key     = "ignition_etcd_${count.index}.json"
  content = "${local.ignition[count.index]}"
  acl     = "private"

  server_side_encryption = "AES256"

  tags = "${merge(map(
      "Name", "${var.tectonic_cluster_name}-ignition-etcd-${count.index}",
      "KubernetesCluster", "${var.tectonic_cluster_name}",
      "tectonicClusterID", "${var.tectonic_cluster_id}"
    ), var.tectonic_aws_extra_tags)}"
}

module "etcd" {
  source = "../../../modules/aws/etcd"

  base_domain             = "${var.tectonic_base_domain}"
  cluster_id              = "${var.tectonic_cluster_id}"
  cluster_name            = "${var.tectonic_cluster_name}"
  container_image         = "${var.tectonic_container_images["etcd"]}"
  container_linux_channel = "${var.tectonic_container_linux_channel}"
  container_linux_version = "${module.container_linux.version}"
  ec2_type                = "${var.tectonic_aws_etcd_ec2_type}"
  extra_tags              = "${var.tectonic_aws_extra_tags}"
  instance_count          = "${module.defaults.etcd_count}"
  region                  = "${var.tectonic_aws_region}"
  root_volume_iops        = "${var.tectonic_aws_etcd_root_volume_iops}"
  root_volume_size        = "${var.tectonic_aws_etcd_root_volume_size}"
  root_volume_type        = "${var.tectonic_aws_etcd_root_volume_type}"
  s3_bucket               = "${local.s3_bucket}"
  sg_ids                  = "${concat(var.tectonic_aws_etcd_extra_sg_ids, list(local.sg_id))}"
  subnets                 = ["${local.subnet_ids_workers}"]
  etcd_iam_role           = "${var.tectonic_aws_etcd_iam_role_name}"
  ec2_ami                 = "${var.tectonic_aws_ec2_ami_override}"
}
