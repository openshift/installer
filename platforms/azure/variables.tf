// The flavor ID as given in `azure flavor list`.
// Specifies the size (CPU/Memory/Drive) of the VM.
variable "tectonic_azure_vm_size" {
  type    = "string"
  default = "Standard_DS2"
}

variable "tectonic_azure_dns_resource_group" {
  type    = "string"
  default = "tectonic-dns-group"
}

// The image ID as given in `azure image list`.
// Specifies the OS image of the VM.
variable "tectonic_azure_image_reference" {
  type = "map"

  default = {
    publisher = "CoreOS"
    offer     = "CoreOS"
    sku       = "Stable"
    version   = "latest"
  }
}

variable "tectonic_azure_location" {
  type = "string"
}

variable "tectonic_ssh_key" {
  type    = "string"
  default = ""
}

// Name of an Azure ssh key to use
// joe-sfo
variable "tectonic_azure_ssh_key" {
  type = "string"
}

variable "tectonic_azure_master_vm_size" {
  type = "string"
}

variable "tectonic_azure_worker_vm_size" {
  type = "string"
}

variable "tectonic_azure_etcd_vm_size" {
  type = "string"
}

variable "tectonic_azure_vnet_cidr_block" {
  type = "string"
}

variable "tectonic_azure_external_vnet_id" {
  type    = "string"
  default = ""
}
