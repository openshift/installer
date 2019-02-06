data "openstack_images_image_v2" "bootstrap_image" {
  name        = "${var.image_name}"
  most_recent = true
}

data "openstack_compute_flavor_v2" "bootstrap_flavor" {
  name = "${var.flavor_name}"
}

data "ignition_systemd_unit" "haproxy_unit" {
  name    = "haproxy.service"
  enabled = true

  content = <<EOF
[Unit]
Description=Load balancer for the OpenShift services

[Service]
ExecStartPre=/sbin/setenforce 0
ExecStartPre=/bin/systemctl disable --now bootkube kubelet progress openshift
ExecStart=/bin/podman run --name haproxy --rm -ti --net=host -v /etc/haproxy:/usr/local/etc/haproxy:ro docker.io/library/haproxy:1.7
ExecStop=/bin/podman stop -t 10 haproxy
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_systemd_unit" "haproxy_unit_watcher" {
  name    = "haproxy-watcher.service"
  enabled = true

  content = <<EOF
[Unit]
Description=HAproxy config updater

[Service]
Type=oneshot
ExecStart=/usr/local/bin/haproxy-watcher.sh

[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_systemd_unit" "haproxy_timer_watcher" {
  name    = "haproxy-watcher.timer"
  enabled = true

  content = <<EOF
[Timer]
OnCalendar=*:0/2

[Install]
WantedBy=timers.target
EOF
}

data "ignition_file" "haproxy_watcher_script" {
  filesystem = "root"
  mode       = "489"                               // 0755
  path       = "/usr/local/bin/haproxy-watcher.sh"

  source {
    source = "data:,%23%21%2Fbin%2Fbash%0A%0Aset%20-x%0A%0Aexport%20KUBECONFIG%3D%2Fopt%2Fopenshift%2Fauth%2Fkubeconfig%0ATEMPLATE%3D%22%7B%7Brange%20.items%7D%7D%7B%7B%5C%24name%3A%3D.metadata.name%7D%7D%7B%7Brange%20.status.conditions%7D%7D%7B%7Bif%20eq%20.type%20%5C%22Ready%5C%22%7D%7D%7B%7Bif%20eq%20.status%20%5C%22True%5C%22%20%7D%7D%7B%7B%5C%24name%7D%7D%7B%7Bend%7D%7D%7B%7Bend%7D%7D%7B%7Bend%7D%7D%20%7B%7Bend%7D%7D%22%0AMASTERS%3D%24%28oc%20get%20nodes%20-l%20node-role.kubernetes.io%2Fmaster%20-ogo-template%3D%22%24TEMPLATE%22%29%0AWORKERS%3D%24%28oc%20get%20nodes%20-l%20node-role.kubernetes.io%2Fworker%20-ogo-template%3D%22%24TEMPLATE%22%29%0A%0Aif%20%5B%5B%20%24MASTERS%20-eq%20%22%22%20%5D%5D%3B%0Athen%0A%20%20%20%20MASTER_LINES%3D%22%0A%20%20%20%20server%20${var.cluster_name}-bootstrap-49500%20${var.cluster_name}-bootstrap.${var.cluster_domain}%20check%20port%2049500%0A%20%20%20%20server%20${var.cluster_name}-bootstrap-6443%20${var.cluster_name}-bootstrap.${var.cluster_domain}%20check%20port%206443%22%0A%20%20%20%20MASTERS%3D%22${var.cluster_name}-master-0%20${var.cluster_name}-master-1%20${var.cluster_name}-master-2%22%0Afi%0A%0Afor%20master%20in%20%24MASTERS%3B%0Ado%0A%20%20%20%20MASTER_LINES%3D%22%24MASTER_LINES%0A%20%20%20%20server%20%24master%20%24master.${var.cluster_domain}%20check%20port%206443%22%0Adone%0A%0Afor%20worker%20in%20%24WORKERS%3B%0Ado%0A%20%20%20%20WORKER_LINES%3D%22%24WORKER_LINES%0A%20%20%20%20server%20%24worker%20%24worker.${var.cluster_domain}%20check%20port%20443%22%0Adone%0A%0Acat%20%3E%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg.new%20%3C%3C%20EOF%0Alisten%20${var.cluster_name}-api-masters%0A%20%20%20%20bind%200.0.0.0%3A6443%0A%20%20%20%20bind%200.0.0.0%3A49500%0A%20%20%20%20mode%20tcp%0A%20%20%20%20balance%20roundrobin%24MASTER_LINES%0A%0Alisten%20${var.cluster_name}-api-workers%0A%20%20%20%20bind%200.0.0.0%3A80%0A%20%20%20%20bind%200.0.0.0%3A443%0A%20%20%20%20mode%20tcp%0A%20%20%20%20balance%20roundrobin%24WORKER_LINES%0AEOF%0A%0A%0Amkdir%20-p%20%2Fetc%2Fhaproxy%0ACHANGED%3D%24%28diff%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg.new%29%0A%0Aif%20%5B%5B%20%21%20-f%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg%20%5D%5D%20%7C%7C%20%5B%5B%20%21%20%24CHANGED%20-eq%20%22%22%20%5D%5D%3B%0Athen%0A%20%20%20%20cp%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg.backup%20%7C%7C%20true%0A%20%20%20%20cp%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg.new%20%2Fetc%2Fhaproxy%2Fhaproxy.cfg%0A%20%20%20%20systemctl%20restart%20haproxy%0Afi%0A"
  }
}

data "ignition_file" "corefile" {
  filesystem = "root"
  mode       = "420"                   // 0644
  path       = "/etc/coredns/Corefile"

  content {
    content = <<EOF
. {
    log
    errors
    reload 10s

    file /etc/coredns/db.${var.cluster_domain} ${var.cluster_name}-api.${var.cluster_domain} {
    }

    file /etc/coredns/db.${var.cluster_domain} _etcd-server-ssl._tcp.${var.cluster_name}.${var.cluster_domain} {
    }

${replace(join("\n", formatlist("    file /etc/coredns/db.${var.cluster_domain} ${var.cluster_name}-etcd-%s.${var.cluster_domain} {\n    upstream /etc/resolv.conf\n    }\n", var.master_port_names)), "master-port-", "")}

    forward . /etc/resolv.conf {
    }
}

${var.cluster_name}.${var.cluster_domain} {
    log
    errors
    reload 10s

    file /etc/coredns/db.${var.cluster_domain} {
        upstream /etc/resolv.conf
    }
}

EOF
  }
}

data "ignition_file" "coredb" {
  filesystem = "root"
  mode       = "420"                                   // 0644
  path       = "/etc/coredns/db.${var.cluster_domain}"

  content {
    content = <<EOF
$ORIGIN ${var.cluster_domain}.
@    3600 IN SOA host-${var.cluster_name}.${var.cluster_domain}. hostmaster (
                                2017042752 ; serial
                                7200       ; refresh (2 hours)
                                3600       ; retry (1 hour)
                                1209600    ; expire (2 weeks)
                                3600       ; minimum (1 hour)
                                )

${var.cluster_name}-api  IN  A  ${var.service_vm_floating_ip}
*.apps.${var.cluster_name}  IN  A  ${var.service_vm_floating_ip}

${replace(join("\n", formatlist("${var.cluster_name}-etcd-%s  IN  CNAME  ${var.cluster_name}-master-%s", var.master_port_names, var.master_port_names)), "master-port-", "")}

${replace(join("\n", formatlist("_etcd-server-ssl._tcp.${var.cluster_name}  8640  IN  SRV  0  10  2380   ${var.cluster_name}-etcd-%s.${var.cluster_domain}.", var.master_port_names)), "master-port-", "")}
EOF
  }
}

data "ignition_systemd_unit" "local_dns" {
  name = "local-dns.service"

  content = <<EOF
[Unit]
Description=Internal DNS serving the required OpenShift records

[Service]
ExecStart=/bin/podman run --rm -i -t -m 128m --net host --cap-add=NET_ADMIN -v /etc/coredns:/etc/coredns:Z openshift/origin-coredns:v4.0 -conf /etc/coredns/Corefile
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
}

data "ignition_user" "core" {
  name = "core"
}

resource "openstack_objectstorage_object_v1" "lb_ignition" {
  container_name = "${var.swift_container}"
  name           = "load-balancer.ign"
  content        = "${var.ignition}"
}

resource "openstack_objectstorage_tempurl_v1" "lb_ignition_tmpurl" {
  container = "${var.swift_container}"
  method    = "get"
  object    = "${openstack_objectstorage_object_v1.lb_ignition.name}"
  ttl       = 3600
}

data "ignition_config" "lb_redirect" {
  append {
    source = "${openstack_objectstorage_tempurl_v1.lb_ignition_tmpurl.url}"
  }

  files = [
    "${data.ignition_file.haproxy_watcher_script.id}",
    "${data.ignition_file.corefile.id}",
    "${data.ignition_file.coredb.id}",
  ]

  systemd = [
    "${data.ignition_systemd_unit.haproxy_unit.id}",
    "${data.ignition_systemd_unit.haproxy_unit_watcher.id}",
    "${data.ignition_systemd_unit.haproxy_timer_watcher.id}",
    "${data.ignition_systemd_unit.local_dns.id}",
  ]

  users = [
    "${data.ignition_user.core.id}",
  ]
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
