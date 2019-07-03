variable "gcp_bootstrap_machine_type" {
  type        = string
  description = "Machine type for the bootstrap node. Example: `n1-standard-2`."
}

variable "gcp_master_machine_type" {
  type        = string
  description = "Machine type for the master node(s). Example: `n1-standard-2`."
}

variable "gcp_extra_labels" {
  type = map(string)

  description = <<EOF
(optional) Extra GCP labels to be applied to created resources.
Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "gcp_image_name_override" {
  type = string
  description = "(optional) Image name override for all nodes. Example: `image-foobar123`."
  default = ""
}

variable "gcp_project" {
  type = string
  description = "The target GCP project for the cluster."
}

variable "gcp_region" {
  type = string
  description = "The target GCP region for the cluster."
}
