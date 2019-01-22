variable "google_config_version" {
  description = <<EOF
(internal) This declares the version of the GCP configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "google_master_instance_type" {
  type        = "string"
  description = "Instance size for the master node(s). Example: `m4.large`."

  # FIXME: get this wired up to the machine default
  default = "n1-standard-2"
}

variable "google_image_name_override" {
  type        = "string"
  description = "(optional) Image name override for all nodes. Example: `image-foobar123`."
  default     = ""
}

variable "google_extra_labels" {
  type = "map"

  description = <<EOF
(optional) Extra GCP labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "google_master_root_volume_type" {
  type        = "string"
  default     = "pd-standard"
  description = "The type of volume for the root block device of master nodes."
}

variable "google_master_root_volume_size" {
  type        = "string"
  default     = "120"
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}


variable "google_region" {
  type        = "string"
  description = "The target GCP region for the cluster."
}
