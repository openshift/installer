package content

var (
	// KubeletSystemdContents is a service for running the kubelet on the
	// bootstrap nodes.
	KubeletSystemdContents = `
[Unit]
Description=Kubernetes Kubelet
Wants=rpc-statd.service

[Service]
Type=notify
ExecStartPre=/bin/mkdir --parents /etc/kubernetes/manifests
ExecStartPre=/usr/bin/bash -c "gawk '/certificate-authority-data/ {print $2}' /etc/kubernetes/kubeconfig | base64 --decode > /etc/kubernetes/ca.crt"
Environment=KUBELET_RUNTIME_REQUEST_TIMEOUT=10m
EnvironmentFile=-/etc/kubernetes/kubelet-env

ExecStart=/usr/bin/hyperkube \
  kubelet \
    --container-runtime=remote \
    --container-runtime-endpoint=/var/run/crio/crio.sock \
    --runtime-request-timeout=${KUBELET_RUNTIME_REQUEST_TIMEOUT} \
    --lock-file=/var/run/lock/kubelet.lock \
    --exit-on-lock-contention \
    --pod-manifest-path=/etc/kubernetes/manifests \
    --allow-privileged \
    --minimum-container-ttl-duration=6m0s \
    --cluster-domain=cluster.local \
    --client-ca-file=/etc/kubernetes/ca.crt \
    --anonymous-auth=false \
    --cgroup-driver=systemd \
    --serialize-image-pulls=false \

Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
`
)
