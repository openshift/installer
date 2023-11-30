data "ibm_resource_group" "group" {
  name = var.resource_group
}

resource "ibm_pi_workspace" "powervs_service_instance" {
  pi_name              = "${var.cluster_id}-power-iaas"
  pi_datacenter        = var.powervs_zone
  pi_resource_group_id = data.ibm_resource_group.group.id
  pi_plan              = "public"
}
