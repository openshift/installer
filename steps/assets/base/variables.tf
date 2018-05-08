# We use the same variables as the platform-specific step, to keep from going insane. Here
# is where we can define variables that the steps can pass directly
variable "etcd_count" {
  type = "string"
}

variable "cloud_provider" {
  type = "string"
}

variable "ingress_kind" {
  type = "string"
}
