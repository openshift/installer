variable "apiserver_cert_pem" {
  type        = "string"
  description = "The API server certificate in PEM format."
}

variable "apiserver_key_pem" {
  type        = "string"
  description = "The API server key in PEM format."
}

variable "openshift_apiserver_cert_pem" {
  type        = "string"
  description = "The Openshift API server certificate in PEM format."
}

variable "openshift_apiserver_key_pem" {
  type        = "string"
  description = "The Openshift API server key in PEM format."
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

variable "etcd_endpoints" {
  description = "List of etcd endpoints to connect with (hostnames/IPs only)"
  type        = "list"
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

variable "aggregator_ca_key_pem" {
  type        = "string"
  description = "The Aggregated API Server CA key in PEM format."
}

variable "kube_ca_cert_pem" {
  type        = "string"
  description = "The Kubernetes CA in PEM format."
}

variable "kube_ca_key_pem" {
  type        = "string"
  description = "The Kubernetes CA key in PEM format."
}

variable "service_serving_ca_cert_pem" {
  type        = "string"
  description = "The Service Serving CA in PEM format."
}

variable "service_serving_ca_key_pem" {
  type        = "string"
  description = "The Service Serving CA key in PEM format."
}

variable "admin_cert_pem" {
  type        = "string"
  description = "The admin certificate in PEM format."
}

variable "admin_key_pem" {
  type        = "string"
  description = "The admin key in PEM format."
}

variable "kubelet_cert_pem" {
  type        = "string"
  description = "The kubelet certificate in PEM format."
}

variable "kubelet_key_pem" {
  type        = "string"
  description = "The kubelet key in PEM format."
}

variable "mcs_cert_pem" {
  type = "string"
}

variable "mcs_key_pem" {
  type = "string"
}

variable "service_account_public_key_pem" {
  type = "string"
}

variable "service_account_private_key_pem" {
  type = "string"
}

variable "oidc_ca_cert" {
  type = "string"
}

variable "clusterapi_ca_cert_pem" {
  type = "string"
}

variable "clusterapi_ca_key_pem" {
  type = "string"
}

variable "service_cidr" {
  description = "A CIDR notation IP range from which to assign service cluster IPs"
  type        = "string"
}

variable "pull_secret" {
  type        = "string"
  description = "Your pull secret.  Obtain this from your Tectonic Account: https://account.coreos.com."
}

variable "worker_ign_config" {
  description = "Worker ignition config"
  type        = "string"
  default     = ""
}

variable "libvirt_tls_ca_pem" {
  type        = "string"
  description = "The libvirt CA certificate in PEM format"
  default     = ""
}

variable "libvirt_tls_cert_pem" {
  type        = "string"
  description = "The libvirt client certificate in PEM format"
  default     = ""
}

variable "libvirt_tls_key_pem" {
  type        = "string"
  description = "The libvirt client private key in PEM format"
  default     = ""
}
