variable "kubevirt_namespace" {
  type        = string
  description = "The namespace/project in the infracluster which all the tenantcluster resources should be created in"
}

variable "kubevirt_source_pvc_name" {
  type        = string
  description = "The Persistent data volume which all the vms (workers/masters) should be cloned from"
}

variable "kubevirt_master_storage" {
  type        = string
  description = "master VM disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "kubevirt_image_url" {
  type        = string
  description = "The source image URL to be used to create the source persistent data volume (all the VMs are cloned from)"
}

variable "kubevirt_master_memory" {
  type        = string
  description = "master VM memory size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "kubevirt_master_cpu" {
  type        = string
  description = "master VM number of cores"
}

variable "kubevirt_storage_class" {
  type        = string
  description = "The \"class\" of the storage located in the infracluster"
}

variable "kubevirt_network_name" {
  type        = string
  description = "The name of the sub network created in the infracluster which should be used by the tenantcluster resources"
}

variable "kubevirt_interface_binding_method" {
  type        = string
  description = "The interface binding method of the nodes of the tenantcluster"
}

variable "kubevirt_pv_access_mode" {
  type        = string
  description = "The access mode which all the persistent volumes should be created with [ReadWriteOnce,ReadOnlyMany,ReadWriteMany]"
}

variable "kubevirt_labels" {
  type = map(string)

  description = <<EOF
(optional) Labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
