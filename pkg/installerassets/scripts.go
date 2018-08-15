package installerassets

func init() {
	Rebuilders["files/usr/local/bin/bootkube.sh"] = TemplateRebuilder(
		"files/usr/local/bin/bootkube.sh",
		map[string]string{
			"BaseDomain":          "base-domain",
			"BootkubeImage":       "image/bootkube",
			"ClusterName":         "cluster-name",
			"EtcdCertSignerImage": "image/etcd-cert-signer",
			"EtcdCtlImage":        "image/etcdctl",
			"ReleaseImage":        "image/release",
			"MasterCount":         "machines/master-count",
		},
		nil,
	)
}
