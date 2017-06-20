# Troubleshooting master nodes using SSH

Tectonic master nodes are usually assigned to a public IP address. To debug a master node, SSH to it or use a VPN connected to the internal network.

View logs on the master node by using `journalctl -xe` or similar tools. [Reading the system log][journalctl] has more information.

If the cluster is deployed on AWS, check if the `init-assets` service started successfully and inspect the downloaded assets:
```sh
$ systemctl status init-assets && journalctl -u init-assets
$ ls /opt/tectonic
```

To examine if the kubelet log, execute:
```sh
$ systemctl status kubelet && journalctl -u kubelet
```

To examine the status and logs of the bootstrap and target control plane containers, execute:
```sh
$ docker ps -a | grep -v pause | grep apiserver
65faeddd2b78        quay.io/coreos/hyperkube@sha256:297f45919160ea076831cd067833ad3b64c789fcb3491016822e6f867d16dcd5                               "/usr/bin/flock /var/"   13 minutes ago      Up 13 minutes                                   k8s_kube-apiserver_kube-apiserver-90pzs_kube-system_2983ff1c-510e-11e7-bc88-063d653969e3_0
$ docker logs 65faeddd2b78
```

The `bootkube` service is responsible for bootstrapping the temporary control plane and to bootstrap a vanilla Kubernetes control plane.
To examine the `bootkube` logs, execute:
```sh
$ journalctl -u bootkube
...
Jun 14 14:31:39 ip-10-0-23-37 bash[1313]: [  219.261765] bootkube[5]:         Pod Status:        pod-checkpointer        Pending
Jun 14 14:31:39 ip-10-0-23-37 bash[1313]: [  219.262217] bootkube[5]:         Pod Status:          kube-apiserver        Running
Jun 14 14:31:39 ip-10-0-23-37 bash[1313]: [  219.262518] bootkube[5]:         Pod Status:          kube-scheduler        Pending
Jun 14 14:31:39 ip-10-0-23-37 bash[1313]: [  219.262746] bootkube[5]:         Pod Status: kube-controller-manager        Pending
...
Jun 14 14:32:44 ip-10-0-23-37 bash[1313]: [  284.264617] bootkube[5]:         Pod Status: kube-controller-manager        Running
Jun 14 14:32:49 ip-10-0-23-37 bash[1313]: [  289.263245] bootkube[5]:         Pod Status:        pod-checkpointer        Running
Jun 14 14:32:49 ip-10-0-23-37 bash[1313]: [  289.263932] bootkube[5]:         Pod Status:          kube-apiserver        Running
Jun 14 14:32:49 ip-10-0-23-37 bash[1313]: [  289.264715] bootkube[5]:         Pod Status: kube-controller-manager        Running
...
Jun 14 14:34:29 ip-10-0-23-37 bash[1313]: [  389.299380] bootkube[5]: Tearing down temporary bootstrap control plane...
Jun 14 14:34:29 ip-10-0-23-37 systemd[1]: Started Bootstrap a Kubernetes cluster.
```

The `tectonic` service is responsible for installing the actual Tectonic assets on the bootstrapped vanilla cluster.
To examine the `tectonic` installation logs, execute:
```sh
$ journalctl -fu tectonic
Jun 14 14:36:22 ip-10-0-23-37 bash[4763]: [  502.655337] hyperkube[5]: Pods not available yet, waiting for 5 seconds (10)
Jun 14 14:36:27 ip-10-0-23-37 bash[4763]: [  507.955606] hyperkube[5]: Tectonic installation is done
Jun 14 14:36:28 ip-10-0-23-37 systemd[1]: Started Bootstrap a Tectonic cluster.
```

[journalctl]: https://github.com/coreos/docs/blob/master/os/reading-the-system-log.md
