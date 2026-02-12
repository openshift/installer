terraform {
  # Infra manager supports specific Terraform veresions; ensure compatibility
  required_version = ">=1.2.3"
  required_providers {
    google = {
      source = "hashicorp/google"
      version = ">= 4.0.0"
    }
  }
}

provider "google-beta" {
  project = "${var.project}"
  region = "${var.region}"
}

variable "infra_id" {
  type        = string
  description = "OpenShift Installer Infrastructure ID"
}

variable "project" {
  type        = string
  description = "Project ID"
}

variable "region" {
  type        = string
  description = "GCP Region where the resources will be created."
  default     = "us-central1"
}

# Terraform handles lists but the infra-manager --input-values only
# supports scalar types.
# If you require more or less zones, you must manually add them below
# as a single variable for each. You must add the zones to the
# locals `zones` list below.
variable "zone_0" {
  type        = string
  description = "Zone 1 for the instance types."
}

variable "zone_1" {
  type        = string
  description = "Zone 2 for the instance types."
}

variable "zone_2" {
  type        = string
  description = "Zone 3 for the instance types."
}

variable "subnet" {
  type        = string
  description = "Control plane subnet."
}

variable "image" {
  type        = string
  description = "Cluster Image."
}

variable "machine_type" {
  type        = string
  description = "Machine type for the control plane machine."
  default     = "n1-standard-4"
}

variable "disk_size" {
  type        = string
  description = "Size in GB for the root volume."
  default     = "128"
}

variable "disk_type" {
  type        = string
  description = "Type of storage disk for the vm."
  default     = "pd-ssd"
}

variable "service_account_email" {
  type        = string
  description = "Email for the service account attached to the control planes."
}

variable "ignition" {
  type        = string
  description = "The name of the ignition file."
}

data "local_file" "ignition_file" {
  filename = "${path.module}/${var.ignition}"
}

resource "google_compute_instance" "master_0" {
  provider = google-beta

  name = "${var.infra_id}-master-0"
  zone = "${var.zone_0}"
  machine_type = "${var.machine_type}"
  tags = [
    "${var.infra_id}-master"
  ]
  boot_disk {
    auto_delete = true
    initialize_params {
      size = "${var.disk_size}"
      image = "${var.image}"
      type = "${var.disk_type}"
    }
  }
  network_interface {
    subnetwork = "${var.subnet}"
  }
  metadata = {
    user-data = data.local_file.ignition_file.content
  }
  service_account {
    email = "${var.service_account_email}"
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}

resource "google_compute_instance" "master_1" {
  provider = google-beta

  name = "${var.infra_id}-master-1"
  zone = "${var.zone_1}"
  machine_type = "${var.machine_type}"
  tags = [
    "${var.infra_id}-master"
  ]
  boot_disk {
    auto_delete = true
    initialize_params {
      size = "${var.disk_size}"
      image = "${var.image}"
      type = "${var.disk_type}"
    }
  }
  network_interface {
    subnetwork = "${var.subnet}"
  }
  metadata = {
    user-data = data.local_file.ignition_file.content
  }
  service_account {
    email = "${var.service_account_email}"
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}

resource "google_compute_instance" "master_2" {
  provider = google-beta

  name = "${var.infra_id}-master-2"
  zone = "${var.zone_2}"
  machine_type = "${var.machine_type}"
  tags = [
    "${var.infra_id}-master"
  ]
  boot_disk {
    auto_delete = true
    initialize_params {
      size = "${var.disk_size}"
      image = "${var.image}"
      type = "${var.disk_type}"
    }
  }
  network_interface {
    subnetwork = "${var.subnet}"
  }
  metadata = {
    user-data = data.local_file.ignition_file.content
  }
  service_account {
    email = "${var.service_account_email}"
    scopes = ["https://www.googleapis.com/auth/cloud-platform"]
  }
}
