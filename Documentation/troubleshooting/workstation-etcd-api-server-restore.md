# Running a Kubernetes API server from an etcd snapshot

Running a Kubernetes API server on a local workstation based on an etcd snapshot can be useful for debugging, and for understanding database backups. In a few simple steps this can be accomplished on Windows, Linux, or OSX.

First, locate the `member/snap/db` file, which is a backup of the etcd data directory. This file is usually stored in `/var/etcd`, `/var/lib/etcd`, or a backup server. Then, use this file to generate a new data-dir that can be run locally.

## Prerequisites

- etcd and etcdctl for your workstation. Find a [release here](https://github.com/coreos/etcd/releases).
- kube-apiserver for your workstation. Go get `k8s.io/kubernetes/cmd/kube-apiserver`.

## Setup

First, use etcdctl to create a copy of your data directory.

```
$ ETCDCTL_API=3 etcdctl snapshot restore kube-system-kube-etcd-0000/member/snap/db --name m1 --initial-cluster m1=http://localhost:2380 --initial-cluster-token etcd-cluster-1 --initial-advertise-peer-urls http://localhost:2380 --skip-hash-check
```

Now, run etcd against the generated data directory.

```
etcd --data-dir m1.etcd/
```

Finally, launch an API server that hits this etcd server running on localhost.

```
kube-apiserver --cert-dir ./kubernetes-temp --etcd-servers=http://localhost:2379 --service-cluster-ip-range 10.3.0.0/12
```

Now, `curl http://localhost:8080/api/v1/events`. Done!
