variable "tectonic_matchbox_http_endpoint" {
  type = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "tectonic_matchbox_rpc_endpoint" {
  type = "string"
  description = "Matchbox gRPC API endpoint (e.g. matchbox.example.com:8081)"
}

variable "tectonic_matchbox_client_cert" {
  type = "string"
  description = "Matchbox client TLS certificate"
}

variable "tectonic_matchbox_client_key" {
  type = "string"
  description = "Matchbox client TLS key"
}

variable "tectonic_matchbox_ca" {
  type = "string"
  description = "Matchbox CA certificate to trust"
}

variable "tectonic_coreos_version" {
  type = "string"
  description = "CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets."
}

variable "tectonic_k8s_domain" {
  type = "string"
  description = "The domain name which resolves to controller node(s)"
}

variable "tectonic_ingress_domain" {
  type = "string"
  description = "The domain name which resolves to Tectonic Ingress (i.e. worker node(s))"
}
