output "bucket_id" {
  value = alicloud_oss_bucket.bucket.id
}

output "bootstrap_role_name" {
  value = alicloud_ram_role.role.name
}

output "bootstrap_role_arn" {
  value = alicloud_ram_role.role.arn
}

output "bootstrap_role_policy_id" {
  value = alicloud_ram_policy.role_policy.id
}

output "bootstrap_sg_id" {
  value = alicloud_security_group.sg_bootstrap.id
}

output "bootstrap_ecs_id" {
  value = alicloud_instance.bootstrap.id
}
