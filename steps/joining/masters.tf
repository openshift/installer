provider "aws" {
  region  = "${var.tectonic_aws_region}"
  profile = "${var.tectonic_aws_profile}"
  version = "1.8.0"
}

resource "aws_autoscaling_group" "masters" {
  name                 = "${var.tectonic_cluster_name}-masters"
  desired_capacity     = "${var.tectonic_master_count}"
  max_size             = "${var.tectonic_master_count * 3}"
  min_size             = "${var.tectonic_master_count}"
  launch_configuration = "${local.aws_launch_configuration_masters}"
  vpc_zone_identifier  = ["${local.subnet_ids_masters}"]

  load_balancers = ["${local.aws_lbs_masters}"]

  tags = [
    {
      key                 = "Name"
      value               = "${var.tectonic_cluster_name}-master"
      propagate_at_launch = true
    },
    {
      key                 = "kubernetes.io/cluster/${var.tectonic_cluster_name}"
      value               = "owned"
      propagate_at_launch = true
    },
    {
      key                 = "tectonicClusterID"
      value               = "${local.cluster_id}"
      propagate_at_launch = true
    },
    "${var.tectonic_autoscaling_group_extra_tags}",
  ]

  lifecycle {
    create_before_destroy = true
  }
}
