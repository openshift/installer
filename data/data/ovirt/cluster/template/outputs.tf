output "releaseimage_template_id" {
  value = var.tmp_import_vm_id == "" ? data.ovirt_templates.finalTemplate.0.templates.0.id : ovirt_template.releaseimage_template.0.id
}
