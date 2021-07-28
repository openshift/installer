################################################################
# Configure the IBM Cloud provider
################################################################
variable "powervs_api_key" {
  type        = string
  description = "IBM Cloud API key associated with user's identity"
  default     = "<key>"
}

variable "powervs_vpc_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "eu-gb"
  ##default     = "eu-gb"
}

variable "powervs_vpc_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = "eu-gb-2"
}

variable "powervs_region" {
  type        = string
  description = "The IBM Cloud region where you want to create the resources"
  default     = "lon"
}

variable "powervs_zone" {
  type        = string
  description = "The zone of an IBM Cloud region where you want to create Power System resources"
  default     = "lon04"
}

variable "powervs_resource_group" {
  type        = string
  description = "The cloud instance resource group"
  default     = ""
}

variable "powervs_cloud_instance_id" {
  type        = string
  description = "The cloud instance ID of your account"
  ## TODO: erase default and set via install-config
  default = "e449d86e-c3a0-4c07-959e-8557fdf55482"
}

################################################################
# Configure storage
################################################################
variable "powervs_cos_instance_location" {
  type        = string
  description = "The location of your COS instance"
  default     = "global"
}

variable "powervs_cos_bucket_location" {
  type        = string
  description = "The location to create your COS bucket"
  default     = "us-east"
}

variable "powervs_cos_storage_class" {
  type        = string
  description = "The plan used for your COS instance"
  default     = "smart"
}

################################################################
# Configure instances
################################################################
variable "powervs_image_name" {
  type        = string
  description = "Name of the image used by all nodes in the cluster."
}

variable "powervs_network_name" {
  type        = string
  description = "Name of the network used by the all nodes in the cluster."
  default     = "pvs-ipi-net"
}

variable "powervs_bootstrap_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by the bootstrap node."
  default     = "32"
}

variable "powervs_bootstrap_processors" {
  type        = string
  description = "Number of processors used by the bootstrap node."
  default     = "0.5"
}

variable "powervs_master_memory" {
  type        = string
  description = "Amount of memory, in  GiB, used by each master node."
  default     = "32"
}

variable "powervs_master_processors" {
  type        = string
  description = "Number of processors used by each master node."
  default     = "0.5"
}

variable "powervs_proc_type" {
  type        = string
  description = "The type of processor mode for all nodes (shared/dedicated)"
  default     = "shared"
}

variable "powervs_sys_type" {
  type        = string
  description = "The type of system (s922/e980)"
  default     = "s922"
}

variable "powervs_base_domain" {
  type        = string
  description = "The base domain name of the cluster"
  default     = "openshift-on-power.com"
}

variable "powervs_cluster_domain" {
  type        = string
  description = "The name of the cluster that all DNS records must belong to."
  default     = "rdr-ipi"
}

variable "powervs_ssh_key" {
  type        = string
  description = "Public key for keypair used to access cluster."
  default     = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDrgIZ+xyn6Hy0DD1UDVgeZjxIUvsdjVa8AyM0gfQlWZAUW7IyAfwpxyZ+To1h90ltINqjpkyiOdMLYkXvB40LCDlq9jR9B2X7cjAZD9ZGJLrWlTgqnSrTtKK5WIPkC5TYLczUGin1BuAxFUb2VAX83omSlVrObPK90JyqCgobh+j3uAZJXrs+5MEcJieobIbxdeLwujsRlC0vzF4fjngRgnWUNyVx04jztyWgAfU3ZrmgxO4+/2puHaPpouxbgUxDXdr+JtwXJ3/zeAO0Zjs1L9xawzbYua+oQD2o7OjM3uJ02wVAcMr/FX7nr4yvxyYOBglXWDEdL8OZtlGCmH1C3aDejsR1GeJikrHg+GhrZ+afRbZTMTlmeZeGOvDRCBgR8ZSqLLDcPOl/y1HFBL9/pQeBFRvEMz8NGncazdSvBHbeFT0XiyYTxXfJUi5cAFQn52tVohzAI4L5gS2WgrNo4jw4YhXcfSoxuhEwbZZtdbht0iQa83zev5+accrFYsW0= bpradipt@Pradiptas-MBP"
}

## TODO: Set this in install-config instead
variable "powervs_vpc_name" {
  type        = string
  description = "Name of the IBM Cloud Virtual Private Cloud (VPC) to setup the load balancer."
  default     = "powervs-ipi"
}

variable "powervs_vpc_subnet_name" {
  type        = string
  description = "Name of the VPC subnet having DirectLink access to the PowerVS private network"
  default     = "subnet2"
}
