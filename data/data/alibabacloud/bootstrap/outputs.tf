output "bootstrap_public_ip" {
  value = data.alicloud_instances.bootstrap_data.instances.0.public_ip
}
