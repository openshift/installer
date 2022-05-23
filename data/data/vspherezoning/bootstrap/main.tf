locals {
  description = "Created By OpenShift Installer"
  vcenter_key = keys(var.vsphere_vcenters)[0]
}

provider "vsphere" {
  user                 = var.vsphere_vcenters[local.vcenter_key].user
  password             = var.vsphere_vcenters[local.vcenter_key].password
  vsphere_server       = var.vsphere_vcenters[local.vcenter_key].server
  allow_unverified_ssl = false
}

resource "vsphere_virtual_machine" "vm_bootstrap" {
  name                        = "${var.cluster_id}-bootstrap"
  resource_pool_id            = var.resource_pool[var.vsphere_control_planes[0].dz_name].id
  datastore_id                = var.datastore[var.vsphere_control_planes[0].dz_name].id
  num_cpus                    = var.vsphere_control_planes[0].provider_spec.numCPUs
  num_cores_per_socket        = var.vsphere_control_planes[0].provider_spec.numCoresPerSocket
  memory                      = var.vsphere_control_planes[0].provider_spec.memoryMiB
  guest_id                    = var.template[var.vsphere_control_planes[0].dz_name].guest_id
  folder                      = var.vsphere_control_planes[0].provider_spec.workspace.folder
  enable_disk_uuid            = "true"
  annotation                  = local.description
  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"
  tags                        = var.tags

  network_interface {
    network_id = var.template[var.vsphere_control_planes[0].dz_name].network_interfaces.0.network_id
  }

  disk {
    label = "disk0"
    size  = var.vsphere_control_planes[0].provider_spec.diskGiB

    eagerly_scrub    = var.template[var.vsphere_control_planes[0].dz_name].disks.0.eagerly_scrub
    thin_provisioned = var.template[var.vsphere_control_planes[0].dz_name].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[var.vsphere_control_planes[0].dz_name].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_bootstrap)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-bootstrap"
    "stealclock.enable"                       = "TRUE"
  }
}
