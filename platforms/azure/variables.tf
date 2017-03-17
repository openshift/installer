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
  type    = "string"
  default = "East US"
}

variable "tectonic_ssh_key" {
  type    = "string"
  default = ""
}
