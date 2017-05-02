variable "tectonic_metal_config_version" {
  description = <<EOF
(internal) This declares the version of the Matchbox configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

variable "tectonic_metal_matchbox_http_endpoint" {
  type        = "string"
  description = "Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)"
}

variable "tectonic_metal_matchbox_rpc_endpoint" {
  type        = "string"
  description = "Matchbox gRPC API endpoint (e.g. matchbox.example.com:8081)"
}

variable "tectonic_metal_matchbox_client_cert" {
  type        = "string"
  description = "Matchbox client TLS certificate"
}

variable "tectonic_metal_matchbox_client_key" {
  type        = "string"
  description = "Matchbox client TLS key"
}

variable "tectonic_metal_matchbox_ca" {
  type        = "string"
  description = "Matchbox CA certificate to trust"
}

variable "tectonic_metal_cl_version" {
  type        = "string"
  description = "CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets and correspond to the tectonic_cl_channel. Example: `1298.7.0`"
}

variable "tectonic_metal_controller_domain" {
  type        = "string"
  description = "The domain name which resolves to controller node(s)"
}

variable "tectonic_metal_ingress_domain" {
  type        = "string"
  description = "The domain name which resolves to Tectonic Ingress (i.e. worker node(s))"
}

variable "tectonic_metal_controller_names" {
  type        = "list"
  description = "Ordered list of controller names. Example: `[\"node1\"]`"
}

variable "tectonic_metal_controller_domains" {
  type        = "list"
  description = "Ordered list of controller domain names. Example: `[\"node2.example.com\", \"node3.example.com\"]`"
}

variable "tectonic_metal_controller_macs" {
  type        = "list"
  description = "Ordered list of controller MAC addresses for matching machines. Example: `[\"52:54:00:a1:9c:ae\"]`"
}

variable "tectonic_metal_worker_names" {
  type        = "list"
  description = "Ordered list of worker names. Example: `[\"node2\", \"node3\"]`"
}

variable "tectonic_metal_worker_domains" {
  type        = "list"
  description = "Ordered list of worker domain names. Example: `[\"node2.example.com\", \"node3.example.com\"]`"
}

variable "tectonic_metal_worker_macs" {
  type        = "list"
  description = "Ordered list of worker MAC addresses for matching machines. Example: `[\"52:54:00:b2:2f:86\", \"52:54:00:c3:61:77\"]`"
}

variable "tectonic_ssh_authorized_key" {
  type        = "string"
  description = "SSH public key to use as an authorized key. Example: `\"ssh-rsa AAAB3N...\"`"
}
