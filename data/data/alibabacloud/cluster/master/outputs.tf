output "master_ecs_ids" {
  value = alicloud_instance.master.*.id
}

output "master_ecs_private_ips" {
  value = data.alicloud_instances.master_data.instances.*.private_ip
}