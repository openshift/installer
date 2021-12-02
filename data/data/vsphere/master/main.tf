locals {
  description = "Created By OpenShift Installer"
  zone_keys = keys(var.vcenter_region_zone_map)
}

provider "vsphere" {
  user                 = var.vsphere_username
  password             = var.vsphere_password
  vsphere_server       = var.vsphere_url
  allow_unverified_ssl = false
}


resource "vsphere_virtual_machine" "vm_master" {
  count = var.master_count

  name                 = "${var.cluster_id}-master-${count.index}"

// TODO: This will need to change to whatever is defined in
// in the new platform vcenters,region,zone section

  resource_pool_id     = var.cluster[local.zone_keys[count.index]].resource_pool_id
  datastore_id         = var.datastore[local.zone_keys[count.index]].id

  num_cpus             = var.vsphere_control_plane_num_cpus
  num_cores_per_socket = var.vsphere_control_plane_cores_per_socket
  memory               = var.vsphere_control_plane_memory_mib

  guest_id             = var.template[local.zone_keys[count.index]].guest_id

// TODO: Change the key if possible.

  folder               = "test-cluster-v79xn-us-east-2"
  enable_disk_uuid     = "true"
  annotation           = local.description

  wait_for_guest_net_timeout  = "0"
  wait_for_guest_net_routable = "false"

  network_interface {

// TODO: don't do this...Use Robert's fix
    network_id = var.template[local.zone_keys[count.index]].network_interfaces.0.network_id
  }

  disk {
    label            = "disk0"
// This needs to come from the machinepool
    size             = var.vsphere_control_plane_disk_gib
    eagerly_scrub    = var.template[local.zone_keys[count.index]].disks.0.eagerly_scrub
    thin_provisioned = var.template[local.zone_keys[count.index]].disks.0.thin_provisioned
  }

  clone {
    template_uuid = var.template[local.zone_keys[count.index]].uuid
  }

  extra_config = {
    "guestinfo.ignition.config.data"          = base64encode(var.ignition_master)
    "guestinfo.ignition.config.data.encoding" = "base64"
    "guestinfo.hostname"                      = "${var.cluster_id}-master-${count.index}"
  }

  tags = var.tags
}

