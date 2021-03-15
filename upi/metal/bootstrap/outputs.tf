output "device_ip" {
  value = metal_device.bootstrap.network[0].address
}

output "device_hostname" {
  value = metal_device.bootstrap.hostname
}

output "device_id" {
  value = metal_device.bootstrap.id
}
