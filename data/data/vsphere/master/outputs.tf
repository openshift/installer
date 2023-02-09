output "control_plane_ips" {
  value = [for k, v in vsphere_virtual_machine.vm_master : v.default_ip_address]
}

output "control_plane_moids" {
  value = [for k, v in vsphere_virtual_machine.vm_master : v.moid]
}
