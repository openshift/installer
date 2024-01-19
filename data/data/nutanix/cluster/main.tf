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

resource "nutanix_category_key" "ocp_category_key" {
  name        = "kubernetes-io-cluster-${var.cluster_id}"
  description = "Openshift Cluster Category Key"
}

resource "nutanix_category_value" "ocp_category_value_owned" {
  name        = nutanix_category_key.ocp_category_key.id
  value       = "owned"
  description = "Openshift Cluster Category Value: resources owned by the cluster"
}

resource "nutanix_category_value" "ocp_category_value_shared" {
  name        = nutanix_category_key.ocp_category_key.id
  value       = "shared"
  description = "Openshift Cluster Category Value: resources used but not owned by the cluster"
}

resource "nutanix_image" "rhcos" {
  name        = var.nutanix_image
  source_uri  = var.nutanix_image_uri
  description = local.description

  categories {
    name  = nutanix_category_key.ocp_category_key.name
    value = nutanix_category_value.ocp_category_value_owned.value
  }
}

data "ignition_file" "hostname" {
  count = var.master_count
  mode  = "420" // 0644
  path  = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-master-${count.index}
EOF
  }
}

data "ignition_config" "master_ignition_config" {
  count = var.master_count

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition_master)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, count.index)
  ]
}

resource "nutanix_virtual_machine" "vm_master" {
  count = var.master_count
  description = local.description
  name = "${var.cluster_id}-master-${count.index}"
  cluster_uuid = var.nutanix_prism_element_uuids[count.index]
  num_vcpus_per_socket = var.nutanix_control_plane_cores_per_socket
  num_sockets = var.nutanix_control_plane_num_cpus
  memory_size_mib = var.nutanix_control_plane_memory_mib
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
      uuid = nutanix_image.rhcos.id
    }
    disk_size_mib = var.nutanix_control_plane_disk_mib
  }

  categories {
    name = nutanix_category_key.ocp_category_key.name
    value = nutanix_category_value.ocp_category_value_owned.value
  }

  dynamic "categories" {
    for_each = (var.nutanix_control_plane_categories == null) ? {} : var.nutanix_control_plane_categories
    content {
      name = categories.key
      value = categories.value
    }
  }

  project_reference = (length(var.nutanix_control_plane_project_uuid) != 0) ? { kind = "project", uuid = var.nutanix_control_plane_project_uuid } : null

  guest_customization_cloud_init_user_data = base64encode(element(data.ignition_config.master_ignition_config.*.rendered, count.index))
  nic_list {
    subnet_uuid = var.nutanix_subnet_uuids[count.index]
  }
}
