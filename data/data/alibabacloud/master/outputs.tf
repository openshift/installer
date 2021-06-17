output "master_ecs_ids" {
  value = alicloud_instance.master.*.id
}
