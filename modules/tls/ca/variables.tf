variable "root_ca_key_alg" {
  description = "Algorithm used to generate root_ca_key (required if root_ca_cert is specified)"
  type        = "string"
  default     = "RSA"
}

variable "kube_ca_key_alg" {
  description = "Algorithm used to generate kube_ca_key (required if root_ca_cert is specified)"
  type        = "string"
  default     = "RSA"
}

variable "aggregator_ca_key_alg" {
  description = "Algorithm used to generate aggregator_ca_key (required if root_ca_cert is specified)"
  type        = "string"
  default     = "RSA"
}

variable "service_serving_ca_key_alg" {
  description = "Algorithm used to generate service_serving_ca_key (required if root_ca_cert is specified)"
  type        = "string"
  default     = "RSA"
}

variable "etcd_ca_key_alg" {
  description = "Algorithm used to generate etcd_ca_key (required if root_ca_cert is specified)"
  type        = "string"
  default     = "RSA"
}

variable "root_ca_cert_pem_path" {
  type    = "string"
  default = ""
}

variable "root_ca_key_pem_path" {
  type    = "string"
  default = ""
}

variable "etcd_ca_cert_pem_path" {
  type    = "string"
  default = ""
}

variable "etcd_ca_key_pem_path" {
  type    = "string"
  default = ""
}

variable "kube_ca_cert_pem_path" {
  type    = "string"
  default = ""
}

variable "kube_ca_key_pem_path" {
  type    = "string"
  default = ""
}

variable "aggregator_ca_cert_pem_path" {
  type    = "string"
  default = ""
}

variable "aggregator_ca_key_pem_path" {
  type    = "string"
  default = ""
}

variable "service_serving_ca_cert_pem_path" {
  type    = "string"
  default = ""
}

variable "service_serving_ca_key_pem_path" {
  type    = "string"
  default = ""
}
