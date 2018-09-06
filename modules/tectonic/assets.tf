# Unique Cluster ID (uuid)
resource "random_id" "cluster_id" {
  byte_length = 16
}

resource "random_string" "ingress_status_password" {
  length  = 6
  special = false
  upper   = false
}

data "ignition_file" "tectonic_sh" {
  filesystem = "root"
  mode       = "0755"
  path       = "/opt/tectonic/tectonic.sh"

  content {
    content = "${file("${path.module}/resources/tectonic.sh")}"
  }
}

data "ignition_systemd_unit" "tectonic_service" {
  name    = "tectonic.service"
  enabled = true
  content = "${file("${path.module}/resources/tectonic.service")}"
}
