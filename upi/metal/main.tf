# ================MATCHBOX=====================

locals {
  kernel_args = [
    "console=tty0",
    "console=ttyS1,115200n8",
    "rd.neednet=1",

    # "rd.break=initqueue"
    "coreos.inst=yes",

    "coreos.inst.image_url=${var.pxe_os_image_url}",
    "coreos.inst.install_dev=sda",
    "coreos.inst.skip_media_check",
  ]

  pxe_kernel = var.pxe_kernel_url
  pxe_initrd = var.pxe_initrd_url
}

provider "matchbox" {
  endpoint    = var.matchbox_rpc_endpoint
  client_cert = file(var.matchbox_client_cert)
  client_key  = file(var.matchbox_client_key)
  ca          = file(var.matchbox_trusted_ca_cert)
}

resource "matchbox_profile" "default" {
  name = var.cluster_id
}

resource "matchbox_group" "default" {
  name    = var.cluster_id
  profile = matchbox_profile.default.name
}

resource "matchbox_profile" "master" {
  name   = "${var.cluster_id}-master"
  kernel = local.pxe_kernel

  initrd = [
    local.pxe_initrd,
  ]

  args = concat(
    local.kernel_args,
    ["coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?cluster_id=${var.cluster_id}&role=master"],
  )

  raw_ignition = file(var.master_ign_file)
}

resource "matchbox_profile" "worker" {
  name   = "${var.cluster_id}-worker"
  kernel = local.pxe_kernel

  initrd = [
    local.pxe_initrd,
  ]

  args = concat(
    local.kernel_args,
    ["coreos.inst.ignition_url=${var.matchbox_http_endpoint}/ignition?cluster_id=${var.cluster_id}&role=worker"],
  )

  raw_ignition = file(var.worker_ign_file)
}

resource "matchbox_group" "master" {
  name    = "${var.cluster_id}-master"
  profile = matchbox_profile.master.name

  selector = {
    cluster_id = var.cluster_id
    role       = "master"
  }
}

resource "matchbox_group" "worker" {
  name    = "${var.cluster_id}-worker"
  profile = matchbox_profile.worker.name

  selector = {
    cluster_id = var.cluster_id
    role       = "worker"
  }
}

# ================PACKET=====================

provider "packet" {}

locals {
  packet_facility = "sjc1"
}

resource "packet_device" "masters" {
  count            = var.master_count
  hostname         = "master-${count.index}.${var.cluster_domain}"
  plan             = "c1.small.x86"
  facilities       = ["any"]
  operating_system = "custom_ipxe"
  ipxe_script_url  = "${var.matchbox_http_endpoint}/ipxe?cluster_id=${var.cluster_id}&role=master"
  billing_cycle    = "hourly"
  project_id       = var.packet_project_id

  depends_on = [matchbox_group.master]
}

resource "packet_device" "workers" {
  count            = var.worker_count
  hostname         = "worker-${count.index}.${var.cluster_domain}"
  plan             = "c1.small.x86"
  facilities       = ["any"]
  operating_system = "custom_ipxe"
  ipxe_script_url  = "${var.matchbox_http_endpoint}/ipxe?cluster_id=${var.cluster_id}&role=worker"
  billing_cycle    = "hourly"
  project_id       = var.packet_project_id

  depends_on = [matchbox_group.worker]
}

# ==============BOOTSTRAP=================

module "bootstrap" {
  source = "./bootstrap"

  pxe_kernel             = local.pxe_kernel
  pxe_initrd             = local.pxe_initrd
  pxe_kernel_args        = local.kernel_args
  matchbox_http_endpoint = var.matchbox_http_endpoint
  igntion_config_content = file(var.bootstrap_ign_file)

  cluster_id = var.cluster_id

  packet_facility   = "any"
  packet_project_id = var.packet_project_id
}

# ================AWS=====================

provider aws {
  region = "us-east-1"
}

locals {
  master_public_ipv4_networks = flatten(packet_device.masters.*.network)
  master_public_ipv4          = data.template_file.master_ips.*.rendered

  worker_public_ipv4_networks = flatten(packet_device.workers.*.network)
  worker_public_ipv4          = data.template_file.worker_ips.*.rendered
  ctrp_records                = compact(concat(var.bootstrap_dns ? [module.bootstrap.device_ip] : [], local.master_public_ipv4))
}

data "template_file" "master_ips" {
  count    = var.master_count
  template = lookup(local.master_public_ipv4_networks[count.index * 3], "address")
}

data "template_file" "worker_ips" {
  count    = var.worker_count
  template = lookup(local.worker_public_ipv4_networks[count.index * 3], "address")
}

data "aws_route53_zone" "public" {
  name = var.public_r53_zone
}

resource "aws_route53_record" "ctrlp" {
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "api.${var.cluster_domain}"

  records = local.ctrp_records
}

resource "aws_route53_record" "ctrlp_int" {
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "api-int.${var.cluster_domain}"

  records = local.ctrp_records
}

resource "aws_route53_record" "apps_noworker" {
  count = var.worker_count < 1 ? 1 : 0

  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "*.apps.${var.cluster_domain}"

  records = local.master_public_ipv4
}

resource "aws_route53_record" "apps_worker" {
  count = var.worker_count > 0 ? 1 : 0

  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "*.apps.${var.cluster_domain}"

  records = local.worker_public_ipv4
}

resource "aws_route53_record" "etcd_a_nodes" {
  count   = var.master_count
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "etcd-${count.index}.${var.cluster_domain}"
  records = [local.master_public_ipv4[count.index]]
}

resource "aws_route53_record" "master_a_nodes" {
  count   = var.master_count
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "master-${count.index}.${var.cluster_domain}"
  records = [local.master_public_ipv4[count.index]]
}

resource "aws_route53_record" "worker_a_nodes" {
  count   = var.worker_count
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "A"
  ttl     = "60"
  name    = "worker-${count.index}.${var.cluster_domain}"
  records = [local.worker_public_ipv4[count.index]]
}

resource "aws_route53_record" "etcd_cluster" {
  zone_id = data.aws_route53_zone.public.zone_id
  type    = "SRV"
  ttl     = "60"
  name    = "_etcd-server-ssl._tcp.${var.cluster_domain}"
  records = formatlist("0 10 2380 %s", aws_route53_record.etcd_a_nodes.*.fqdn)
}
