variable "container_images" {
    description = "Container images to use"
    type = "map"
}

variable "versions" {
    description = "Versions of the components to use"
    type = "map"
}

variable "platform" {
    description = "Platform on which Tectonic is being installed (e.g. bare-metal, aws)"
    type        = "string"
}

variable "ingress_kind" {
    description = "Type of Ingress mapping to use (e.g. HostPort, NodePort)"
    type        = "string"
}

variable "license_path" {
    description = "Path to the license issued to run Tectonic"
    type        = "string"
}

variable "pull_secret_path" {
    description = "Path to the authentication secret used to pull container images"
    type        = "string"
}

variable "ca_generated" {
    description = "Define whether the CA has been generated or user-provided"
    type        = "string"
}

variable "ca_cert" {
    description = "PEM-encoded CA certificate, used to generate Tectonic Console's server certificate"
    type        = "string"
}

variable "ca_key_alg" {
    description = "Algorithm used to generate ca_key"
    type        = "string"
}

variable "ca_key" {
    description = "PEM-encoded CA key, used to generate Tectonic Console's server certificate"
    type        = "string"
}

variable "base_address" {
    description = "Base address used to access the Tectonic Console, without protocol nor trailing forward slash (may contain a port)"
    type        = "string"
}

variable "admin_email" {
    description = "E-mail address used to login to the Tectonic Console"
    type        = "string"
}

variable "admin_password_hash" {
    description = "Password used to login to the Tectonic Console, hashed by bcrypt"
    type        = "string"
}

variable "update_server" {
    description = "Server contacted to pull updates from"
    type        = "string"
}

variable "update_channel" {
    description = "Channel to pull updates from"
    type        = "string"
}

variable "update_app_id" {
    description = "Application identifier to pull updates for"
    type        = "string"
}

variable "console_client_id" {
    description = "OIDC identifier for the Tectonic Console"
    type        = "string"
}

variable "kubectl_client_id" {
    description = "OIDC identifier for kubectl"
    type        = "string"
}

variable "kube_apiserver_url" {
    description = "URL used to reach kube-apiserver"
    type        = "string"
}