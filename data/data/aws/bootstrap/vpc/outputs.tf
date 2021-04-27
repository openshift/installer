output "vpc_id" {
  value = data.aws_vpc.cluster_vpc.id
}

output "az_to_private_subnet_id" {
  value = zipmap(data.aws_subnet.private.*.availability_zone, data.aws_subnet.private.*.id)
}

output "az_to_public_subnet_id" {
  value = zipmap(data.aws_subnet.public.*.availability_zone, data.aws_subnet.public.*.id)
}

output "master_sg_id" {
  value = data.aws_security_group.master.id
}

output "aws_lb_target_group_arns" {
  // The order of the list is very important because the consumers assume the 3rd item is the external aws_lb_target_group
  // Because of the issue https://github.com/hashicorp/terraform/issues/12570, the consumers cannot use a dynamic list for count
  // and therefore are force to implicitly assume that the list is of aws_lb_target_group_arns_length - 1, in case there is no api_external
  value = compact(
    concat(
      data.aws_lb_target_group.api_internal.*.arn,
      data.aws_lb_target_group.services.*.arn,
      data.aws_lb_target_group.api_external.*.arn,
    ),
  )
}

output "aws_lb_target_group_arns_length" {
  // 2 for private endpoints and 1 for public endpoints
  value = "3"
}
