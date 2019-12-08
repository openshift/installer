output "releaseimage_template_id" {
  value = data.ovirt_templates.finalTemplate.templates.0.id
}
