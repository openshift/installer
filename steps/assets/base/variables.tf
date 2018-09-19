# We use the same variables as the platform-specific step, to keep from going insane. Here
# is where we can define variables that the steps can pass directly
variable "cloud_provider" {
  type = "string"
}

variable "ingress_kind" {
  type = "string"
}

variable "aws_worker_ign_config" {
  type    = "string"
  default = ""
}

variable "libvirt_tls_ca_pem" {
  type    = "string"
  default = ""
}

variable "libvirt_tls_cert_pem" {
  type    = "string"
  default = ""
}

variable "libvirt_tls_key_pem" {
  type    = "string"
  default = ""
}
