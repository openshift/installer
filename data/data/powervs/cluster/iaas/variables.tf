variable "cluster_id" {
  type        = string
  description = "The ID created by the installer to uniquely identify the created cluster."
}

variable "resource_group" {
  type        = string
  description = "The name of the Power VS resource group to which the user belongs."
}

variable "powervs_zone" {
  type        = string
  description = "The Power VS zone in which to create resources."
}

variable "service_instance_name" {
  type        = string
  description = "Optionally, the service instance name of an existing object before cluster creation"
}

variable "wait_for_workspace" {
  type        = string
  description = "Time to wait after creating a Power VS workspace to allow PER to initialize"
}
