output "control_plane_vm_ids" {
  value = ovirt_vm_start.master.*.id
}
