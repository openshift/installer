output "aws_launch_configuration" {
  value = "${aws_launch_configuration.worker_conf.id}"
}

output "subnet_ids" {
  value = "${var.subnet_ids}"
}

output "aws_lbs" {
  value = "${var.load_balancers}"
}

output "cluster_id" {
  value = "${var.cluster_id}"
}
