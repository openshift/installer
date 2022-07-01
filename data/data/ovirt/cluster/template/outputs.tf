output "releaseimage_template_id" {
  value = var.tmp_import_vm_id == "" ? one(data.ovirt_templates.finalTemplate.0.templates.*.id) : ovirt_template.releaseimage_template.0.id
}
