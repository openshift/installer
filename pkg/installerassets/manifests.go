package installerassets

func init() {
	Rebuilders["files/opt/tectonic/manifests/cvo-overrides.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/cvo-overrides.yaml",
		map[string]string{
			"ClusterID": "cluster-id",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/etcd-service-endpoints.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/etcd-service-endpoints.yaml",
		map[string]string{
			"BaseDomain":  "base-domain",
			"ClusterName": "cluster-name",
			"MasterCount": "machines/master-count",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/host-etcd-service-endpoints.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/host-etcd-service-endpoints.yaml",
		map[string]string{
			"BaseDomain":  "base-domain",
			"ClusterName": "cluster-name",
			"MasterCount": "machines/master-count",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/legacy-cvo-overrides.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/legacy-cvo-overrides.yaml",
		map[string]string{
			"ClusterID": "cluster-id",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/kube-system-configmap-etcd-serving-ca.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/kube-system-configmap-etcd-serving-ca.yaml",
		map[string]string{
			"Cert": "tls/etcd-ca.crt",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/kube-system-configmap-root-ca.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/kube-system-configmap-root-ca.yaml",
		map[string]string{
			"Cert": "tls/root-ca.crt",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/kube-system-secret-etcd-client.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/kube-system-secret-etcd-client.yaml",
		map[string]string{
			"Cert": "tls/etcd-client.crt",
			"Key":  "tls/etcd-client.key",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/machine-config-server-tls-secret.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/machine-config-server-tls-secret.yaml",
		map[string]string{
			"Cert": "tls/machine-config-server.crt",
			"Key":  "tls/machine-config-server.key",
		},
		nil,
	)

	Rebuilders["files/opt/tectonic/manifests/openshift-service-cert-signer-ca-secret.yaml"] = TemplateRebuilder(
		"files/opt/tectonic/manifests/openshift-service-cert-signer-ca-secret.yaml",
		map[string]string{
			"Cert": "tls/service-serving-ca.crt",
			"Key":  "tls/service-serving-ca.key",
		},
		nil,
	)
}
