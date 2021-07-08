output "bootstrap_ip" {
  value = libvirt_domain.bootstrap.*.network_interface.0.addresses[0]
}
