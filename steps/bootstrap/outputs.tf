# Masters
output "aws_launch_configuration_masters" {
  value = "${module.masters.aws_launch_configuration}"
}

output "subnet_ids_masters" {
  value = "${module.masters.subnet_ids}"
}

output "aws_lbs_masters" {
  value = "${module.masters.aws_lbs}"
}
