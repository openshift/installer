variable "cluster_id" {
  description = "The ID of OpenShift cluster"
}

variable "master_count" {
  description = "The number of master vm instances"
}

variable "namespace" {
  type        = string
  description = "The namespace/project in the infracluster which all the tenantcluster resources should be created in"
}

variable "storage" {
  type        = string
  description = "master VM disk size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "memory" {
  type        = string
  description = "master VM memory size, of type Quantity (see: https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go)"
}

variable "cpu" {
  type        = string
  description = "master VM number of cores"
}

variable "ignition_data" {
  type        = string
  description = "Ignition config file contents of the master VMs"
}

variable "storage_class" {
  type        = string
  description = "The \"class\" of the storage located in the infracluster"
}

variable "network_name" {
  type        = string
  description = "The name of the sub network created in the infracluster which should be used by the tenantcluster resources"
}

variable "pv_access_mode" {
  type        = string
  description = "The access mode which all the persistant volumes should be created with [ReadWriteOnce,ReadWriteMany]"
}

variable "pvc_name" {
  type        = string
  description = "The Persistant data volume which all the vms (workers/masters) should be cloned from"
}

variable "labels" {
  type = map(string)

  description = <<EOF
(optional) Labels to be applied to created resources.

Example: `{ "key" = "value", "foo" = "bar" }`
EOF

  default = {}
}
