variable "project_id" {
  type        = string
  description = "The target GCP project for the cluster."
}

variable "cluster_id" {
  type        = string
  description = "The name of the cluster."
}

variable "ignition" {
  type        = string
  description = "The content of the masters ignition file."
}

variable "image" {
  type        = string
  description = "The image for the master instances."
}

variable "instance_count" {
  type        = string
  description = "The number of master instances to launch."
}

variable "gcp_extra_labels" {
  type        = map(string)
  description = "GCP labels to be applied to created resources."
  default     = {}
}

variable "machine_type" {
  type        = string
  description = "The machine type for the master instances."
}

variable "service_account" {
  type        = string
  description = "The service account used by the instances."
}
variable "subnet" {
  type        = string
  description = "The subnetwork the master instances will be added to."
}

variable "tags" {
  type        = list(string)
  description = "The list of network tags which will be added to the control plane instances."
}

variable "root_volume_size" {
  type        = string
  description = "The size of the volume in gigabytes for the root block device."
}

variable "root_volume_type" {
  type        = string
  description = "The type of volume for the root block device."
}

variable "root_volume_kms_key_link" {
  type        = string
  description = "The GCP self link of KMS key to encrypt the volume."
  default     = null
}

variable "zones" {
  type = list
}

variable "secure_boot" {
  type        = string
  description = "Verify the digital signature of all boot components."
  default     = ""
}

variable "confidential_compute" {
  type        = string
  description = "Defines whether the instance should have confidential compute enabled."
  default     = ""
}

variable "on_host_maintenance" {
  type        = string
  description = "The behavior when a maintenance event occurs."
  default     = ""
}

variable "gcp_extra_tags" {
  type        = map(string)
  description = <<EOF
(optional) Extra GCP tags to be applied to the created resources.
Example: `{ "tagKeys/123" = "tagValues/456", "tagKeys/456" = "tagValues/789" }`
EOF
  default = {}
}
