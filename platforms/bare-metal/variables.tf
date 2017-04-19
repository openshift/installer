variable "matchbox_http_endpoint" {
  type = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "matchbox_rpc_endpoint" {
  type = "string"
  description = "Matchbox gRPC API endpoint (e.g. matchbox.example.com:8081)"
}

variable "matchbox_client_cert" {
  type = "string"
  description = "Matchbox client TLS certificate"
}

variable "matchbox_client_key" {
  type = "string"
  description = "Matchbox client TLS key"
}

variable "matchbox_ca" {
  type = "string"
  description = "Matchbox CA certificate to trust"
}

variable "coreos_version" {
  type = "string"
  description = "CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets."
}
