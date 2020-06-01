output "releaseimage_template_id" {
  value = data.ovirt_templates.finalTemplate.templates.0.id
}

output "tmp_import_vm" {
  value = data.tmp_import_vm.id
}
