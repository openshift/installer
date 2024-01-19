data "ibm_resource_group" "group" {
  name = var.resource_group
}

data "ibm_resource_instance" "powervs_service_instance" {
  name = var.service_instance_name == "" ? "${var.cluster_id}-power-iaas" : var.service_instance_name
}
