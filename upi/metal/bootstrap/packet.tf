resource "packet_device" "bootstrap" {
  hostname                = "${var.cluster_id}-bootstrap"
  plan                    = var.packet_plan
  facilities              = [var.packet_facility]
  operating_system        = "custom_ipxe"
  ipxe_script_url         = "${var.matchbox_http_endpoint}/ipxe?cluster_id=${var.cluster_id}&role=bootstrap"
  billing_cycle           = "hourly"
  project_id              = var.packet_project_id
  hardware_reservation_id = var.packet_hardware_reservation_id

  depends_on = [matchbox_group.bootstrap]
}
