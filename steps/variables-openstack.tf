variable "tectonic_openstack_cloud" {
  description = "(optional) An entry in a clouds.yaml file."
  type        = "string"
  default     = "default"
}

variable "tectonic_openstack_etcd_extra_sg_ids" {
  description = "(optional) List of additional security group IDs for etcd nodes."
  type        = "list"
  default     = []
}

variable "tectonic_openstack_etcd_flavor_name" {
  description = "Name of the openstack image flavor to use."
  type        = "string"
}

variable "tectonic_openstack_etcd_image_name" {
  description = "(optional) Name of the openstack image name to use."
  type        = "string"
  default     = "RHCOS"
}

variable "tectonic_openstack_key_pair" {
  description = "(optional) Name of the openstack key pair to use."
  type        = "string"
  default     = ""
}

variable "tectonic_openstack_master_extra_sg_ids" {
  description = "(optional) List of additional security group IDs for master nodes."
  type        = "list"
  default     = []
}

variable "tectonic_openstack_master_flavor_name" {
  description = "Name of the openstack image flavor to use."
  type        = "string"
}

variable "tectonic_openstack_master_image_name" {
  description = "(optional) Name of the openstack image name to use."
  type        = "string"
  default     = "RHCOS"
}

variable "tectonic_openstack_ssh_key" {
  description = "Contents of an SSH public key to allow access to the core user"
  type        = "string"
}

variable "tectonic_openstack_worker_extra_sg_ids" {
  description = "(optional) List of additional security group IDs for worker nodes."
  type        = "list"
  default     = []
}

variable "tectonic_openstack_worker_flavor_name" {
  description = "Name of the openstack image flavor to use."
  type        = "string"
}

variable "tectonic_openstack_worker_image_name" {
  description = "(optional) Name of the openstack image name to use."
  type        = "string"
  default     = "RHCOS"
}
