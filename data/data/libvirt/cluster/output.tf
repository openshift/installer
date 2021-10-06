output "pool" {
  value = libvirt_pool.storage_pool.name
}

output "base_volume_id" {
  value = libvirt_volume.coreos_base.id
}

output "network_id" {
  value = libvirt_network.net.id
}

output "control_plane_ips" {
  value = var.libvirt_master_ips
}
