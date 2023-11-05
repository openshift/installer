data "ibm_resource_group" "group" {
  name = var.resource_group
}

resource "ibm_resource_instance" "powervs_service_instance" {
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
