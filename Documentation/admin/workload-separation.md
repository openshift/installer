# Workload separation

Tectonic 1.6 and later leverages Kubernetes' [new __beta__][doc-taints] support 
for tolerations, taints and pod affinity, to guarantee a clear separation 
between control plane and user workloads, as well as improve the reliability 
of the control plane by spreading the services across multiple nodes.

Every Tectonic _master_ node registers itself to the cluster tainted.
Taints are an attribute or marker applied to a node that describe a special 
property. Before a pod can be scheduled to this node, it must tolerate the 
taint. Since most pods will not tolerate this taint, they will be scheduled only
to _worker_ nodes. This behavior also includes pods generated from Deployments, 
ReplicaSets, and [DaemonSets][ds-tolerations].

Components of the control plane, however, are specifically marked to tolerate
the taint and therefore are allowed to be scheduled on _master_ nodes. 
Furthermore, most of them are constrained to exclusively run on these nodes and
thus won't be collocated with any other workloads (see below).

Finally, in order to limit the impact of host failures, the Kubernetes scheduler
is hinted, using the pod anti-affinity mechanism, to avoid running multiple 
instances of the same components on the same nodes. While not strictly related 
to workloads separation, disruption control requires the cluster to always have
at least one instance of the component running at any given time, regardless of
update / re-scheduling operations.

### Control plane's components scheduling

The following table provides a reference for the scheduling of the control 
plane's components, as of Tectonic 1.6.2.

|            Component           |    Type    | Replicas (default) | Node role | Anti-Affinity | Disruption Control |
|:------------------------------:|:----------:|:------------------:|:---------:|:-------------:|:------------------:|
|         kube-apiserver         |  DaemonSet |          n         |   master  |      N/A      | N/A                |
|     kube-controller-manager    | Deployment |          2         |   master  |      yes      | yes                |
|         kube-scheduler         | Deployment |          2         |   master  |      yes      | yes                |
|            kube-dns            | Deployment |          1         |   master  |      N/A      | no                 |
|          kube-flannel          |  DaemonSet |          n         |    all    |      N/A      | N/A                |
|           kube-proxy           |  DaemonSet |          n         |    all    |      N/A      | N/A                |
| kube-etcd-network-checkpointer |  DaemonSet |          n         |   master  |      N/A      | N/A                |
|        pod-checkpointer        |  DaemonSet |          n         |   master  |      N/A      | N/A                |
|          etcd-operator         | Deployment |          1         |    any    |      N/A      | no                 |
|            kube-etcd-*         |    Pod*    |         1*         |   master  |      yes      | N/A                |
|            heapster            | ReplicaSet |          1         |    any    |      N/A      | N/A                |
|       All other workloads      |     Any    |          ?         |   worker  |       ?       | ?                  |

_n: number of qualified nodes_
_*: managed by an operator_

[doc-taints]: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#taints-and-tolerations-beta-feature
[ds-tolerations]: https://github.com/kubernetes/kubernetes/pull/41172
