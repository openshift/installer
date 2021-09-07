variable "pvc_name" {
  type        = string
  description = "The Persistent data volume name"
}

variable "namespace" {
  type        = string
  description = "The namespace/project in the infracluster which all the tenantcluster resources should be created in"
}

variable "storage" {
  type        = string
  description = "persistent data volume disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "20Gi"
}

variable "pv_access_mode" {
  type        = string
  description = "The access mode which all the persistent volumes should be created with [ReadWriteOnce,ReadOnlyMany,ReadWriteMany]"
}

variable "storage_class" {
  type        = string
  description = "The \"class\" of the storage located in the infracluster"
}

variable "image_url" {
  type        = string
  description = "The source image URL to be used to create the source persistent data volume (all the VMs are cloned from)"
}

variable "labels" {
  type = map(string)

  description = <<EOF
(optional) Labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
