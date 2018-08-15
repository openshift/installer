package installerassets

func init() {
	Defaults["image/bootkube"] = ConstantDefault([]byte("quay.io/coreos/bootkube:v0.14.0"))
	Defaults["image/etcd-cert-signer"] = ConstantDefault([]byte("quay.io/coreos/kube-etcd-signer-server:678cc8e6841e2121ebfdb6e2db568fce290b67d6"))
	Defaults["image/etcdctl"] = ConstantDefault([]byte("quay.io/coreos/etcd:v3.2.14"))
	Defaults["image/release"] = ConstantDefault([]byte("registry.svc.ci.openshift.org/openshift/origin-release:v4.0"))
}
