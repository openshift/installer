output "control_plane_ips" {
  value = vsphere_virtual_machine.vm_master.*.default_ip_address
}
