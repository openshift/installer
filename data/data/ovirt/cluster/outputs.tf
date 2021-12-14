output "control_plane_vm_ids" {
  value = module.masters.control_plane_vm_ids
}

output "release_image_template_id" {
  value = module.template.releaseimage_template_id
}
