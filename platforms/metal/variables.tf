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
  description = "CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets and correspond to the tectonic_cl_channel"
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
  description = "Ordered list of controller names (e.g. node1)"
}

variable "tectonic_metal_controller_domains" {
  type        = "list"
  description = "Ordered list of controller domain names (e.g. node1.example.com)"
}

variable "tectonic_metal_controller_macs" {
  type        = "list"
  description = "Ordered list of controller MAC addresses for matching machines"
}

variable "tectonic_metal_worker_names" {
  type        = "list"
  description = "Ordered list of worker names (e.g. node2,node3)"
}

variable "tectonic_metal_worker_domains" {
  type        = "list"
  description = "Ordered list of worker domain names (e.g. node2.example.com,node3.example.com)"
}

variable "tectonic_metal_worker_macs" {
  type        = "list"
  description = "Ordered list of worker MAC addresses for matching machines"
}

variable "tectonic_ssh_authorized_key" {
  type        = "string"
  description = "SSH public key to use as an authorized key"
}
