// The flavor ID as given in `openstack flavor list`.
// Specifies the size (CPU/Memory/Drive) of the VM.
variable "tectonic_openstack_flavor_id" {
  type = "string"
}

// The image ID as given in `openstack image list`.
// Specifies the OS image of the VM.
variable "tectonic_openstack_image_id" {
  type = "string"
}
