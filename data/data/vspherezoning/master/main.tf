locals {
  description = "Created By OpenShift Installer"
  //control_plane_keys = keys(var.vsphere_control_planes_zone)
  vcenter_key = keys(var.vsphere_vcenters)[0]
}


//   for_each = var.vsphere_deployment_zone

//



provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}


// We are supposed to be using master_count...
/*

var.vsphere_control_planes[].DeploymentZone
var.vsphere_control_planes[].ControlPlaneConfig.



*/

resource "vsphere_virtual_machine" "vm_master" {
  count = var.master_count

  name                        = "${var.cluster_id}-master-${count.index}"
  resource_pool_id            = var.resource_pool[var.vsphere_control_planes[count.index].dz_name].id
  datastore_id                = var.datastore[var.vsphere_control_planes[count.index].dz_name].id
  num_cpus                    = var.vsphere_control_planes[count.index].provider_spec.numCPUs
  num_cores_per_socket        = var.vsphere_control_planes[count.index].provider_spec.numCoresPerSocket
  memory                      = var.vsphere_control_planes[count.index].provider_spec.memoryMiB
  guest_id                    = var.template[var.vsphere_control_planes[count.index].dz_name].guest_id
  folder                      = var.vsphere_control_planes[count.index].provider_spec.workspace.folder
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags

  network_interface {
    network_id = var.template[var.vsphere_control_planes[count.index].dz_name].network_interfaces.0.network_id
  }

  disk {
    label = "disk0"
    size  = var.vsphere_control_planes[count.index].provider_spec.diskGiB

    eagerly_scrub    = var.template[var.vsphere_control_planes[count.index].dz_name].disks.0.eagerly_scrub
    thin_provisioned = var.template[var.vsphere_control_planes[count.index].dz_name].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[var.vsphere_control_planes[count.index].dz_name].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-master-${count.index}"
    "stealclock.enable"                       = "TRUE"
  }
}


/*

resource "vsphere_virtual_machine" "vm_master" {

  count = var.master_count
  //for_each = var.vsphere_control_planes_zone

  name                        = "${var.cluster_id}-master-${index(local.control_plane_keys, each.key)}"
  resource_pool_id            = var.resource_pool[each.key].id
  datastore_id                = var.datastore[each.key].id
  num_cpus                    = each.value.numCPUs
  num_cores_per_socket        = each.value.numCoresPerSocket
  memory                      = each.value.memoryMiB
  guest_id                    = var.template[each.key].guest_id
  folder                      = each.value.workspace.folder
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags

  network_interface {
    network_id = var.template[each.key].network_interfaces.0.network_id
  }

  disk {
    label = "disk0"
    size  = each.value.diskGiB

    eagerly_scrub    = var.template[each.key].disks.0.eagerly_scrub
    thin_provisioned = var.template[each.key].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[each.key].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-master-${index(local.control_plane_keys, each.key)}"
    "stealclock.enable"                       = "TRUE"
  }
}

*/