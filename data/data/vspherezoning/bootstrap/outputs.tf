output "bootstrap_ip" {
  value = vsphere_virtual_machine.vm_bootstrap.default_ip_address
}

output "bootstrap_moid" {
  value = vsphere_virtual_machine.vm_bootstrap.moid
}
