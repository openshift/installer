output "tmp_import_vm_id" {
  value      = length(var.ovirt_base_image_name) == 0 ? ovirt_vm.tmp_import_vm.0.id : ""
  depends_on = [ovirt_nic.tmp_import_vm, ovirt_disk_attachment.tmp_import_vm]
}
