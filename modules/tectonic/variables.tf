variable "container_images" {
  description = "Container images to use. Leave blank for defaults."
  type        = "map"
}

variable "container_base_images" {
  description = "Container base images to use. Leave blank for defaults."
  type        = "map"
}

variable "versions" {
  description = "Versions of the components to use. Leave blank for defaults."
  type        = "map"
}

variable "platform" {
  description = "Platform on which Tectonic is being installed. Example: bare-metal or aws."
  type        = "string"
}

variable "ingress_kind" {
  description = "Type of Ingress mapping to use. Example: HostPort or NodePort."
  type        = "string"
}

variable "license_path" {
  description = "Path on disk to your Tectonic license. Obtain this from your Tectonic Account: https://account.coreos.com."
  type        = "string"
  default     = "/Users/coreos/Desktop/tectonic-license.txt"
}

variable "pull_secret_path" {
  type        = "string"
  description = "Path on disk to your Tectonic pull secret. Obtain this from your Tectonic Account: https://account.coreos.com."
  default     = "/Users/coreos/Desktop/config.json"
}

variable "ca_generated" {
  description = "Define whether the CA has been generated or user-provided."
  type        = "string"
}

variable "identity_client_ca_cert" {
  description = "A PEM-encoded CA bundle, used to verify identity server."
  type        = "string"
}

variable "identity_server_ca_cert" {
  description = "A PEM-encoded CA certificate, used to authenticate identity client."
  type        = "string"
}

variable "base_address" {
  description = "Base address used to access the Tectonic Console, without protocol nor trailing forward slash (may contain a port). Example: console.example.com:30000."
  type        = "string"
}

variable "admin_email" {
  description = "E-mail address used to log in to the Tectonic Console."
  type        = "string"
  default     = "admin@example.com"
}

variable "update_server" {
  description = "Server contacted to request Tectonic software updates. Leave blank for defaults."
  type        = "string"
}

variable "update_channel" {
  description = "Release channel used to request Tectonic software updates. Leave blank for defaults. Example: Tectonic-1.5"
  type        = "string"
}

variable "update_app_id" {
  description = "Application identifier used to request Tectonic software updates. Leave blank for defaults."
  type        = "string"
}

variable "ingress_ca_cert_pem" {
  type = "string"
}

variable "ingress_cert_pem" {
  type = "string"
}

variable "ingress_key_pem" {
  type = "string"
}

variable "identity_server_cert_pem" {
  type = "string"
}

variable "identity_server_key_pem" {
  type = "string"
}

variable "identity_client_cert_pem" {
  type = "string"
}

variable "identity_client_key_pem" {
  type = "string"
}
