variable "apiserver_cert_pem" {
  type        = "string"
  description = "The API server certificate in PEM format."
}

variable "apiserver_key_pem" {
  type        = "string"
  description = "The API server key in PEM format."
}

variable "apiserver_proxy_cert_pem" {
  type        = "string"
  description = "The API server proxy certificate in PEM format."
}

variable "apiserver_proxy_key_pem" {
  type        = "string"
  description = "The API server proxy key in PEM format."
}

variable "cloud_provider_config" {
  description = "Content of cloud provider config"
  type        = "string"
  default     = ""
}

variable "cluster_name" {
  type = "string"
}

variable "container_images" {
  description = "Container images to use"
  type        = "map"
}

variable "etcd_ca_cert_pem" {
  type        = "string"
  description = "The etcd CA certificate in PEM format."
}

variable "etcd_client_cert_pem" {
  type        = "string"
  description = "The etcd client certificate in PEM format."
}

variable "etcd_client_key_pem" {
  type        = "string"
  description = "The etcd client key in PEM format."
}

variable "kube_apiserver_url" {
  description = "URL used to reach kube-apiserver"
  type        = "string"
}

variable "root_ca_cert_pem" {
  type        = "string"
  description = "The Root CA in PEM format."
}

variable "aggregator_ca_cert_pem" {
  type        = "string"
  description = "The Aggregated API Server CA in PEM format."
}

variable "kube_ca_cert_pem" {
  type        = "string"
  description = "The Kubernetes CA in PEM format."
}

variable "kube_ca_key_pem" {
  type        = "string"
  description = "The Kubernetes CA key in PEM format."
}

variable "admin_cert_pem" {
  type        = "string"
  description = "The kubelet certificate in PEM format."
}

variable "admin_key_pem" {
  type        = "string"
  description = "The kubelet key in PEM format."
}

variable "oidc_ca_cert" {
  type = "string"
}

variable "service_cidr" {
  description = "A CIDR notation IP range from which to assign service cluster IPs"
  type        = "string"
}

variable "pull_secret_path" {
  type        = "string"
  description = "Path on disk to your Tectonic pull secret. Obtain this from your Tectonic Account: https://account.coreos.com."
  default     = "/Users/coreos/Desktop/config.json"
}
