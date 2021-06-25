output "tmp_import_vm_id" {
  value = length(var.ovirt_base_image_name) == 0 ? ovirt_vm.tmp_import_vm.0.id : ""
}