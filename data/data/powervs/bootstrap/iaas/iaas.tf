data "ibm_resource_group" "group" {
  name = var.resource_group
}

data "ibm_pi_workspace" "powervs_service_instance" {
  pi_cloud_instance_id = var.service_instance_guid == "" ? var.service_instance_guid : var.service_instance_guid
}