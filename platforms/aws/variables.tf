variable "tectonic_aws_ssh_key" {
  type        = "string"
  description = "Name of an SSH key located within the AWS region. Example: coreos-user."
}

variable "tectonic_aws_master_ec2_type" {
  type        = "string"
  description = "Instance size for the master node(s). Example: t2.medium."
  default     = "t2.medium"
}

variable "tectonic_aws_worker_ec2_type" {
  type        = "string"
  description = "Instance size for the worker node(s). Example: t2.medium."
  default     = "t2.medium"
}

variable "tectonic_aws_etcd_ec2_type" {
  type        = "string"
  description = "Instance size for the etcd node(s). Example: t2.medium."
  default     = "t2.medium"
}

variable "tectonic_aws_vpc_cidr_block" {
  type        = "string"
  default     = "10.0.0.0/16"
  description = "Block of IP addresses used by the VPC. This should not overlap with any other networks, such as a private datacenter connected via Direct Connect."
}

variable "tectonic_aws_az_count" {
  type        = "string"
  default     = "3"
  description = "Number of Availability Zones your EC2 instances will be deployed across. This should be less than or equal to the total number available in the region. Be aware that some regions only have 2."
}

variable "tectonic_aws_external_vpc_id" {
  type        = "string"
  description = "ID of an existing VPC to launch nodes into. Example: vpc-123456. Leave blank to create a new VPC."
  default     = ""
}

variable "tectonic_aws_external_vpc_public" {
  description = "If set to true and an external VPC id is given, create public facing ingress resources (ELB, A-records)."
  default     = true
}

variable "tectonic_aws_external_master_subnet_ids" {
  type        = "list"
  description = "List of subnet IDs within an existing VPC to deploy master nodes into. Required to use an existing VPC and the list must match the AZ count. Example: [\"subnet-111111\", \"subnet-222222\", \"subnet-333333\"]"
  default     = [""]
}

variable "tectonic_aws_external_worker_subnet_ids" {
  type        = "list"
  description = "List of subnet IDs within an existing VPC to deploy worker nodes into. Required to use an existing VPC and the list must match the AZ count. Example: [\"subnet-111111\", \"subnet-222222\", \"subnet-333333\"]"
  default     = [""]
}
