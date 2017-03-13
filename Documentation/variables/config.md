# Common Tectonic Terraform variables
All the common Tectonic SDK variables used for *all* platforms.

## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| tectonic_admin_email | e-mail address used to login to Tectonic | - | yes |
| tectonic_admin_password_hash | bcrypt hash of admin password to use with Tectonic Console | - | yes |
| tectonic_base_domain | Base address used to access the Tectonic Console, without protocol nor trailing forward slash | - | yes |
| tectonic_ca_cert | PEM-encoded CA certificate, used to generate Tectonic Console's server certificate. Optional, if left blank, a CA certificate will be automatically generated. | - | yes |
| tectonic_ca_key | PEM-encoded CA key, used to generate Tectonic Console's server certificate. Optional if tectonic_ca_cert is left blank | - | yes |
| tectonic_ca_key_alg | Algorithm used to generate tectonic_ca_key. Optional if tectonic_ca_cert is left blank. | `RSA` | no |
| tectonic_cl_channel |  | `stable` | no |
| tectonic_cluster_cidr | A CIDR notation IP range from which to assign pod IPs | `10.2.0.0/16` | no |
| tectonic_cluster_name | The name of the cluster. This will be prepended to `tectonic_base_domain` resulting in the URL to the Tectonic console. | - | yes |
| tectonic_container_images | Container images to use | `<map>` | no |
| tectonic_etcd_count | The number of etcd nodes to be created. | `1` | no |
| tectonic_etcd_servers | List of extenral etcd v3 servers to connect with (scheme://ip:port). Optionally use if providing external etcd. | - | yes |
| tectonic_ingress_type | Type of Ingress mapping to use (e.g. HostPort, NodePort) | `HostPort` | no |
| tectonic_kube_apiserver_service_ip | Service IP used to reach kube-apiserver inside the cluster | `10.3.0.1` | no |
| tectonic_kube_apiserver_url | URL used to reach kube-apiserver | `https://10.1.1.1:443` | no |
| tectonic_kube_dns_service_ip | Service IP used to reach kube-dns | `10.3.0.10` | no |
| tectonic_license |  | - | yes |
| tectonic_master_count | The number of master nodes to be created. | `1` | no |
| tectonic_pull_secret |  | - | yes |
| tectonic_service_cidr | A CIDR notation IP range from which to assign service cluster IPs | `10.3.0.0/16` | no |
| tectonic_update_app_id |  | `6bc7b986-4654-4a0f-94b3-84ce6feb1db4` | no |
| tectonic_update_channel |  | `tectonic-1.5` | no |
| tectonic_update_server |  | `https://public.update.core-os.net` | no |
| tectonic_versions | Versions of the components to use | `<map>` | no |
| tectonic_worker_count | The number of worker nodes to be created. | `3` | no |

