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

  content {
    content = <<TFEOF
#!/bin/bash

set -x

export KUBECONFIG=/opt/openshift/auth/kubeconfig
TEMPLATE="{{range .items}}{{\$name:=.metadata.name}}{{range .status.conditions}}{{if eq .type \"Ready\"}}{{if eq .status \"True\" }}{{\$name}}{{end}}{{end}}{{end}} {{end}}"
MASTERS=$(oc get nodes -l node-role.kubernetes.io/master -ogo-template="$TEMPLATE")
WORKERS=$(oc get nodes -l node-role.kubernetes.io/worker -ogo-template="$TEMPLATE")

if [[ $MASTERS -eq "" ]];
then
    MASTER_LINES="
    server ${var.cluster_name}-bootstrap-22623 ${var.cluster_name}-bootstrap.${var.cluster_domain} check port 22623
    server ${var.cluster_name}-bootstrap-6443 ${var.cluster_name}-bootstrap.${var.cluster_domain} check port 6443"
    MASTERS="${var.cluster_name}-master-0 ${var.cluster_name}-master-1 ${var.cluster_name}-master-2"
fi

for master in $MASTERS;
do
    MASTER_LINES="$MASTER_LINES
    server $master $master.${var.cluster_domain} check port 6443"
done

for worker in $WORKERS;
do
    WORKER_LINES="$WORKER_LINES
    server $worker $worker.${var.cluster_domain} check port 443"
done

cat > /etc/haproxy/haproxy.cfg.new << EOF
listen ${var.cluster_name}-api-masters
    bind 0.0.0.0:6443
    bind 0.0.0.0:22623
    mode tcp
    balance roundrobin$MASTER_LINES

listen ${var.cluster_name}-api-workers
    bind 0.0.0.0:80
    bind 0.0.0.0:443
    mode tcp
    balance roundrobin$WORKER_LINES
EOF


mkdir -p /etc/haproxy
CHANGED=$(diff /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.new)

if [[ ! -f /etc/haproxy/haproxy.cfg ]] || [[ ! $CHANGED -eq "" ]];
then
    cp /etc/haproxy/haproxy.cfg /etc/haproxy/haproxy.cfg.backup || true
    cp /etc/haproxy/haproxy.cfg.new /etc/haproxy/haproxy.cfg
    systemctl restart haproxy
fi
TFEOF
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

resource "openstack_objectstorage_object_v1" "service_ignition" {
  container_name = "${var.swift_container}"
  name           = "load-balancer.ign"
  content        = "${var.ignition}"
}

resource "openstack_objectstorage_tempurl_v1" "service_ignition_tmpurl" {
  container = "${var.swift_container}"
  method    = "get"
  object    = "${openstack_objectstorage_object_v1.service_ignition.name}"
  ttl       = 3600
}

data "ignition_config" "service_redirect" {
  append {
    source = "${openstack_objectstorage_tempurl_v1.service_ignition_tmpurl.url}"
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

  user_data = "${data.ignition_config.service_redirect.rendered}"

  network {
    port = "${var.service_port_id}"
  }

  metadata {
    Name = "${var.cluster_name}-bootstrap"

    # "kubernetes.io/cluster/${var.cluster_name}" = "owned"
    openshiftClusterID = "${var.cluster_id}"
  }
}
