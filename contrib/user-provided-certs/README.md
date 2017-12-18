This package contains a sample terraform project to generate user provided certificates.
It is meant to be a starter for customizations.

If `tectonic_cluster_name` is set to `cluster`,
`tectonic_base_domain` is set to `base.domain`,
and `tectonic_etcd_count` is set to `3`
then the following values have to be configured:

```
$ cat terraform.tfvars
api_fqdn = "cluster-api.base.domain"

console_fqdn = "cluster.base.domain"

etcd_dns_names = [
  "cluster-etcd-0.base.domain",
  "cluster-etcd-1.base.domain",
  "cluster-etcd-2.base.domain",
]

service_cidr = "10.3.0.0/16"
```

Once configured, execute `terraform apply`. The folder `generated/tls` will contain the generated TLS certificates.

The following table gives an overview which generated certificates have to be configured for the corresponding variables for the `user-provided` flavor of TLS modules:

Certificate           | TLS module                           | Variable
----------------------|--------------------------------------|---------
`apiserver.crt`       | `modules/tls/kube/user-provided`     | `apiserver_cert_pem_path`
`apiserver.key`       | `modules/tls/kube/user-provided`     | `apiserver_key_pem_path`
`ca.crt`              | `modules/tls/kube/user-provided`     | `ca_cert_pem_path`
`ca.crt`              | `modules/tls/ingress/user-provided`  | `ca_cert_pem_path`
`etcd-ca.crt`         | `modules/tls/etcd/user-provided`     | `etcd_ca_crt_pem_path`
`etcd-client.crt`     | `modules/tls/etcd/user-provided`     | `etcd_client_crt_pem_path`
`etcd-client.key`     | `modules/tls/etcd/user-provided`     | `etcd_client_key_pem_path`
`etcd-peer.crt`       | `modules/tls/etcd/user-provided`     | `etcd_peer_crt_pem_path`
`etcd-peer.key`       | `modules/tls/etcd/user-provided`     | `etcd_peer_key_pem_path`
`etcd-server.crt`     | `modules/tls/etcd/user-provided`     | `etcd_server_crt_pem_path`
`etcd-server.key`     | `modules/tls/etcd/user-provided`     | `etcd_server_key_pem_path`
`identity-client.crt` | `modules/tls/identity/user-provided` | `client_cert_pem_path`
`identity-client.key` | `modules/tls/identity/user-provided` | `client_key_pem_path`
`identity-server.crt` | `modules/tls/identity/user-provided` | `server_cert_pem_path`
`identity-server.key` | `modules/tls/identity/user-provided` | `server_key_pem_path`
`ingress.crt`         | `modules/tls/ingress/user-provided`  | `cert_pem_path`
`ingress.key`         | `modules/tls/ingress/user-provided`  | `key_pem_path`
`kubelet.crt`         | `modules/tls/kube/user-provided`     | `kubelet_cert_pem_path`
`kubelet.key`         | `modules/tls/kube/user-provided`     | `kubelet_key_pem_path`

