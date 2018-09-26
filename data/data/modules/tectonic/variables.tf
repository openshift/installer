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
  description = "Platform on which Tectonic is being installed. Example: aws or libvirt."
  type        = "string"
}

variable "ingress_kind" {
  description = "Type of Ingress mapping to use. Example: HostPort or NodePort."
  type        = "string"
}

variable "pull_secret" {
  type        = "string"
  description = "Your pull secret. Obtain this from your Tectonic Account: https://account.coreos.com."
}

variable "base_address" {
  description = "Base address used to access the Tectonic Console, without protocol nor trailing forward slash (may contain a port). Example: console.example.com:30000."
  type        = "string"
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

variable "ingress_bundle_pem" {
  type = "string"
}
