variable "kube_ca_cert_pem" {
  description = "PEM-encoded CA certificate"
  type        = "string"
}

variable "kube_ca_key_alg" {
  description = "Algorithm used to generate kube_ca_key"
  type        = "string"
}

variable "kube_ca_key_pem" {
  description = "PEM-encoded CA key"
  type        = "string"
}

variable "aggregator_ca_cert_pem" {
  description = "PEM-encoded CA certificate"
  type        = "string"
}

variable "aggregator_ca_key_alg" {
  description = "Algorithm used to generate aggregator_ca_key"
  type        = "string"
}

variable "aggregator_ca_key_pem" {
  description = "PEM-encoded CA key"
  type        = "string"
}

variable "service_serving_ca_cert_pem" {
  description = "PEM-encoded CA certificate"
  type        = "string"
}

variable "service_serving_ca_key_alg" {
  description = "Algorithm used to generate service_serving_ca_key"
  type        = "string"
}

variable "service_serving_ca_key_pem" {
  description = "PEM-encoded CA key"
  type        = "string"
}

variable "kube_apiserver_url" {
  type = "string"
}

variable "service_cidr" {
  type = "string"
}
