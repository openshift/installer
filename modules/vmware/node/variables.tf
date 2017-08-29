variable "base_domain" {
  type = "string"
}

variable "bootkube_service" {
  type        = "string"
  description = "The content of the bootkube systemd service unit"
}

variable "container_images" {
  description = "Container images to use"
  type        = "map"
}

variable "ign_kubelet_env_id" {
  type = "string"
}

variable "image_re" {
  description = "(internal) Regular expression used to extract repo and tag components from image strings"
  type        = "string"
}

variable "instance_count" {
  type        = "string"
  description = "Number of nodes to be created."
}

variable "kubeconfig" {
  type        = "string"
  description = "Contents of Kubeconfig"
}

variable "private_key" {
  type        = "string"
  description = "SSH private key file in .pem format corresponding to tectonic_vmware_ssh_authorized_key. If not provided, SSH agent will be used."
  default     = ""
}

variable "tectonic_service" {
  type        = "string"
  description = "The content of the tectonic installer systemd service unit"
}

variable "tectonic_service_disabled" {
  description = "Specifies whether the tectonic installer systemd unit will be disabled. If true, no tectonic assets will be deployed"
  default     = false
}

variable "vmware_folder" {
  type        = "string"
  description = "Name of the VMware folder to create objects in"
}

variable "core_public_keys" {
  type        = "list"
  description = "Public Key for Core User"
}

variable "dns_server" {
  type        = "string"
  description = "DNS Server of the nodes"
}

variable "gateway" {
  type        = "string"
  description = "Gateway of the node"
}

variable "hostname" {
  type        = "map"
  description = "Hostname of the node"
}

variable "ip_address" {
  type        = "map"
  description = "IP Address of the node"
}

variable "vm_disk_datastore" {
  type        = "string"
  description = "Datastore to create VM(s) in "
}

variable "vm_disk_template" {
  type        = "string"
  description = "Disk template to use for cloning CoreOS Container Linux"
}

variable "vm_disk_template_folder" {
  type        = "string"
  description = "vSphere Folder CoreOS Container Linux is located in"
}

variable "vm_memory" {
  type        = "string"
  description = "VMs Memory size in MB"
}

variable "vm_network_label" {
  type        = "string"
  description = "VMs PortGroup"
}

variable "vm_vcpu" {
  type        = "string"
  description = "VMs vCPU count"
}

variable "vmware_cluster" {
  type        = "string"
  description = "vSphere Cluster to create VMs in"
}

variable "vmware_datacenter" {
  type        = "string"
  description = "vSphere Datacenter to create VMs in"
}

variable "ign_kubelet_env_service_id" {
  type        = "string"
  description = "The kubelet env service to use"
}
