data "ignition_systemd_unit" "tx-off" {
  name    = "tx-off.service"
  enable  = true
  content = "${file("${path.module}/resources/tx-off.service")}"
}
