variable "master_sg_ids" {
  type        = list(string)
  default     = ["default"]
  description = "The security group IDs to be applied to the master nodes."
}

variable "master_port_ids" {
  type        = list(string)
  description = "List of port ids for the master nodes"
}
