# Disaster recovery of scheduler and controller manager pods

It is possible to recover a Kubernetes cluster from the failure of certain control plane components by first replacing failed control plane pods with manually scheduled temporary instances. These temporary instances will then schedule permanent replacements for the failed components.

If the API server is still running, a failed `kube-scheduler`, or `kube-controller-manager`, or both may be replaced with the following process.

## Recovering a scheduler

To recover a failed scheduler pod, place a new `kube-scheduler` pod into the cluster. Delete the temporary recovery pod once it has created a new `kube-scheduler`, and inspect the health of the pod to confirm.

First, use kubectl to get a copy of the manifest:
```bash
kubectl --namespace=kube-system get deployment kube-scheduler -o yaml
```

Then, extract the pod `spec` by running the following command:

`label=kube-scheduler ; namespace=kube-system ; kubectl get deploy --namespace=$namespace -l k8s-app=${label} -o json --export | jq --arg namespace $namespace --arg name ${label}-rescue --arg node $(kubectl get node -l master -o jsonpath='{.items[0].metadata.name}') '.items[0].spec.template | .kind = "Pod" | .apiVersion = "v1" | del(.metadata, .spec.nodeSelector) | .metadata.namespace = $namespace | .metadata.name = $name | .spec.containers[0].name = $name | .spec.nodeName = $node | .spec.serviceAccount = "default" | .spec.serviceAccountName = "default" ' | kubectl convert -f-`

(Or, manually copy the`spec` section under `template`.)

Executing this command creates the following rescue pod, which is used as a temporary instance of the scheduler:
```yaml
apiVersion: v1
items:
- apiVersion: v1
  kind: Pod
  metadata:
    creationTimestamp: null
    name: kube-scheduler-rescue
    namespace: kube-system
  spec:
    containers:
    - command:
      - ./hyperkube
      - scheduler
      - --leader-elect=true
      image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
      imagePullPolicy: IfNotPresent
      name: kube-scheduler-rescue
      resources: {}
      terminationMessagePath: /dev/termination-log
    dnsPolicy: ClusterFirst
    nodeName: ip-10-0-60-193.us-west-2.compute.internal
    restartPolicy: Always
    securityContext: {}
    serviceAccount: default
    serviceAccountName: default
    terminationGracePeriodSeconds: 30
  status: {}
kind: List
metadata: {}
```

Next, create a pod `spec`. For example,`recovery-pod.yaml`:
```yaml
spec:
  nodeName:
  containers:
    - command:
    - ./hyperkube
    - scheduler
    - --leader-elect=true
    image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
    imagePullPolicy: IfNotPresent
    name: kube-scheduler
```

Then, get the name of the master node:
```bash
kubectl get nodes -l master=true
```

If using AWS EC2, your master node name will be returned with the format: `ip-12-34-56-78.us-west-2.compute.internal`. Use this value as the master `nodeName` when creating your temporary scheduler pod, as described below.
Wrap the `spec` in a pod header and specify the name of the master node in `nodeName`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: kube-scheduler-rescue
spec:
  nodeName: <master-node-name-returned-above>
  containers:
    - command:
    - ./hyperkube
    - scheduler
    - --leader-elect=true
  image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
  imagePullPolicy: IfNotPresent
  name: kube-scheduler
```

Use kubectl create to inject the pod into the cluster:
```bash
kubectl create -f recovery-pod.yaml
```
This pod acts as a rescue `kube-scheduler`, and temporarily performs all the operations of a scheduler.

When your rescue scheduler has successfully created a new kube-scheduler, delete the recovery pod.
```bash
kubectl delete -f recovery-pod.yaml
```

Finally, inspect the health of the pod to make certain that all pods are up and running as expected.
```bash
kubectl get pods --all-namespaces
```

## Recovering a controller manager

The recovery process for a controller manager is essentially the same as that for a scheduler.

Place a new `kube-controller-manager` pod into the cluster. Delete the temporary recovery pod once it has created a new `kube-controller-manager`, and inspect the health of the pod to confirm.

First, place a `kube-controller-manager` pod into the cluster.

Run `kubectl` to get a copy of the manifest:

```bash
kubectl --namespace=kube-system get deployment kube-controller-manager -o yaml
```
Extract the pod `spec` by running the following command:

`kubectl --namespace=kube-system get deployment kube-controller-manager -ojson | jq '.spec.template.apiVersion = "v1" | .spec.template.kind = "Pod" | .spec.template.metadata = {"namespace": .metadata.namespace} | .spec.template.metadata.name = .metadata.name + "-recovery" | .spec.template'`

Then, create an example `spec`, `recovery-pod.yaml`:
```yaml
spec:
  nodeName:
  containers:
    - command:
    - ./hyperkube
    - controller-manager
    - --leader-elect=true
    image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
    imagePullPolicy: IfNotPresent
    name: kube-controller-manager
```
Get the name of the master node:

```bash
kubectl get nodes -l master=true
```

If using AWS EC2, your master node name will be returned with the format: `ip-12-34-56-78.us-west-2.compute.internal`. Use this value as the master `nodeName` when creating your temporary controller manager pod, as described below.

Then, wrap the pod `spec` in a pod header and specify the name of the master node in `nodeName`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: recovery-controller
spec:
  nodeName: <master-node-name-returned-above>
  containers:
    - command:
    - ./hyperkube
    - controller-manager
    - --leader-elect=true
  image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
  imagePullPolicy: IfNotPresent
  name: kube-controller-manager
```

And finally, inject the pod into the cluster:
```bash
kubectl create -f recovery-pod.yaml
```

This pod acts as a temporary `kube-controller-manager`, which would convert the existing `kube-controller-manager` into pods. These pods will then be scheduled.

When a new kube-controller-manager has been scheduled, delete the temporary recovery pod:
```bash
kubectl delete -f recovery-pod.yaml
```

Once the temporary controller manager has been removed, inspect the health of the pod to confirm that all pods are up and running as expected:
```bash
kubectl get pods --all-namespaces
```
