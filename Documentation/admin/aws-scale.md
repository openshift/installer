# Scaling Tectonic AWS clusters

This document explains how to add and remove worker and controller nodes to scale a Tectonic cluster, and describes connecting to either an existing external etcd cluster, or one created by Tectonic installer.

## Scaling Kubernetes nodes

During installation, Tectonic creates two AWS Auto Scaling Groups (ASGs). The `AutoScaleWorker` group contains worker nodes, while the `AutoScaleController` group contains cluster controller nodes. The `AutoScaleWorker` group can be scaled up or down to add or remove cluster work capacity. The `AutoScaleController` group can be scaled up to increase control plane availability.

Set the number of members of either group by visiting the [AWS EC2 console][aws-ec2], under *Auto Scaling* in the left hand sidebar. Edit the *min*, *max* and *desired* fields for the respective ASG to match capacity needs for workers, or availability demands for more than one controller.

After increasing the number of nodes in either ASG, new nodes will boot and join the cluster within a few minutes.

Scaling down will immediately terminate the node VMs, causing replicated workloads to be restarted on other nodes automatically. The Kubernetes node draining facility is not employed when removing nodes from an ASG.

The Console shows both types of nodes, with the control plane denoted with the label `master=true`.

![Scaled Nodes][scaled-nodes]

## External etcd cluster

Tectonic Installer creates an external etcd cluster by default. Set the size of this cluster greater than the default 3 nodes in the case of high demand. Alternatively, Tectonic Installer will accept the client endpoint of an existing etcd v3 cluster to which it should connect instead.

![Tectonic installer: etcd cluster choice][byo-etcd]

[byo-etcd]: ../img/byo-etcd.png
[scaled-nodes]: ../img/scaled-nodes.png
[aws-autoscaling]: http://docs.aws.amazon.com/autoscaling/latest/userguide/WhatIsAutoScaling.html
[aws-ec2]: https://console.aws.amazon.com/ec2/
