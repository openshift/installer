
locals {
  arch = "x86_64"
  // TODO(displague) use an EquinixMetal proxy
  /*
  coreos_baseurl = "http://mirror.openshift.com/pub/openshift-v4/${local.arch}/dependencies/rhcos"
  coreos_url     = "${local.coreos_baseurl}/${var.ocp_version}/${var.ocp_version}.${var.ocp_version_zstream}"
  coreos_filenm  = "rhcos-${var.ocp_version}.${var.ocp_version_zstream}-${local.arch}"
  coreos_img     = "${local.coreos_filenm}-metal.${local.arch}.raw.gz"
  coreos_kernel  = "${local.coreos_filenm}-installer-kernel-${local.arch}"
  coreos_initrd  = "${local.coreos_filenm}-installer-initramfs.${local.arch}.img"
  */

  // extracting "api.<clustername>" from <clusterdomain>
  external_name = "api-int.${replace(var.cluster_domain, ".${var.base_domain}", "")}.${var.base_domain}"
}

/*

data "template_file" "user_data" {
  template = file("${path.module}/templates/user_data_${var.operating_system}.sh")
}

data "template_file" "ipxe_script" {
  depends_on = [packet_device.bootstrap]
  for_each   = toset(var.nodes)
  template   = file("${path.module}/templates/ipxe.tpl")

  vars = {
    node_type           = each.value
    bootstrap_ip          = packet_device.bootstrap.access_public_ipv4
    ocp_version         = var.ocp_version
    ocp_version_zstream = var.ocp_version_zstream
  }
}

data "template_file" "ignition_append" {
  depends_on = [packet_device.bootstrap]
  for_each   = toset(var.nodes)
  template   = file("${path.module}/templates/ignition-append.json.tpl")

  vars = {
    node_type          = each.value
    bootstrap_ip         = packet_device.bootstrap.access_public_ipv4
    cluster_name       = var.cluster_name
    cluster_basedomain = var.cluster_basedomain
  }
}
*/

resource "packet_device" "bootstrap" {
  hostname   = local.external_name
  plan       = var.plan
  facilities = [var.facility]
  // metro            = var.metro
  operating_system = "custom_ipxe"
  billing_cycle    = var.billing_cycle
  project_id       = var.project_id
  ipxe_script_url  = "https://gist.githubusercontent.com/displague/5282172449a83c7b83821f8f8333a072/raw/0f0d50c744bb758689911d1f8d421b7730c0fb3e/rhcos.ipxe"

  // user_data        = data.template_file.user_data.rendered
  user_data = var.ignition
}

resource "packet_ip_attachment" "node-address" {
  device_id     = packet_device.bootstrap.id
  cidr_notation = "${var.ip_address}/32"
}

/*
resource "null_resource" "dircheck" {

  provisioner "remote-exec" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }


    inline = [
      "while [ ! -d /usr/share/nginx/html ]; do sleep 2; done; ls /usr/share/nginx/html/",
      "while [ ! -f /usr/lib/systemd/system/nfs-server.service ]; do sleep 2; done; ls /usr/lib/systemd/system/nfs-server.service"
    ]
  }
}

resource "null_resource" "ocp_install_ignition" {

  depends_on = [null_resource.dircheck]


  provisioner "remote-exec" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }


    inline = [
      "curl -o /usr/share/nginx/html/${local.coreos_img} ${local.coreos_url}/${local.coreos_img}",
      "curl -o /usr/share/nginx/html/${local.coreos_kernel} ${local.coreos_url}/${local.coreos_kernel}",
      "curl -o /usr/share/nginx/html/${local.coreos_initrd} ${local.coreos_url}/${local.coreos_initrd}",
      "chmod -R 0755 /usr/share/nginx/html/"
    ]
  }
}

resource "null_resource" "ipxe_files" {

  depends_on = [null_resource.dircheck]
  for_each   = data.template_file.ipxe_script

  provisioner "file" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }

    content     = each.value.rendered
    destination = "/usr/share/nginx/html/${each.key}.ipxe"
  }

  provisioner "remote-exec" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }


    inline = [
      "chmod -R 0755 /usr/share/nginx/html/",
    ]
  }
}

resource "null_resource" "ignition_append_files" {

  depends_on = [null_resource.dircheck]
  for_each   = data.template_file.ignition_append

  provisioner "file" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }

    content     = each.value.rendered
    destination = "/usr/share/nginx/html/${each.key}-append.ign"
  }

  provisioner "remote-exec" {

    connection {
      private_key = file(var.ssh_private_key_path)
      host        = packet_device.bootstrap.access_public_ipv4
    }


    inline = [
      "chmod -R 0755 /usr/share/nginx/html/",
    ]
  }
}


output "finished" {
  depends_on = [null_resource.file_uploads, null_resource.ipxe_files]
  value      = "Loadbalancer provisioning finished."
}

*/
