output "device_ip" {
  value = packet_device.bootstrap.network[0].address
}

output "device_hostname" {
  value = packet_device.bootstrap.hostname
}

output "device_id" {
  value = packet_device.bootstrap.id
}
