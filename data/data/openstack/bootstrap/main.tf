resource "openstack_objectstorage_object_v1" "ignition" {
  container_name = var.swift_container
  name           = "bootstrap.ign"
  content        = var.ignition
}

resource "openstack_objectstorage_tempurl_v1" "ignition_tmpurl" {
  container = var.swift_container
  method    = "get"
  object    = openstack_objectstorage_object_v1.ignition.name
  ttl       = 3600
}

data "ignition_config" "redirect" {
  append {
    source = openstack_objectstorage_tempurl_v1.ignition_tmpurl.url
  }

  files = [
    data.ignition_file.hostname.id,
    data.ignition_file.bootstrap_ifcfg.id,
  ]
}

data "ignition_file" "bootstrap_ifcfg" {
  filesystem = "root"
  mode       = "420" // 0644
  path       = "/etc/sysconfig/network-scripts/ifcfg-eth0"

  content {
    content = <<EOF
DEVICE="eth0"
BOOTPROTO="dhcp"
ONBOOT="yes"
TYPE="Ethernet"
PERSISTENT_DHCLIENT="yes"
DNS1="${var.service_vm_fixed_ip}"
PEERDNS="no"
NM_CONTROLLED="yes"
EOF

  }
}

data "ignition_file" "hostname" {
  filesystem = "root"
  mode = "420" // 0644
  path = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-bootstrap
EOF

}
}

data "openstack_images_image_v2" "bootstrap_image" {
name        = var.image_name
most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
name = var.flavor_name
}

resource "openstack_compute_instance_v2" "bootstrap" {
name      = "${var.cluster_id}-bootstrap"
flavor_id = data.openstack_compute_flavor_v2.bootstrap_flavor.id
image_id  = data.openstack_images_image_v2.bootstrap_image.id

user_data = data.ignition_config.redirect.rendered

network {
port = var.bootstrap_port_id
}

metadata = {
Name = "${var.cluster_id}-bootstrap"
# "kubernetes.io/cluster/${var.cluster_id}" = "owned"
openshiftClusterID = var.cluster_id
}
}

