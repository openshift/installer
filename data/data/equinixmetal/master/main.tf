resource "packet_device" "node" {
  count = var.node_count

  depends_on       = [var.depends]
  hostname         = format("master%01d.%s.%s", count.index, replace(var.cluster_domain, ".${var.base_domain}", ""), var.base_domain)
  operating_system = "custom_ipxe"
  // ipxe_script_url  = "http://${var.bootstrap_ip}:8080/${var.node_type}.ipxe"
  billing_cycle = "hourly"
  project_id    = var.project_id

  plan       = var.plan
  facilities = [var.facility]
  // metro           = var.metro
  ipxe_script_url = "https://gist.githubusercontent.com/displague/5282172449a83c7b83821f8f8333a072/raw/0f0d50c744bb758689911d1f8d421b7730c0fb3e/rhcos.ipxe"


  user_data = var.ignition
}

resource "packet_ip_attachment" "node-address" {
  count     = var.node_count
  device_id = packet_device.node[count.index].id
  cidr      = "${var.ip_addresses[count.index]}/32"
}

