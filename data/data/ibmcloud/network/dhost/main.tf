locals {
  prefix               = var.cluster_id
  dhosts_master_create = [for dhost in var.dedicated_hosts_master : dhost if lookup(dhost, "id", "") == ""]
  dhosts_master_zones  = [for i, dhost in var.dedicated_hosts_master : var.zones_master[i] if lookup(dhost, "id", "") == ""]
  dhosts_worker_create = [for dhost in var.dedicated_hosts_worker : dhost if lookup(dhost, "id", "") == ""]
  dhosts_worker_zones  = [for i, dhost in var.dedicated_hosts_worker : var.zones_worker[i] if lookup(dhost, "id", "") == ""]
  dhosts_master_merged = [
    for i, dhost in var.dedicated_hosts_master :
    lookup(dhost, "id", "") == ""
    ? ibm_is_dedicated_host.control_plane[index(ibm_is_dedicated_host.control_plane.*.zone, var.zones_master[i])].id
    : dhost.id
  ]
}

############################################
# Dedicated hosts (Control Plane)
############################################

data "ibm_is_dedicated_host_profile" "control_plane" {
  count = length(local.dhosts_master_create)
  name  = local.dhosts_master_create[count.index].profile
}

resource "ibm_is_dedicated_host_group" "control_plane" {
  count = length(local.dhosts_master_create)

  name           = "${local.prefix}-dgroup-control-plane-${local.dhosts_master_zones[count.index]}"
  class          = data.ibm_is_dedicated_host_profile.control_plane[count.index].class
  family         = data.ibm_is_dedicated_host_profile.control_plane[count.index].family
  resource_group = var.resource_group_id
  zone           = local.dhosts_master_zones[count.index]
}

resource "ibm_is_dedicated_host" "control_plane" {
  count = length(local.dhosts_master_create)

  name           = "${local.prefix}-dhost-control-plane-${local.dhosts_master_zones[count.index]}"
  host_group     = ibm_is_dedicated_host_group.control_plane[count.index].id
  profile        = local.dhosts_master_create[count.index].profile
  resource_group = var.resource_group_id

  instance_placement_enabled = true
}

############################################
# Dedicated hosts (Compute)
############################################

data "ibm_is_dedicated_host_profile" "compute" {
  count = length(local.dhosts_worker_create)
  name  = local.dhosts_worker_create[count.index].profile
}

resource "ibm_is_dedicated_host_group" "compute" {
  count = length(local.dhosts_worker_create)

  name           = "${local.prefix}-dgroup-compute-${local.dhosts_worker_zones[count.index]}"
  class          = data.ibm_is_dedicated_host_profile.compute[count.index].class
  family         = data.ibm_is_dedicated_host_profile.compute[count.index].family
  resource_group = var.resource_group_id
  zone           = local.dhosts_worker_zones[count.index]
}

resource "ibm_is_dedicated_host" "compute" {
  count = length(local.dhosts_worker_create)

  name           = "${local.prefix}-dhost-compute-${local.dhosts_worker_zones[count.index]}"
  host_group     = ibm_is_dedicated_host_group.compute[count.index].id
  profile        = local.dhosts_worker_create[count.index].profile
  resource_group = var.resource_group_id

  instance_placement_enabled = true
}