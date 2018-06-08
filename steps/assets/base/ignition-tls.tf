# CA ignition file assets
data "ignition_file" "root_ca_cert" {
  filesystem = "root"
  mode       = "0644"
  path       = "/opt/tectonic/tls/root-ca.crt"

  content {
    content = "${local.root_ca_cert_pem}"
  }
}

data "ignition_file" "kube_ca_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.kube_ca_key_pem}"
  }

  path = "/opt/tectonic/tls/kube-ca.key"
}

data "ignition_file" "kube_ca_cert" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.kube_ca_cert_pem}"
  }

  path = "/opt/tectonic/tls/kube-ca.crt"
}

data "ignition_file" "aggregator_ca_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.aggregator_ca_key_pem}"
  }

  path = "/opt/tectonic/tls/aggregator-ca.key"
}

data "ignition_file" "aggregator_ca_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.aggregator_ca_cert_pem}"
  }

  path = "/opt/tectonic/tls/aggregator-ca.crt"
}

data "ignition_file" "service_serving_ca_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.service_serving_ca_key_pem}"
  }

  path = "/opt/tectonic/tls/service-serving-ca.key"
}

data "ignition_file" "service_serving_ca_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.service_serving_ca_cert_pem}"
  }

  path = "/opt/tectonic/tls/service-serving-ca.crt"
}

data "ignition_file" "etcd_ca_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.etcd_ca_key_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client-ca.key"
}

data "ignition_file" "etcd_ca_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.etcd_ca_cert_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client-ca.crt"
}

# etcd ignition assets
data "ignition_file" "etcd_client_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.etcd_client_cert_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client.crt"
}

data "ignition_file" "etcd_client_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.etcd_client_key_pem}"
  }

  path = "/opt/tectonic/tls/etcd-client.key"
}

# Kube TLS ignition assets
data "ignition_file" "apiserver_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.apiserver_key_pem}"
  }

  path = "/opt/tectonic/tls/apiserver.key"
}

data "ignition_file" "apiserver_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.apiserver_cert_pem}"
  }

  path = "/opt/tectonic/tls/apiserver.crt"
}

data "ignition_file" "openshift_apiserver_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.openshift_apiserver_key_pem}"
  }

  path = "/opt/tectonic/tls/openshift-apiserver.key"
}

data "ignition_file" "openshift_apiserver_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.openshift_apiserver_cert_pem}"
  }

  path = "/opt/tectonic/tls/openshift-apiserver.crt"
}

data "ignition_file" "apiserver_proxy_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.apiserver_proxy_key_pem}"
  }

  path = "/opt/tectonic/tls/apiserver-proxy.key"
}

data "ignition_file" "apiserver_proxy_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.apiserver_proxy_cert_pem}"
  }

  path = "/opt/tectonic/tls/apiserver-proxy.crt"
}

data "ignition_file" "admin_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${local.admin_key_pem}"
  }

  path = "/opt/tectonic/tls/admin.key"
}

data "ignition_file" "admin_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.admin_cert_pem}"
  }

  path = "/opt/tectonic/tls/admin.crt"
}

data "ignition_file" "kubelet_key" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.kubelet_key_pem}"
  }

  path = "/opt/tectonic/tls/kubelet.key"
}

data "ignition_file" "kubelet_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${local.kubelet_cert_pem}"
  }

  path = "/opt/tectonic/tls/kubelet.crt"
}

locals {
  ca_certs_ignition_file_id_list = [
    "${data.ignition_file.root_ca_cert.id}",
    "${data.ignition_file.kube_ca_key.id}",
    "${data.ignition_file.kube_ca_cert.id}",
    "${data.ignition_file.aggregator_ca_key.id}",
    "${data.ignition_file.aggregator_ca_cert.id}",
    "${data.ignition_file.service_serving_ca_key.id}",
    "${data.ignition_file.service_serving_ca_cert.id}",
    "${data.ignition_file.etcd_ca_key.id}",
    "${data.ignition_file.etcd_ca_cert.id}",
  ]

  etcd_certs_ignition_file_id_list = [
    "${data.ignition_file.etcd_client_cert.id}",
    "${data.ignition_file.etcd_client_key.id}",
  ]

  kube_certs_ignition_file_id_list = [
    "${data.ignition_file.apiserver_key.id}",
    "${data.ignition_file.apiserver_cert.id}",
    "${data.ignition_file.openshift_apiserver_key.id}",
    "${data.ignition_file.openshift_apiserver_cert.id}",
    "${data.ignition_file.apiserver_proxy_key.id}",
    "${data.ignition_file.apiserver_proxy_cert.id}",
    "${data.ignition_file.admin_key.id}",
    "${data.ignition_file.admin_cert.id}",
    "${data.ignition_file.kubelet_key.id}",
    "${data.ignition_file.kubelet_cert.id}",
  ]
}
