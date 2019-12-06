resource "openstack_objectstorage_container_v1" "container" {
  name = var.cluster_id

  # "kubernetes.io/cluster/${var.cluster_id}" = "owned"
  metadata = merge(
    {
      "Name"               = "${var.cluster_id}-ignition"
      "openshiftClusterID" = var.cluster_id
    },

    var.extra_tags,
  )
}

resource "openstack_networking_port_v2" "bootstrap_port" {
  name = "${var.cluster_id}-bootstrap-port"

  admin_state_up     = "true"
  network_id         = var.private_network_id
  security_group_ids = [var.master_sg_id]
  tags               = ["openshiftClusterID=${var.cluster_id}"]

  extra_dhcp_option {
    name  = "domain-search"
    value = var.cluster_domain
  }

  fixed_ip {
    subnet_id = var.nodes_subnet_id
  }

  allowed_address_pairs {
    ip_address = var.api_int_ip
  }

  allowed_address_pairs {
    ip_address = var.node_dns_ip
  }
}

resource "openstack_networking_floatingip_v2" "bootstrap_fip" {
  description = "${var.cluster_id}-bootstrap-fip"
  pool        = var.external_network
  port_id     = openstack_networking_port_v2.bootstrap_port.id
  tags        = ["openshiftClusterID=${var.cluster_id}"]
}


resource "openstack_objectstorage_object_v1" "ignition" {
  container_name = openstack_objectstorage_container_v1.container.name
  name           = "bootstrap.ign"
  content        = var.ignition
}

resource "openstack_objectstorage_tempurl_v1" "ignition_tmpurl" {
  container = openstack_objectstorage_container_v1.container.name
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
    data.ignition_file.dns_conf.id,
    data.ignition_file.dhcp_conf.id,
  ]
}

data "ignition_file" "dhcp_conf" {
  filesystem = "root"
  mode       = "420"
  path       = "/etc/NetworkManager/conf.d/dhcp-client.conf"

  content {
    content = <<EOF
[main]
dhcp=dhclient
EOF
  }
}

data "ignition_file" "dns_conf" {
  filesystem = "root"
  mode = "420"
  path = "/etc/dhcp/dhclient.conf"

  content {
    content = <<EOF
send dhcp-client-identifier = hardware;
prepend domain-name-servers 127.0.0.1;
EOF
  }
}

data "ignition_file" "hostname" {
  filesystem = "root"
  mode       = "420" // 0644
  path       = "/etc/hostname"

  content {
    content = <<EOF
${var.cluster_id}-bootstrap
EOF
  }
}

data "openstack_images_image_v2" "bootstrap_image" {
  name = var.image_name
  most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = var.flavor_name
}

resource "openstack_compute_instance_v2" "bootstrap" {
  name = "${var.cluster_id}-bootstrap"
  flavor_id = data.openstack_compute_flavor_v2.bootstrap_flavor.id
  image_id = data.openstack_images_image_v2.bootstrap_image.id

  user_data = data.ignition_config.redirect.rendered

  network {
    port = openstack_networking_port_v2.bootstrap_port.id
  }

  metadata = {
    Name = "${var.cluster_id}-bootstrap"
    # "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    openshiftClusterID = var.cluster_id
  }
}
