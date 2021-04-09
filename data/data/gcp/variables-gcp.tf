variable "gcp_project_id" {
  type        = string
  description = "The target GCP project for the cluster."
}

variable "gcp_service_account" {
  type        = string
  description = "The service account for authenticating with GCP APIs."
}

variable "gcp_region" {
  type        = string
  description = "The target GCP region for the cluster."
}

variable "gcp_extra_labels" {
  type = map(string)

  description = <<EOF
(optional) Extra GCP labels to be applied to created resources.
Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}

variable "gcp_bootstrap_instance_type" {
  type = string
  description = "Instance type for the bootstrap node. Example: `n1-standard-4`"
}

variable "gcp_master_instance_type" {
  type = string
  description = "Instance type for the master node(s). Example: `n1-standard-4`"
}

variable "gcp_image_uri" {
  type = string
  description = "URL to Raw Image for all nodes. This is used in case a new image needs to be generated for the nodes."
}

variable "gcp_image" {
  type = string
  description = "URL to the Image for all nodes."
}

variable "gcp_preexisting_image" {
  type = bool
  default = true
  description = "Specifies whether an existing GCP Image should be used or a new one created for installation"
}

variable "gcp_master_root_volume_type" {
  type = string
  description = "The type of volume for the root block device of master nodes."
}

variable "gcp_master_root_volume_size" {
  type = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "gcp_public_dns_zone_name" {
  type = string
  default = null
  description = "The name of the public DNS zone to use for this cluster"
}

variable "gcp_master_availability_zones" {
  type = list(string)
  description = "The availability zones in which to create the masters. The length of this list must match master_count."
}

variable "gcp_preexisting_network" {
  type = bool
  default = false
  description = "Specifies whether an existing network should be used or a new one created for installation."
}

variable "gcp_cluster_network" {
  type = string
  description = "The name of the cluster network, either existing or to be created."
}

variable "gcp_control_plane_subnet" {
  type = string
  description = "The name of the subnet for the control plane, either existing or to be created."
}

variable "gcp_compute_subnet" {
  type = string
  description = "The name of the subnet for worker nodes, either existing or to be created"
}

variable "gcp_publish_strategy" {
  type = string
  description = "The cluster publishing strategy, either Internal or External"
}

variable "gcp_image_licenses" {
  type = list(string)
  description = "The licenses to use when creating compute instances"
  default = []
}

variable "gcp_root_volume_kms_key_link" {
  type = string
  description = "The GCP self link of KMS key to encrypt the volume."
  default = null
}
