variable "api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
}

variable "memory" {
  type        = string
  description = "The amount of memory to assign to each node in GB."
}

variable "processors" {
  type        = string
  description = "The processor count for each node."
}

variable "ignition" {
  type        = string
  description = "The contents of the ignition file."
}

variable "cloud_instance_id" {
  type        = string
  description = "The Power VS Service Instance (aka Cloud Instance) ID."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "image_id" {
  type        = string
  description = "The ID of the Power VS boot image for the nodes."
}

variable "proc_type" {
  type        = string
  description = "The type of processor to be assigned (e.g. capped, dedicated, shared) to the nodes."
}

variable "sys_type" {
  type        = string
  description = "The type of system on which to provision the nodes (e.g s922)."
}

variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "ssh_key_name" {
  type        = string
  description = "The SSH Key name."
}

variable "cos_instance_location" {
  type        = string
  description = "Specifies whether the Cloud Object Store instance is global or in a specific region. Used for the ignition file."
}

variable "cos_bucket_location" {
  type        = string
  description = "The region in which to create the Cloud Object Store bucket. Used for the igntion file."
}

variable "cos_storage_class" {
  type        = string
  description = "The storage class for the Cloud Object Store instance."
}

variable "vpc_region" {
  type        = string
  description = "The IBM Cloud region in which the VPC is created."
}

variable "vpc_zone" {
  type        = string
  description = "The IBM Cloud zone in which the VPC is created."
}

variable "vpc_id" {
  type        = string
  description = "The ID of the VPC created for the cluster load balancers."
}

variable "powervs_region" {
  type        = string
  description = "The Power VS region in which to create resources."
}

variable "powervs_zone" {
  type        = string
  description = "The Power VS zone in which to create resources."

}

variable "dhcp_id" {
  type        = string
  description = "The ID of the Power VS DHCP Service."
}

variable "dhcp_network_id" {
  type        = string
  description = "The ID of the Power VS DHCP network."
}

variable "lb_ext_id" {
  type        = string
  description = "The ID of the external load balancer in the IBM Cloud VPC"
}

variable "lb_int_id" {
  type        = string
  description = "The ID of the private load balancer in the IBM Cloud VPC"
}

variable "machine_cfg_pool_id" {
  type        = string
  description = "The ID of the load balancer pool for the machine-config server."
}

variable "api_pool_int_id" {
  type        = string
  description = "The ID of the private load balancer pool for the API server."
}

variable "api_pool_ext_id" {
  type        = string
  description = "The ID of the public load balancer pool for the API server."
}
