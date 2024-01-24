output "bootstrap_ip" {
  value = local.is_external ? data.alicloud_instances.bootstrap_data.instances.0.public_ip : data.alicloud_instances.bootstrap_data.instances.0.private_ip
}

output "bootstrap_ecs_id" {
  value = alicloud_instance.bootstrap.id
}
