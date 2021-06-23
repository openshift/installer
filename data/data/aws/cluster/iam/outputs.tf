output "bootstrap_instance_profile_name" {
  value = var.include_bootstrap ? aws_iam_instance_profile.bootstrap[0].name : ""
}