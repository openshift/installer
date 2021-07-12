variable "server_group_names" {
  type        = set(string)
  description = "Names of the server groups for the worker nodes."
}

variable "server_group_policy" {
  type        = string
  description = "Policy of the server groups for the worker nodes."
}
