data "ignition_file" "hostname" {
  count = var.master_count
  mode  = "420"
  path  = "/etc/hostname"

  content {
    content = "${var.cluster_id}-master-${count.index}"
  }
}

data "ignition_config" "master_ignition_config" {
  count = var.master_count

  merge {
    source = "data:text/plain;charset=utf-8;base64,${base64encode(var.ignition_data)}"
  }

  files = [
    element(data.ignition_file.hostname.*.rendered, count.index)
  ]
}

resource "kubernetes_secret" "master_ignition" {
  count = var.master_count

  metadata {
    name      = "${var.cluster_id}-master-${count.index}-ignition"
    namespace = var.namespace
    labels    = var.labels
  }
  data = {
    "userdata" = element(
      data.ignition_config.master_ignition_config.*.rendered,
      count.index,
    )
  }
}

locals {
  anti_affinity_label = {
    "anti-affinity-tag-${var.cluster_id}" = "master"
  }
}

resource "kubevirt_virtual_machine" "master_vm" {
  count = var.master_count

  metadata {
    name      = "${var.cluster_id}-master-${count.index}"
    namespace = var.namespace
    labels    = merge(var.labels, local.anti_affinity_label)
  }
  spec {
    run_strategy = "Always"
    data_volume_templates {
      metadata {
        name      = "${var.cluster_id}-master-${count.index}-bootvolume"
        namespace = var.namespace
      }
      spec {
        source {
          pvc {
            name      = var.pvc_name
            namespace = var.namespace
          }
        }
        pvc {
          access_modes = [var.pv_access_mode]
          resources {
            requests = {
              storage = var.storage
            }
          }
          storage_class_name = var.storage_class
        }
      }
    }
    template {
      metadata {
        labels = {
          "kubevirt.io/vm" = "${var.cluster_id}-master-${count.index}"
        }
      }
      spec {
        termination_grace_period_seconds = 600
        volume {
          name = "datavolumedisk1"
          volume_source {
            data_volume {
              name = "${var.cluster_id}-master-${count.index}-bootvolume"
            }
          }
        }
        volume {
          name = "cloudinitdisk"
          volume_source {
            cloud_init_config_drive {
              user_data_secret_ref {
                name = kubernetes_secret.master_ignition[count.index].metadata[0].name
              }
            }
          }
        }
        domain {
          resources {
            requests = {
              memory = var.memory
              cpu    = var.cpu
            }
          }
          devices {
            disk {
              name = "datavolumedisk1"
              disk_device {
                disk {
                  bus = "virtio"
                }
              }
            }
            disk {
              name = "cloudinitdisk"
              disk_device {
                disk {
                  bus = "virtio"
                }
              }
            }
            interface {
              name                     = "main"
              interface_binding_method = var.interface_binding_method
            }
          }
        }
        network {
          name = "main"
          network_source {
            multus {
              network_name = var.network_name
            }
          }
        }
        affinity {
          pod_anti_affinity {
            preferred_during_scheduling_ignored_during_execution {
              weight = 100
              pod_affinity_term {
                label_selector {
                  match_labels = local.anti_affinity_label
                }
                topology_key = "kubernetes.io/hostname"
              }
            }
          }
        }
      }
    }
  }
}
