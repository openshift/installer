# Tectonic with CentOS 7 Worker

## Disclaimer

This is unofficial. The author does not update, maintain, or support this setup. Please direct all questions to your support provider.

## CoreOS Tectonic

Provision a Tectonic `1.6.4-tectonic.1` bare-metal cluster (Container Linux, 1 controller, 2 workers) in the usual way with [matchbox](https://github.com/coreos/matchbox) and the Tectonic [Installer](https://coreos.com/tectonic/docs/latest/install/bare-metal/index.html).

Locally, you may PXE boot QEMU/KVM nodes via Tectonic and matchbox on the `metal0` bridge.

## CentOS 7

Provision a CentOS 7 machine on the same network. See [Downloads](https://www.centos.org/download/). Install essentials and add SSH keys.

```
sudo yum install -y openssh-server wget vim-minimal vim
sudo systemctl start sshd
sudo systemctl enable sshd
```

### rkt

Install rkt from the RPM.

```
gpg2 --recv-key 18AD5014C99EF7E3BA5F6CE950BDD3E0FC8A365E
wget https://github.com/coreos/rkt/releases/download/v1.23.0/rkt-1.23.0-1.x86_64.rpm
wget https://github.com/coreos/rkt/releases/download/v1.23.0/rkt-1.23.0-1.x86_64.rpm.asc
gpg2 --verify rkt-1.23.0-1.x86_64.rpm.asc
sudo rpm -Uvh rkt-1.23.0-1.x86_64.rpm
```

Add the rkt group and add yourself to it.

```
sudo gpasswd -a $(whoami) rkt
```

Set SELinux to Permissive mode persistently. Edit /etc/selinux/config to ensure this change is persistent.

```
$ getenforce
Permissive
```

Verify a container image can be run.

```
sudo rkt run --insecure-option=image --interactive docker://alpine --exec /bin/sh
```

### Docker

Install Docker

    sudo yum install -y docker
    sudo systemctl start docker
    sudo systemctl enable docker

### Kubelet

Add the [kubelet-wrapper](https://github.com/coreos/coreos-overlay/blob/master/app-admin/kubelet-wrapper/files/kubelet-wrapper) script to `/opt/bin/kubelet-wrapper`.

    sudo mkdir /opt/bin
    sudo vim /opt/bin/kubelet-wrapper
    sudo chmod +x /opt/bin/kubelet-wrapper

Edit the rkt `--mount` sections.

* /usr/share/ca-certificates to /etc/pki/tls/certs
* /usr/lib/os-release to /etc/os-release

Assign the machine a static IP in your DHCP server and a stable, convenience DNS name in your DNS server. This is the name the kubelet will use to register itself in the cluster. Use the name in your `kubelet.service` file below.

    $ hostname
    node4.example.com

Add `/etc/systemd/system/kubelet.service`.

```
[Unit]
Description=Kubelet via Hyperkube ACI
[Service]
Environment="RKT_OPTS=--uuid-file-save=/var/run/kubelet-pod.uuid \
  --volume=resolv,kind=host,source=/etc/resolv.conf \
  --mount volume=resolv,target=/etc/resolv.conf \
  --volume var-lib-cni,kind=host,source=/var/lib/cni \
  --mount volume=var-lib-cni,target=/var/lib/cni \
  --volume var-log,kind=host,source=/var/log \
  --mount volume=var-log,target=/var/log"
EnvironmentFile=/etc/kubernetes/kubelet.env
ExecStartPre=/bin/mkdir -p /etc/kubernetes/manifests
ExecStartPre=/bin/mkdir -p /srv/kubernetes/manifests
ExecStartPre=/bin/mkdir -p /etc/kubernetes/checkpoint-secrets
ExecStartPre=/bin/mkdir -p /etc/kubernetes/cni/net.d
ExecStartPre=/bin/mkdir -p /var/lib/cni
ExecStartPre=-/usr/bin/rkt rm --uuid-file=/var/run/kubelet-pod.uuid
ExecStart=/opt/bin/kubelet-wrapper \
  --kubeconfig=/etc/kubernetes/kubeconfig \
  --require-kubeconfig \
  --client-ca-file=/etc/kubernetes/ca.crt \
  --anonymous-auth=false \
  --cni-conf-dir=/etc/kubernetes/cni/net.d \
  --network-plugin=cni \
  --lock-file=/var/run/lock/kubelet.lock \
  --exit-on-lock-contention \
  --pod-manifest-path=/etc/kubernetes/manifests \
  --allow-privileged \
  --hostname-override=${TODO_ACTUAL_HOSTNAME} \
  --minimum-container-ttl-duration=6m0s \
  --cluster_dns=10.3.0.10 \
  --cluster_domain=cluster.local
ExecStop=-/usr/bin/rkt stop --uuid-file=/var/run/kubelet-pod.uuid
Restart=always
RestartSec=10
[Install]
WantedBy=multi-user.target
```

Add a `/etc/kubernetes/kubelet.env` file.

```
    sudo mkdir -p /etc/kubernetes
    sudo bash -c 'cat > /etc/kubernetes/kubelet.env << EOF
KUBELET_ACI=quay.io/coreos/hyperkube
KUBELET_VERSION=v1.6.2_coreos.0
EOF'
```

Copy the `ca.crt` and `kubeconfig` from an existing Container Linux worker.

    sudo vim /etc/kubernetes/ca.crt
    sudo vim /etc/kubenretes/kubeconfig

Start the Kubelet and inspect logs.

    sudo systemctl start kubelet
    sudo systemctl enable kubelet

Verify the Kubelet downloads the hyperkube ACI (its big) and is running.

    sudo journalctl -u kubelet -f

Verify the kubelet registers with an existing kube-apiserver.

    $ kubectl get nodes
    NAME                    STATUS    AGE
    node1.example.com       Ready     5h
    node2.example.com       Ready     5h
    node3.example.com       Ready     5h
    node4.example.com       Ready     2m

Verify all daemonset pods are Running.

```
kubectl get pods --all-namespaces
NAMESPACE         NAME                                       READY     STATUS    RESTARTS   AGE
kube-system       checkpoint-installer-0rk8n                 1/1       Running   0          5h
kube-system       heapster-4141187441-pm3c9                  2/2       Running   2          5h
kube-system       kube-apiserver-k7spt                       1/1       Running   3          5h
kube-system       kube-controller-manager-1472766980-5qc5w   1/1       Running   0          5h
kube-system       kube-controller-manager-1472766980-r1bhv   1/1       Running   0          5h
kube-system       kube-dns-4101612645-63wfm                  4/4       Running   4          5h
kube-system       kube-flannel-3n7pp                         2/2       Running   1          3m
kube-system       kube-flannel-brs8g                         2/2       Running   3          5h
kube-system       kube-flannel-ph0zn                         2/2       Running   3          5h
kube-system       kube-flannel-slt9s                         2/2       Running   0          5h
kube-system       kube-proxy-c887z                           1/1       Running   0          3m
kube-system       kube-proxy-dpspm                           1/1       Running   1          5h
kube-system       kube-proxy-pkft7                           1/1       Running   0          5h
kube-system       kube-proxy-zxwms                           1/1       Running   1          5h
kube-system       kube-scheduler-3027616201-2gkkw            1/1       Running   1          5h
kube-system       kube-scheduler-3027616201-xc4x7            1/1       Running   0          5h
kube-system       pod-checkpointer-node1.example.com         1/1       Running   0          5h
tectonic-system   default-http-backend-1636328216-d6q29      1/1       Running   1          5h
tectonic-system   node-exporter-0x7vn                        1/1       Running   0          3m
tectonic-system   node-exporter-gx39l                        1/1       Running   1          5h
tectonic-system   node-exporter-r8m51                        1/1       Running   0          5h
tectonic-system   node-exporter-x9rkt                        1/1       Running   1          5h
tectonic-system   prometheus-k8s-0                           2/2       Running   2          5h
tectonic-system   prometheus-operator-3833970704-9k5rj       1/1       Running   1          5h
tectonic-system   tectonic-console-3283709919-gr63q          1/1       Running   2          5h
tectonic-system   tectonic-console-3283709919-r22h2          1/1       Running   2          5h
tectonic-system   tectonic-identity-4277829508-8m139         1/1       Running   1          5h
tectonic-system   tectonic-ingress-controller-rsq3g          1/1       Running   2          5h
tectonic-system   tectonic-ingress-controller-z2cv8          1/1       Running   3          3m
tectonic-system   tectonic-ingress-controller-zpmkc          1/1       Running   2          5h
tectonic-system   tectonic-stats-emitter-1843093653-sw7cp    2/2       Running   2          5h
```

### Tectonic Console

https://storage.googleapis.com/dghubble/tectonic-worker.png

### Optional / Other

* CentOS doesn't auto-update, no need for locksmith or etcd3 gateway
* CentOS doesn't support Ignition so you'll need to provision with some other system
* Raise max user watches to 16184 to match Container Linux workers
* Adjust Kubelet `--cluster_dns` flag if using custom service CIDR
