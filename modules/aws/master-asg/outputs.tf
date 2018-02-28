output "aws_launch_configuration" {
  value = "${aws_launch_configuration.master_conf.id}"
}

output "subnet_ids" {
  value = "${var.subnet_ids}"
}

output "aws_lbs" {
  value = "${var.aws_lbs}"
}

output "cluster_id" {
  value = "${var.cluster_id}"
}
