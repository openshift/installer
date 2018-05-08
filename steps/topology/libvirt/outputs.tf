output "libvirt_network_id" {
  value = "${libvirt_network.tectonic_net.id}"
}

output "libvirt_base_volume_id" {
  value = "${module.libvirt_base_volume.coreos_base_volume_id}"
}
