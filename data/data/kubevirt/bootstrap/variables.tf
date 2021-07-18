variable "cluster_id" {
  description = "The ID of OpenShift cluster"
}

variable "namespace" {
  type        = string
  description = "The namespace/project in the infra cluster, in which all the tenant cluster resources should be created"
}

variable "storage" {
  type        = string
  description = "bootstrap VM disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "35Gi"
}

variable "memory" {
  type        = string
  description = "bootstrap VM memory size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
  default     = "8G"
}

variable "cpu" {
  type        = string
  description = "bootstrap VM number of cores"
  default     = "4"
}

variable "ignition_data" {
  type        = string
  description = "Ignition config file contents of the bootstrap VM"
}

variable "storage_class" {
  type        = string
  description = "The \"class\" of the storage located in the infra cluster"
}

variable "network_name" {
  type        = string
  description = "The name of the sub network created in the infracluster which should be used by the tenant cluster resources"
}

variable "interface_binding_method" {
  type        = string
  description = "The interface binding method of the nodes of the tenantcluster"
}

variable "pv_access_mode" {
  type        = string
  description = "The access mode which all the persistent volumes should be created with [ReadWriteOnce,ReadOnlyMany,ReadWriteMany]"
}

variable "pvc_name" {
  type        = string
  description = "The Persistent data volume which bootstrap VM should be cloned from"
}

variable "labels" {
  type = map(string)

  description = <<EOF
(optional) Labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
