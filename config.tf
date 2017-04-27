variable "tectonic_config_version" {
  description = <<EOF
This declares the version of the global configuration variables.
It has no impact on generated assets but declares the version contract of the configuration.
EOF

  default = "1.0"
}

terraform {
  required_version = ">= 0.9.4"
}

variable "tectonic_container_images" {
  description = "Container images to use"
  type        = "map"

  default = {
    hyperkube                    = "quay.io/coreos/hyperkube:v1.6.2_coreos.0"
    pod_checkpointer             = "quay.io/coreos/pod-checkpointer:20cf8b9a6018731a0770192f30dfa7a1941521e3"
    bootkube                     = "quay.io/coreos/bootkube:v0.4.1"
    console                      = "quay.io/coreos/tectonic-console:v1.4.2"
    identity                     = "quay.io/coreos/dex:v2.3.0"
    kube_version_operator        = "quay.io/coreos/kube-version-operator:7da46d189c36092f43d07ca381a61897402fa13c"
    tectonic_channel_operator    = "quay.io/coreos/tectonic-channel-operator:15c001bd7c008a04394390d08ac71046e723ac48"
    node_agent                   = "quay.io/coreos/node-agent:53f6c8dcc7657b49d1468f7e24933d3897ae8ea7"
    prometheus_operator          = "quay.io/coreos/prometheus-operator:v0.8.2"
    tectonic_prometheus_operator = "quay.io/coreos/tectonic-prometheus-operator:v1.1.0"
    node_exporter                = "quay.io/prometheus/node-exporter:v0.13.0"
    config_reload                = "quay.io/coreos/configmap-reload:v0.0.1"
    heapster                     = "gcr.io/google_containers/heapster:v1.3.0"
    addon_resizer                = "gcr.io/google_containers/addon-resizer:1.7"
    stats_emitter                = "quay.io/coreos/tectonic-stats:6e882361357fe4b773adbf279cddf48cb50164c1"
    stats_extender               = "quay.io/coreos/tectonic-stats-extender:487b3da4e175da96dabfb44fba65cdb8b823db2e"
    error_server                 = "quay.io/coreos/tectonic-error-server:1.0"
    ingress_controller           = "gcr.io/google_containers/nginx-ingress-controller:0.9.0-beta.3"
    kubedns                      = "gcr.io/google_containers/k8s-dns-kube-dns-amd64:1.14.1"
    kubednsmasq                  = "gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:1.14.1"
    kubedns_sidecar              = "gcr.io/google_containers/k8s-dns-sidecar-amd64:1.14.1"
    flannel                      = "quay.io/coreos/flannel:v0.7.1-amd64"
    etcd                         = "quay.io/coreos/etcd:v3.1.6"
    awscli                       = "quay.io/coreos/awscli:025a357f05242fdad6a81e8a6b520098aa65a600"
  }
}

variable "tectonic_versions" {
  description = "Versions of the components to use"
  type        = "map"

  default = {
    etcd       = "v3.1.6"
    prometheus = "v1.6.1"
    monitoring = "1.1.0"
    kubernetes = "1.6.2+tectonic.0"
    tectonic   = "1.6.2-tectonic.0"
  }
}

variable "tectonic_kube_apiserver_service_ip" {
  type        = "string"
  description = "Service IP used to reach kube-apiserver inside the cluster"
  default     = "10.3.0.1"
}

variable "tectonic_kube_dns_service_ip" {
  type        = "string"
  description = "Service IP used to reach kube-dns"
  default     = "10.3.0.10"
}

variable "tectonic_service_cidr" {
  description = "A CIDR notation IP range from which to assign service cluster IPs"
  type        = "string"
  default     = "10.3.0.0/16"
}

variable "tectonic_cluster_cidr" {
  description = "A CIDR notation IP range from which to assign pod IPs"
  type        = "string"
  default     = "10.2.0.0/16"
}

variable "tectonic_master_count" {
  type        = "string"
  description = "The number of master nodes to be created."
  default     = "1"
}

variable "tectonic_worker_count" {
  type        = "string"
  description = "The number of worker nodes to be created."
  default     = "3"
}

variable "tectonic_etcd_count" {
  type        = "string"
  default     = "-1"
  description = "The number of etcd nodes to be created. If not set, the count of etcd nodes will be determined automatically (currently only supported on AWS)."
}

variable "tectonic_etcd_servers" {
  description = "List of external etcd v3 servers to connect with (hostnames/IPs only). Optionally used if using an external etcd cluster."
  type        = "list"
  default     = [""]
}

variable "tectonic_etcd_ca_cert_path" {
  description = "The path to the etcd CA certificate for TLS communication with etcd (optional)."
  type        = "string"
  default     = ""
}

variable "tectonic_etcd_client_cert_path" {
  description = "The path to the etcd client certificate for TLS communication with etcd (optional)."
  type        = "string"
  default     = ""
}

variable "tectonic_etcd_client_key_path" {
  description = "The path to the etcd client key for TLS communication with etcd (optional)."
  type        = "string"
  default     = ""
}

variable "tectonic_base_domain" {
  type        = "string"
  description = "The base DNS domain of the cluster. Example: `openstack.dev.coreos.systems`."
}

variable "tectonic_cluster_name" {
  type        = "string"
  description = "The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console."
}

variable "tectonic_pull_secret_path" {
  type        = "string"
  description = "The path the pull secret file in JSON format."
}

variable "tectonic_license_path" {
  type        = "string"
  description = "The path to the tectonic licence file."
}

variable "tectonic_cl_channel" {
  type    = "string"
  default = "stable"

  description = <<EOF
The Container Linux update channel.
Examples: `stable`, `beta`, `alpha`
EOF
}

variable "tectonic_update_server" {
  type        = "string"
  default     = "https://public.update.core-os.net"
  description = "The URL of the Tectonic Omaha update server"
}

variable "tectonic_update_channel" {
  type        = "string"
  default     = "tectonic-1.5"
  description = "The Tectonic Omaha update channel"
}

variable "tectonic_update_app_id" {
  type        = "string"
  default     = "6bc7b986-4654-4a0f-94b3-84ce6feb1db4"
  description = "The Tectonic Omaha update App ID"
}

variable "tectonic_admin_email" {
  type        = "string"
  description = "e-mail address used to login to Tectonic"
}

variable "tectonic_admin_password_hash" {
  type        = "string"
  description = "bcrypt hash of admin password to use with Tectonic Console"
}

variable "tectonic_ca_cert" {
  type        = "string"
  description = "PEM-encoded CA certificate, used to generate Tectonic Console's server certificate. Optional, if left blank, a CA certificate will be automatically generated."
  default     = ""
}

variable "tectonic_ca_key" {
  type        = "string"
  description = "PEM-encoded CA key, used to generate Tectonic Console's server certificate. Optional if tectonic_ca_cert is left blank"
  default     = ""
}

variable "tectonic_ca_key_alg" {
  type        = "string"
  description = "Algorithm used to generate tectonic_ca_key. Optional if tectonic_ca_cert is left blank."
  default     = "RSA"
}

variable "tectonic_vanilla_k8s" {
  description = "If set to true, a vanilla Kubernetes cluster will be deployed, omitting the tectonic assets."
  default     = false
}

variable "tectonic_experimental" {
  description = "If set to true, experimental Tectonic assets are being deployed."
  default     = false
}
