output "master_ecs_ids" {
  value = alicloud_instance.master.*.id
}

output "master_ecs_private_ips" {
  value = { for ecs in data.alicloud_instances.master_data.instances : ecs.name => ecs.private_ip }
}