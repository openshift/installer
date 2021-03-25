/*
variable "metal_cf_email" {
  description = "Your Cloudflare email address"
}

variable "metal_cf_api_key" {
  description = "Your Cloudflare API key"
}

variable "metal_cf_zone_id" {
  description = "Your Cloudflare Zone"
}
*/

variable "metal_auth_token" {
  description = "Your Equinix Metal API key"
}

variable "metal_project_id" {
  description = "Your Equinix Metal Project ID"
}

variable "metal_ssh_private_key_path" {
  description = "Your SSH private key path (used locally only)"
  default     = "~/.ssh/id_rsa"
}

variable "metal_ssh_public_key_path" {
  description = "Your SSH public key path (used for install-config.yaml)"
  default     = "~/.ssh/id_rsa.pub"
}

variable "metal_bootstrap_operating_system" {
  description = "Your preferred bootstrap operating systems (RHEL or CentOS)"
  default     = "rhel_7"
}

variable "metal_facility" {
  description = "Your primary facility"
}

variable "metal_billing_cycle" {
  description = "Your billing cycle (hourly)"
}

variable "metal_machine_type" {
  description = "Plan for Compute Nodes"
}

/*
variable "metal_ocp_version" {
  default     = "4.6"
  description = "OpenShift minor release version"
}

variable "metal_ocp_version_zstream" {
  default     = "3"
  description = "OpenShift zstream version"
}


variable "metal_ocp_cluster_manager_token" {
  description = "OpenShift Cluster Manager API Token used to generate your pullSecret (https://cloud.redhat.com/openshift/token)"
}

variable "metal_ocp_storage_nfs_enable" {
  description = "Enable configuration of NFS and NFS-related k8s provisioner/storageClass"
  default     = true
}
variable "metal_ocp_storage_ocs_enable" {
  description = "Enable installation of OpenShift Container Storage via operator. This requires a minimum of 3 worker nodes"
  default     = false
}

variable "metal_ocp_virtualization_enable" {
  description = "Enable installation of OpenShift Virtualization via operator. This requires storage provided by OCS, NFS, and/or hostPath provisioner(s)"
  default     = false
}
*/