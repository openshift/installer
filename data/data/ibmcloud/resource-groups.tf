############################################
# IBM Cloud resource groups
############################################

# Decide which resource group to use.
# Either create new, use existing, fetch default, or auto-generate new resource group.
locals {
  resource_group_id = (
    var.ibmcloud_resource_group_create ? ibm_resource_group.new.0.id :
      var.ibmcloud_resource_group_name != "" ? data.ibm_resource_group.existing.0.id : data.ibm_resource_group.default.0.id
  )
  resource_group_name = var.ibmcloud_resource_group_name != "" ? var.ibmcloud_resource_group_name : var.cluster_id
}

resource "ibm_resource_group" "new" {
  count = var.ibmcloud_resource_group_create ? 1 : 0
  name  = local.resource_group_name
}

data "ibm_resource_group" "existing" {
  count = !var.ibmcloud_resource_group_create && var.ibmcloud_resource_group_name != "" ? 1 : 0
  name  = local.resource_group_name
}

data "ibm_resource_group" "default" {
  count = !var.ibmcloud_resource_group_create && var.ibmcloud_resource_group_name == "" ? 1 : 0
  is_default = true
}
