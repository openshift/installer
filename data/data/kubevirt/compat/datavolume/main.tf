resource "kubevirt_data_volume" "data_volume" {
  metadata {
    name      = var.pvc_name
    namespace = var.namespace
    labels    = var.labels
    annotations = {
      "cdi.kubevirt.io/storage.bind.immediate.requested" = "true"
    }
  }
  spec {
    source {
      http {
        url = var.image_url
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

