data "ibm_resource_group" "group" {
  name = var.resource_group
}

data "ibm_pi_workspace" "existing_workspace" {
  count                = var.service_instance_guid != "" ? 1 : 0
  pi_cloud_instance_id = var.service_instance_guid
}

resource "ibm_pi_workspace" "created_service_instance" {
  count                = var.service_instance_guid == "" ? 1 : 0
  pi_name              = "${var.cluster_id}-power-iaas"
  pi_datacenter        = var.powervs_zone
  pi_resource_group_id = data.ibm_resource_group.group.id
  pi_plan              = "public"
}

resource "time_sleep" "wait_for_workspace" {
  count           = var.service_instance_guid == "" ? 1 : 0
  depends_on      = [ibm_pi_workspace.created_workspace]
  create_duration = var.wait_for_workspace
}
data "ibm_pi_workspace" "created_workspace" {
  count                = var.service_instance_guid != "" ? 0 : 1
  pi_cloud_instance_id = resource.ibm_pi_workspace.created_service_instance[0].id
}
