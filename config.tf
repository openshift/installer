variable "tectonic_container_images" {
  description = "Container images to use"
  type        = "map"

  default = {
    hyperkube                 = "quay.io/coreos/hyperkube:v1.5.5_coreos.0"
    pod_checkpointer          = "quay.io/coreos/pod-checkpointer:5b585a2d731173713fa6871c436f6c53fa17f754"
    bootkube                  = "quay.io/coreos/bootkube:v0.3.12"
    console                   = "quay.io/coreos/tectonic-console:v1.1.1"
    identity                  = "quay.io/coreos/dex:v2.2.4"
    kube_version_operator     = "quay.io/coreos/kube-version-operator:7da46d189c36092f43d07ca381a61897402fa13c"
    tectonic_channel_operator = "quay.io/coreos/tectonic-channel-operator:15c001bd7c008a04394390d08ac71046e723ac48"
    node_agent                = "quay.io/coreos/node-agent:96db3df24c2ddf0b82c38111c71179c5b10de0c9"
    prometheus_operator       = "quay.io/coreos/prometheus-operator:v0.6.0"
    node_exporter             = "quay.io/prometheus/node-exporter:v0.13.0"
    config_reload             = "quay.io/coreos/configmap-reload:v0.0.1"
    heapster                  = "gcr.io/google_containers/heapster:v1.3.0-beta.0"
    addon_resizer             = "gcr.io/google_containers/addon-resizer:1.7"
    stats_emitter             = "quay.io/coreos/tectonic-stats:6e882361357fe4b773adbf279cddf48cb50164c1"
    stats_extender            = "quay.io/coreos/tectonic-stats-extender:487b3da4e175da96dabfb44fba65cdb8b823db2e"
    error_server              = "quay.io/coreos/tectonic-error-server:1.0"
    ingress_controller        = "gcr.io/google_containers/nginx-ingress-controller:0.9.0-beta.3"
    kubedns                   = "gcr.io/google_containers/kubedns-amd64:1.9"
    kubednsmasq               = "gcr.io/google_containers/kube-dnsmasq-amd64:1.4"
    dnsmasq_metrics           = "gcr.io/google_containers/dnsmasq-metrics-amd64:1.0"
    exechealthz               = "gcr.io/google_containers/exechealthz-amd64:1.2"
    flannel                   = "quay.io/coreos/flannel:v0.7.0-amd64"
    etcd                      = "quay.io/coreos/etcd:v3.1.2"
    awscli                    = "quay.io/coreos/awscli:025a357f05242fdad6a81e8a6b520098aa65a600"
  }
}

variable "tectonic_versions" {
  description = "Versions of the components to use"
  type        = "map"

  default = {
    etcd       = "v3.1.2"
    prometheus = "v1.5.2"
    kubernetes = "1.5.5+tectonic.2"
    tectonic   = "1.5.5-tectonic.2"
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

// The amount of master nodes to be created.
// Example: `1`
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

// Example: `1`
variable "tectonic_etcd_count" {
  type        = "string"
  default     = "1"
  description = "The number of etcd nodes to be created."
}

variable "tectonic_etcd_servers" {
  description = "List of external etcd v3 servers to connect with (hostnames/IPs only). Optionally used if using an external etcd cluster."
  type        = "list"
  default     = [""]
}

// The base DNS domain of the cluster.
// Example: `openstack.dev.coreos.systems`
variable "tectonic_base_domain" {
  type = "string"
}

// Example: `demo`
variable "tectonic_cluster_name" {
  type        = "string"
  description = "The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console."
}

variable "tectonic_pull_secret_path" {
  type = "string"
}

variable "tectonic_license_path" {
  type = "string"
}

variable "tectonic_cl_channel" {
  type    = "string"
  default = "stable"
}

variable "tectonic_update_server" {
  type    = "string"
  default = "https://public.update.core-os.net"
}

variable "tectonic_update_channel" {
  type    = "string"
  default = "tectonic-1.5"
}

variable "tectonic_update_app_id" {
  type    = "string"
  default = "6bc7b986-4654-4a0f-94b3-84ce6feb1db4"
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
  default = ""
}

variable "tectonic_ca_key" {
  type        = "string"
  description = "PEM-encoded CA key, used to generate Tectonic Console's server certificate. Optional if tectonic_ca_cert is left blank"
  default = ""
}

variable "tectonic_ca_key_alg" {
  type        = "string"
  description = "Algorithm used to generate tectonic_ca_key. Optional if tectonic_ca_cert is left blank."
  default     = "RSA"
}