provider "ibm" {
  ibmcloud_api_key = var.powervs_api_key
}

data "ibm_resource_group" "group" {
  name = var.powervs_resource_group
}

resource "ibm_resource_instance" "powervs_service_instance" {
  count             = var.powervs_cloud_instance_id == "" ? 1 : 0
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
