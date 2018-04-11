resource "local_file" "apiserver_key" {
  content  = "${tls_private_key.apiserver.private_key_pem}"
  filename = "./generated/tls/apiserver.key"
}

data "ignition_file" "apiserver_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${tls_private_key.apiserver.private_key_pem}"
  }

  path = "/opt/tectonic/tls/apiserver.key"
}

data "template_file" "apiserver_cert" {
  template = "${join("", list(tls_locally_signed_cert.apiserver.cert_pem, var.kube_ca_cert_pem))}"
}

resource "local_file" "apiserver_cert" {
  content  = "${data.template_file.apiserver_cert.rendered}"
  filename = "./generated/tls/apiserver.crt"
}

data "ignition_file" "apiserver_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${data.template_file.apiserver_cert.rendered}"
  }

  path = "/opt/tectonic/tls/apiserver.crt"
}

resource "local_file" "apiserver_proxy_key" {
  content  = "${tls_private_key.apiserver_proxy.private_key_pem}"
  filename = "./generated/tls/apiserver-proxy.key"
}

data "ignition_file" "apiserver_proxy_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${tls_private_key.apiserver_proxy.private_key_pem}"
  }

  path = "/opt/tectonic/tls/apiserver-proxy.key"
}

resource "local_file" "apiserver_proxy_cert" {
  content  = "${tls_locally_signed_cert.apiserver_proxy.cert_pem}"
  filename = "./generated/tls/apiserver-proxy.crt"
}

data "ignition_file" "apiserver_proxy_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${tls_locally_signed_cert.apiserver_proxy.cert_pem}"
  }

  path = "/opt/tectonic/tls/apiserver-proxy.crt"
}

resource "local_file" "admin_key" {
  content  = "${tls_private_key.admin.private_key_pem}"
  filename = "./generated/tls/admin.key"
}

data "ignition_file" "admin_key" {
  filesystem = "root"
  mode       = "0600"

  content {
    content = "${tls_private_key.admin.private_key_pem}"
  }

  path = "/opt/tectonic/tls/admin.key"
}

resource "local_file" "admin_cert" {
  content  = "${tls_locally_signed_cert.admin.cert_pem}"
  filename = "./generated/tls/admin.crt"
}

data "ignition_file" "admin_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${tls_locally_signed_cert.admin.cert_pem}"
  }

  path = "/opt/tectonic/tls/admin.crt"
}

resource "local_file" "kubelet_key" {
  content  = "${tls_private_key.kubelet.private_key_pem}"
  filename = "./generated/tls/kubelet.key"
}

data "ignition_file" "kubelet_key" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${tls_private_key.kubelet.private_key_pem}"
  }

  path = "/opt/tectonic/tls/kubelet.key"
}

resource "local_file" "kubelet_cert" {
  content  = "${tls_locally_signed_cert.kubelet.cert_pem}"
  filename = "./generated/tls/kubelet.crt"
}

data "ignition_file" "kubelet_cert" {
  filesystem = "root"
  mode       = "0644"

  content {
    content = "${tls_locally_signed_cert.kubelet.cert_pem}"
  }

  path = "/opt/tectonic/tls/kubelet.crt"
}
