variable "gcp_project_id" {
  type        = string
  description = "The target GCP project for the cluster."
}

variable "gcp_network_project_id" {
  type        = string
  description = "The project that the network and subnets exist in when they are not in the main ProjectID."
  default     = ""
}

variable "gcp_service_account" {
  type        = string
  description = "The service account for authenticating with GCP APIs."
  default     = null
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

variable "gcp_bootstrap_lb" {
  type = bool
  description = "Setting this to false allows the bootstrap resources to be removed from the cluster load balancers."
  default = true
}

variable "gcp_bootstrap_instance_type" {
  type = string
  description = "Instance type for the bootstrap node. Example: `n1-standard-4`"
}

variable "gcp_master_instance_type" {
  type = string
  description = "Instance type for the master node(s). Example: `n1-standard-4`"
}

variable "gcp_image" {
  type = string
  description = "URL to the Image for all nodes."
}

variable "gcp_instance_service_account" {
  type = string
  description = "The service account used by the instances."
  default = ""
}

variable "gcp_master_root_volume_type" {
  type = string
  description = "The type of volume for the root block device of master nodes."
}

variable "gcp_master_root_volume_size" {
  type = string
  description = "The size of the volume in gigabytes for the root block device of master nodes."
}

variable "gcp_public_zone_name" {
  type = string
  default = null
  description = "The name of the public DNS zone to use for this cluster"
}

variable "gcp_private_zone_name" {
  type = string
  default = ""
  description = "The name of the private DNS zone to use for this cluster, if one already exists"
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

variable "gcp_root_volume_kms_key_link" {
  type = string
  description = "The GCP self link of KMS key to encrypt the volume."
  default = null
}

variable "gcp_control_plane_tags" {
  type = list(string)
  description = "The list of network tags which will be added to the control plane instances."

}

variable "gcp_create_firewall_rules" {
  type = bool
  default = true
  description = "Create the cluster's network firewall rules."
}

variable "gcp_master_secure_boot" {
  type = string
  description = "Verify the digital signature of all boot components."
  default = ""
}

variable "gcp_master_confidential_compute" {
  type = string
  description = "Defines whether the instance should have confidential compute enabled."
  default = ""
}

variable "gcp_master_on_host_maintenance" {
  type = string
  description = "The behavior when a maintenance event occurs."
  default = ""
}

variable "gcp_user_provisioned_dns" {
  type = bool
  default = false
  description = <<EOF
When true the user has selected to configure their own dns solution, and no dns records will be created.
EOF
}

variable "gcp_extra_tags" {
  type        = map(string)
  description = <<EOF
(optional) Extra GCP tags to be applied to the created resources.
Example: `{ "tagKeys/123" = "tagValues/456", "tagKeys/456" = "tagValues/789" }`
EOF
  default = {}
}

variable "gcp_ignition_shim" {
  type = string
  description = "Ignition stub containing the signed url that points to the bucket containing the ignition data."
  default = ""
}

variable "gcp_signed_url" {
  type = string
  description = "Presigned url for bootstrap ignition link to the bucket where the ignition shim is stored."
}