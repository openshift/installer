locals {
  description = "Created By OpenShift Installer"
}

provider "nutanix" {
  wait_timeout = 60
  username     = var.nutanix_username
  password     = var.nutanix_password
  endpoint     = var.nutanix_prism_central_address
  port         = var.nutanix_prism_central_port
}

resource "nutanix_image" "bootstrap_ignition" {
  name        = var.nutanix_bootstrap_ignition_image
  source_path = var.nutanix_bootstrap_ignition_image_filepath
  description = local.description

  categories {
    name  = var.ocp_category_key_id
    value = var.ocp_category_value_owned_id
  }
}

resource "nutanix_virtual_machine" "vm_bootstrap" {
  name                 = "${var.cluster_id}-bootstrap"
  description          = local.description
  cluster_uuid         = var.nutanix_prism_element_uuids[0]
  num_vcpus_per_socket = 4
  num_sockets          = 1
  memory_size_mib      = 16384
  boot_device_order_list = [
    "DISK",
    "CDROM",
    "NETWORK"
  ]

  disk_list {
    device_properties {
      device_type = "DISK"
      disk_address = {
        device_index = 0
        adapter_type = "SCSI"
      }
    }
    data_source_reference = {
      kind = "image"
      uuid = var.image_id
    }
    disk_size_mib = var.nutanix_control_plane_disk_mib
  }

  disk_list {
    device_properties {
      device_type = "CDROM"
      disk_address = {
        adapter_type = "IDE"
        device_index = 0
      }
    }
    data_source_reference = {
      kind = "image"
      uuid = nutanix_image.bootstrap_ignition.id
    }
  }

  categories {
    name  = var.ocp_category_key_id
    value = var.ocp_category_value_owned_id
  }

  nic_list {
    subnet_uuid = var.nutanix_subnet_uuids[0]
  }
}
