// The flavor ID as given in `azure flavor list`.
// Specifies the size (CPU/Memory/Drive) of the VM.
variable "tectonic_azure_vm_size" {
  type    = "string"
  default = "Standard_D2_v2"
}

// The image ID as given in `azure image list`.
// Specifies the OS image of the VM.
variable "tectonic_azure_image_reference" {
  type    = "map"
  default = {
     publisher = "CoreOS"
     offer     = "CoreOS"
     sku       = "Stable"
     version   = "latest"
  }
}

// The hyperkube "quay.io/coreos/hyperkube" image version.
variable "tectonic_kube_version" {
  type = "string"
  default = "v1.5.3_coreos.0"
}

// The amount of master nodes to be created.
// Example: `1`
variable "tectonic_master_count" {
  type = "string"
  default = "1"
}

// The amount of worker nodes to be created.
// Example: `3`
variable "tectonic_worker_count" {
  type = "string"
  default = "1"
}

// The amount of etcd nodes to be created.
// Example: `1`
variable "tectonic_etcd_count" {
  type    = "string"
  default = "1"
}

// The base DNS domain of the cluster.
// Example: `azure.dev.coreos.systems`
variable "tectonic_base_domain" {
  type = "string"
  default = "jzawesome"
}

// The name of the cluster.
// This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console.
// Example: `demo`
variable "tectonic_cluster_name" {
  type = "string"
  default = "ecawesome"
}

variable "tectonic_ssh_key" {
  type = "string"
  default = "/Users/jimzim/.ssh/id_rsa.pub"
}

variable "tectonic_region" {
  type = "string"
  default = "East US"
}
