provider "kubevirt" {
}

resource "kubevirt_data_volume" "data_volume" {
  metadata {
    name      = var.kubevirt_source_pvc_name
    namespace = var.kubevirt_namespace
    labels    = var.kubevirt_labels
    annotations = {
      "cdi.kubevirt.io/storage.bind.immediate.requested" = "true"
    }
  }
  spec {
    source {
      http {
        url = var.kubevirt_image_url
      }
    }
    pvc {
      access_modes = [var.kubevirt_pv_access_mode]
      resources {
        requests = {
          storage = var.storage
        }
      }
      storage_class_name = var.kubevirt_storage_class
    }
  }
}

