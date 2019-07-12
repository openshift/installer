variable "base_image" {
  type = string
}

variable "bootstrap_ip" {
  type = string
}

variable "cluster_id" {
  type        = string
  description = "The identifier for the cluster."
}

variable "cluster_domain" {
  type        = string
  description = "The domain name of the cluster. All DNS records must be under this domain."
}

variable "flavor_name" {
  type = string
}

variable "instance_count" {
  type = string
}

variable "lb_floating_ip" {
  type = string
}

variable "master_ips" {
  type = list(string)
}

variable "master_sg_ids" {
  type        = list(string)
  default     = ["default"]
  description = "The security group IDs to be applied to the master nodes."
}

variable "master_port_ids" {
  type        = list(string)
  description = "List of port ids for the master nodes"
}

variable "master_port_names" {
  type = list(string)
}

variable "user_data_ign" {
  type = string
}

variable "api_int_ip" {
  type = string
}

variable "node_dns_ip" {
  type = string
}

variable "service_vm_fixed_ip" {
  type = string
}
