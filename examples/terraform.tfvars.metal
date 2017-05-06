
// e-mail address used to login to Tectonic
tectonic_admin_email = ""

// bcrypt hash of admin password to use with Tectonic Console
tectonic_admin_password_hash = ""

// The base DNS domain of the cluster. Example: `openstack.dev.coreos.systems`.
tectonic_base_domain = ""

// (optional) PEM-encoded CA certificate, used to generate Tectonic Console's server certificate. Optional, if left blank, a CA certificate will be automatically generated.
// tectonic_ca_cert = ""

// (optional) PEM-encoded CA key, used to generate Tectonic Console's server certificate. Optional if tectonic_ca_cert is left blank
// tectonic_ca_key = ""

// Algorithm used to generate tectonic_ca_key. Optional if tectonic_ca_cert is left blank.
tectonic_ca_key_alg = "RSA"

// The Container Linux update channel.
// Examples: `stable`, `beta`, `alpha`
tectonic_cl_channel = "stable"

// A CIDR notation IP range from which to assign pod IPs
tectonic_cluster_cidr = "10.2.0.0/16"

// The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console.
tectonic_cluster_name = ""

// (optional) The path to the etcd CA certificate for TLS communication with etcd.
// tectonic_etcd_ca_cert_path = ""

// (optional) The path to the etcd client certificate for TLS communication with etcd.
// tectonic_etcd_client_cert_path = ""

// (optional) The path to the etcd client key for TLS communication with etcd.
// tectonic_etcd_client_key_path = ""

// The number of etcd nodes to be created.
// If set to zero, the count of etcd nodes will be determined automatically (currently only supported on AWS).
tectonic_etcd_count = "0"

// (optional) List of external etcd v3 servers to connect with (hostnames/IPs only).
// Needs to be set if using an external etcd cluster.
// Example: `["etcd1", "etcd2", "etcd3"]`
// tectonic_etcd_servers = ""

// If set to true, experimental Tectonic assets are being deployed.
tectonic_experimental = false

// Service IP used to reach kube-apiserver inside the cluster
tectonic_kube_apiserver_service_ip = "10.3.0.1"

// Service IP used to reach kube-dns
tectonic_kube_dns_service_ip = "10.3.0.10"

// Service IP used to reach self-hosted etcd
tectonic_kube_etcd_service_ip = "10.3.0.15"

// The path to the tectonic licence file.
tectonic_license_path = ""

// The number of master nodes to be created.
tectonic_master_count = "1"

// CoreOS kernel/initrd version to PXE boot. Must be present in matchbox assets and correspond to the tectonic_cl_channel. Example: `1298.7.0`
tectonic_metal_cl_version = ""

// The domain name which resolves to controller node(s)
tectonic_metal_controller_domain = ""

// Ordered list of controller domain names. Example: `["node2.example.com", "node3.example.com"]`
tectonic_metal_controller_domains = ""

// Ordered list of controller MAC addresses for matching machines. Example: `["52:54:00:a1:9c:ae"]`
tectonic_metal_controller_macs = ""

// Ordered list of controller names. Example: `["node1"]`
tectonic_metal_controller_names = ""

// The domain name which resolves to Tectonic Ingress (i.e. worker node(s))
tectonic_metal_ingress_domain = ""

// Matchbox CA certificate to trust
tectonic_metal_matchbox_ca = ""

// Matchbox client TLS certificate
tectonic_metal_matchbox_client_cert = ""

// Matchbox client TLS key
tectonic_metal_matchbox_client_key = ""

// Matchbox HTTP read-only endpoint (e.g. http://matchbox.example.com:8080)
tectonic_metal_matchbox_http_url = ""

// Matchbox gRPC API endpoint (e.g. matchbox.example.com:8081)
tectonic_metal_matchbox_rpc_endpoint = ""

// Ordered list of worker domain names. Example: `["node2.example.com", "node3.example.com"]`
tectonic_metal_worker_domains = ""

// Ordered list of worker MAC addresses for matching machines. Example: `["52:54:00:b2:2f:86", "52:54:00:c3:61:77"]`
tectonic_metal_worker_macs = ""

// Ordered list of worker names. Example: `["node2", "node3"]`
tectonic_metal_worker_names = ""

// The path the pull secret file in JSON format.
tectonic_pull_secret_path = ""

// A CIDR notation IP range from which to assign service cluster IPs
tectonic_service_cidr = "10.3.0.0/16"

// SSH public key to use as an authorized key. Example: `"ssh-rsa AAAB3N..."`
tectonic_ssh_authorized_key = ""

// If set to true, a vanilla Kubernetes cluster will be deployed, omitting the tectonic assets.
tectonic_vanilla_k8s = false

// The number of worker nodes to be created.
tectonic_worker_count = "3"
