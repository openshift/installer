data "openstack_images_image_v2" "bootstrap_image" {
  name        = "${var.image_name}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = "${var.flavor_name}"
}

data "ignition_systemd_unit" "haproxy_unit" {
  name    = "bootkube-haproxy.service"
  enabled = true

  content = <<EOF
[Unit]
Description=Load balancer for the OpenShift services

[Service]
ExecStartPre=/sbin/setenforce 0
ExecStart=/bin/podman run --name haproxy --rm -ti --net=host -v /etc/haproxy:/usr/local/etc/haproxy:ro docker.io/library/haproxy:1.7
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_file" "haproxy_conf" {
  filesystem = "root"
  path       = "/etc/haproxy/haproxy.cfg"

  source {
    source = "data:,listen%20${var.cluster_name}-api-80%0D%0A%20%20%20%20bind%200.0.0.0%3A80%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20${var.cluster_name}-bootstrap%20${var.cluster_name}-bootstrap.${var.cluster_domain}%3A80%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-0%20${var.cluster_name}-master-0.${var.cluster_domain}%3A80%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-1%20${var.cluster_name}-master-1.${var.cluster_domain}%3A80%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-2%20${var.cluster_name}-master-2.${var.cluster_domain}%3A80%20check%0D%0A%0D%0Alisten%20${var.cluster_name}-api-6443%0D%0A%20%20%20%20bind%200.0.0.0%3A6443%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20${var.cluster_name}-bootstrap%20${var.cluster_name}-bootstrap.${var.cluster_domain}%3A6443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-0%20${var.cluster_name}-master-0.${var.cluster_domain}%3A6443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-1%20${var.cluster_name}-master-1.${var.cluster_domain}%3A6443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-2%20${var.cluster_name}-master-2.${var.cluster_domain}%3A6443%20check%0D%0A%0D%0Alisten%20${var.cluster_name}-api-443%0D%0A%20%20%20%20bind%200.0.0.0%3A443%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20${var.cluster_name}-bootstrap%20${var.cluster_name}-bootstrap.${var.cluster_domain}%3A443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-0%20${var.cluster_name}-master-0.${var.cluster_domain}%3A443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-1%20${var.cluster_name}-master-1.${var.cluster_domain}%3A443%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-2%20${var.cluster_name}-master-2.${var.cluster_domain}%3A443%20check%0D%0A%0D%0Alisten%20${var.cluster_name}-api-49500%0D%0A%20%20%20%20bind%200.0.0.0%3A49500%0D%0A%20%20%20%20mode%20tcp%0D%0A%20%20%20%20stats%20enable%0D%0A%20%20%20%20stats%20uri%20%2Fhaproxy%3Fstatus%0D%0A%20%20%20%20balance%20roundrobin%0D%0A%20%20%20%20server%20${var.cluster_name}-bootstrap%20${var.cluster_name}-bootstrap.${var.cluster_domain}%3A49500%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-0%20${var.cluster_name}-master-0.${var.cluster_domain}%3A49500%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-1%20${var.cluster_name}-master-1.${var.cluster_domain}%3A49500%20check%0D%0A%20%20%20%20server%20${var.cluster_name}-master-2%20${var.cluster_name}-master-2.${var.cluster_domain}%3A49500%20check"
  }
}

data "ignition_file" "openshift_hosts" {
  filesystem = "root"
  mode       = "420"                  // 0644
  path       = "/etc/openshift-hosts"

  content {
    content = <<EOF
${replace(join("\n", formatlist("%s ${var.cluster_name}-etcd-%s.${var.cluster_domain}", var.master_ips, var.master_port_names)), "master-port-", "")}
EOF
  }
}

data "ignition_systemd_unit" "local_dns" {
  name = "local-dns.service"

  content = <<EOF
[Unit]
Description=Internal DNS server for running OpenShift on OpenStack

[Service]
ExecStart=/bin/podman run --name bootstrap-dns --rm -t -i -p 53:53/tcp -p 53:53/udp -v /etc/openshift-hosts:/etc/openshift-hosts:z --cap-add=NET_ADMIN docker.io/andyshinn/dnsmasq:latest --keep-in-foreground --log-facility=- --log-queries --no-resolv --addn-hosts=/etc/openshift-hosts --server=10.0.0.2 ${replace(join(" ", formatlist("--srv-host=_etcd-server-ssl._tcp.${var.cluster_name}.${var.cluster_domain},${var.cluster_name}-etcd-%s.${var.cluster_domain},2380,0,10", var.master_port_names)), "master-port-", "")}

[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_user" "core" {
  name = "core"
}

data "ignition_config" "config" {
  files = [
    "${data.ignition_file.haproxy_conf.id}",
    "${data.ignition_file.openshift_hosts.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.haproxy_unit.id}",
    "${data.ignition_systemd_unit.local_dns.id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
}

resource "openstack_objectstorage_object_v1" "lb_ignition" {
  container_name = "${var.swift_container}"
  name           = "load-balancer.ign"
  content        = "${data.ignition_config.config.rendered}"
}

resource "openstack_objectstorage_tempurl_v1" "lb_ignition_tmpurl" {
  container = "${var.swift_container}"
  method    = "get"
  object    = "${openstack_objectstorage_object_v1.lb_ignition.name}"
  ttl       = 3600
}

data "ignition_config" "lb_redirect" {
  replace {
    source = "${openstack_objectstorage_tempurl_v1.lb_ignition_tmpurl.url}"
  }
}

resource "openstack_compute_instance_v2" "load_balancer" {
  name      = "${var.cluster_name}-api"
  flavor_id = "${data.openstack_compute_flavor_v2.bootstrap_flavor.id}"
  image_id  = "${data.openstack_images_image_v2.bootstrap_image.id}"

  user_data = "${data.ignition_config.lb_redirect.rendered}"

  network {
    port = "${var.lb_port_id}"
  }

  metadata {
    Name = "${var.cluster_name}-bootstrap"

    # "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    openshiftClusterID = "${var.cluster_id}"
  }
}
