//////
// Nutanix variables
//////

variable "nutanix_prism_central_address" {
  type        = string
  description = "Address to connect to Prism Central."
}

variable "nutanix_prism_central_port" {
  type        = string
  description = "Port to connect to Prism Central."
}

variable "nutanix_username" {
  type        = string
  description = "Prism Central user for the environment."
}

variable "nutanix_password" {
  type        = string
  description = "Prism Central user password"
}

variable "nutanix_prism_element_uuid" {
  type        = string
  description = "This is the uuid of the Prism Element cluster."
}

variable "nutanix_image_filepath" {
  type        = string
  description = "This is the filepath to the image file that will be imported into Prism Central."
}

variable "nutanix_image" {
  type        = string
  description = "This is the name to the image that will be imported into Prism Central."
}

variable "nutanix_subnet_uuid" {
  type        = string
  description = "This is the uuid of the publicly accessible subnet for cluster ingress and access."
}

variable "nutanix_bootstrap_ignition_image" {
  type        = string
  description = "Name of the image containing the bootstrap ignition files"
}

variable "nutanix_bootstrap_ignition_image_filepath" {
  type        = string
  description = "Path to the image containing the bootstrap ignition files"
}

///////////
// Control Plane machine variables
///////////

variable "nutanix_control_plane_memory_mib" {
  type = number
}

variable "nutanix_control_plane_disk_mib" {
  type = number
}

variable "nutanix_control_plane_num_cpus" {
  type = number
}

variable "nutanix_control_plane_cores_per_socket" {
  type = number
}
