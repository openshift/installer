variable "packet_cf_email" {
  description = "Your Cloudflare email address"
}

variable "packet_cf_api_key" {
  description = "Your Cloudflare API key"
}

variable "packet_cf_zone_id" {
  description = "Your Cloudflare Zone"
}

variable "packet_cluster_basedomain" {
  description = "Your Cloudflare Base domain for your cluster"
}


variable "packet_auth_token" {
  description = "Your Packet API key"
}

variable "packet_project_id" {
  description = "Your Packet Project ID"
}

variable "packet_ssh_private_key_path" {
  description = "Your SSH private key path (used locally only)"
  default     = "~/.ssh/id_rsa"
}

variable "packet_ssh_public_key_path" {
  description = "Your SSH public key path (used for install-config.yaml)"
  default     = "~/.ssh/id_rsa.pub"
}

variable "packet_bastion_operating_system" {
  description = "Your preferred bastion operating systems (RHEL or CentOS)"
  default     = "rhel_7"
}

variable "packet_facility" {
  description = "Your primary facility"
  default     = "dfw2"
}

variable "packet_plan_master" {
  description = "Plan for Master Nodes"
  default     = "c3.medium.x86"
}

variable "packet_plan_compute" {
  description = "Plan for Compute Nodes"
  default     = "c2.medium.x86"
}

variable "packet_count_bootstrap" {
  default     = "1"
  description = "Number of Master Nodes."
}

variable "packet_count_master" {
  default     = "3"
  description = "Number of Master Nodes."
}

variable "packet_count_compute" {
  default     = "2"
  description = "Number of Compute Nodes"
}

variable "packet_cluster_name" {
  default     = "jr"
  description = "Cluster name label"
}

variable "packet_ocp_version" {
  default     = "4.4"
  description = "OpenShift minor release version"
}

variable "packet_ocp_version_zstream" {
  default     = "3"
  description = "OpenShift zstream version"
}

variable "packet_ocp_cluster_manager_token" {
  description = "OpenShift Cluster Manager API Token used to generate your pullSecret (https://cloud.redhat.com/openshift/token)"
}

variable "packet_ocp_storage_nfs_enable" {
  description = "Enable configuration of NFS and NFS-related k8s provisioner/storageClass"
  default     = true
}
variable "packet_ocp_storage_ocs_enable" {
  description = "Enable installation of OpenShift Container Storage via operator. This requires a minimum of 3 worker nodes"
  default     = false
}

variable "packet_ocp_virtualization_enable" {
  description = "Enable installation of OpenShift Virtualization via operator. This requires storage provided by OCS, NFS, and/or hostPath provisioner(s)"
  default     = false
}