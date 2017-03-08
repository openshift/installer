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
}

// The amount of master nodes to be created.
// Example: `1`
variable "tectonic_master_count" {
  type = "string"
}

// The amount of worker nodes to be created.
// Example: `3`
variable "tectonic_worker_count" {
  type = "string"
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
}

// The name of the cluster.
// This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console.
// Example: `demo`
variable "tectonic_cluster_name" {
  type = "string"
}

// ID of existing VPC to build the cluster into.
// Example: `vpc-5c73a334`
variable "tectonic_aws_external_vpc_id" {
  type    = "string"
  default = ""
}

// IP address range to use when creating the cluster VPC.
// Example: `10.0.0.0/16`
variable "tectonic_aws_vpc_cidr_block" {
  type    = "string"
  default = "10.0.0.0/16"
}

// Number of availability zones the cluster should span.
// Example: `3`
variable "tectonic_aws_az_count" {
  type = "string"
}

// EC2 instance type to use for master nodes.
// Example: `m4.large`
variable "tectonic_aws_master_ec2_type" {
  type = "string"
}

// EC2 instance type to use for worker nodes. 
// Example: `m4.large`
variable "tectonic_aws_worker_ec2_type" {
  type = "string"
}

variable "tectonic_ssh_key" {
  type = "string"
}
