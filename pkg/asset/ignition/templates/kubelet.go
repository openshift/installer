package templates

const (
	// KubeletSystemdContents is a service for running the kubelet on the
	// bootstrap nodes.
	KubeletSystemdContents = `[Unit]
Description=Kubernetes Kubelet
Wants=rpc-statd.service

[Service]
ExecStartPre=/bin/mkdir --parents /etc/kubernetes/manifests
ExecStartPre=/usr/bin/bash -c "gawk '/certificate-authority-data/ {print $2}' /etc/kubernetes/kubeconfig | base64 --decode > /etc/kubernetes/ca.crt"

ExecStart=/usr/bin/hyperkube \
  kubelet \
    --bootstrap-kubeconfig=/etc/kubernetes/kubeconfig \
    --kubeconfig=/var/lib/kubelet/kubeconfig \
    --rotate-certificates \
    --cni-conf-dir=/etc/kubernetes/cni/net.d \
    --cni-bin-dir=/var/lib/cni/bin \
    --network-plugin=cni \
    --lock-file=/var/run/lock/kubelet.lock \
    --exit-on-lock-contention \
    --pod-manifest-path=/etc/kubernetes/manifests \
    --allow-privileged \
    --node-labels=node-role.kubernetes.io/bootstrap \
    --register-with-taints=node-role.kubernetes.io/bootstrap=:NoSchedule \
    --minimum-container-ttl-duration=6m0s \
    --cluster-dns={{.ClusterDNSIP}} \
    --cluster-domain=cluster.local \
    --client-ca-file=/etc/kubernetes/ca.crt \
    --cloud-provider={{.CloudProvider}} \
    --anonymous-auth=false \
    --cgroup-driver=systemd \
    {{.CloudProviderConfig}} \
    {{.DebugConfig}} \

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target`
)
