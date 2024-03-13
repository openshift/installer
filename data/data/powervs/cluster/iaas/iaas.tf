data "ibm_resource_group" "group" {
  name = var.resource_group
}

data "ibm_resource_instance" "existing_service_instance" {
  count             = var.service_instance_name != "" ? 1 : 0
  name              = var.service_instance_name
  service           = "power-iaas"
  resource_group_id = data.ibm_resource_group.group.id
}

resource "ibm_resource_instance" "created_service_instance" {
  count             = var.service_instance_name == "" ? 1 : 0
  name              = "${var.cluster_id}-power-iaas"
  service           = "power-iaas"
  plan              = "power-virtual-server-group"
  location          = var.powervs_zone
  tags              = ["${var.cluster_id}-power-iaas", "${var.cluster_id}"]
  resource_group_id = data.ibm_resource_group.group.id

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}

resource "time_sleep" "wait_for_workspace" {
  count           = var.service_instance_name == "" ? 1 : 0
  depends_on      = [ibm_resource_instance.created_service_instance]
  create_duration = var.wait_for_workspace
}
