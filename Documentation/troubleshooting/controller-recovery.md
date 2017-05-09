# Disaster recovery of Scheduler and Controller Manager pods

It is possible to recover a Kubernetes cluster from the failure of certain control plane components. Provided the API server is still running, a failed `kube-scheduler` or `kube-controller-manager`, or both the pods can be replaced with the process described in this document. At a high level, these failed control plane pods are replaced with manually scheduled temporary instances. These temporary instances in turn schedule permanent replacements for the failed control plane components.

## Recovering a Scheduler

1. Place a `kube-scheduler` pod into the cluster:

    a. Run the following to get a copy of the manifest:

        kubectl --namespace=kube-system get deployment kube-scheduler -o yaml

    b. Extract the pod `spec` by running the following command:

      ```label=kube-scheduler ; namespace=kube-system ; kubectl get deploy --namespace=$namespace -l k8s-app=${label} -o json --export | jq --arg namespace $namespace --arg name ${label}-rescue --arg node $(kubectl get node -l master -o jsonpath='{.items[0].metadata.name}') '.items[0].spec.template | .kind = "Pod" | .apiVersion = "v1" | del(.metadata, .spec.nodeSelector) | .metadata.namespace = $namespace | .metadata.name = $name | .spec.containers[0].name = $name | .spec.nodeName = $node | .spec.serviceAccount = "default" | .spec.serviceAccountName = "default" ' | kubectl convert -f-```

     Alternatively, the`spec` section under `template` can also be manually copied.

     Executing this command gives the following rescue pod:

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

    c. Create a pod `spec`. For example,`recovery-pod.yaml`:

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


   d. Get the name of the master node:

     kubectl get nodes -l master=true

   e. Wrap up the `spec` in a pod header and specify the name of the master node in `nodeName`:

			   apiVersion: v1
			   kind: Pod
			   metadata:
			     name: kube-scheduler-rescue
			   spec:
			      nodeName: <name-of-the-tectonic-master-node>
			   containers:
			    - command:
			    - ./hyperkube
			    - scheduler
			    - --leader-elect=true
			    image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
			    imagePullPolicy: IfNotPresent
			    name: kube-scheduler


   f. Inject the pod into the cluster:

		  kubectl create -f recovery-pod.yaml

   This pod acts as a rescue `kube-scheduler`, and  temporarily performs all the operations of a scheduler.

2. Delete the recovery pod once the existing `kube-scheduler` has been scheduled:

  `kubectl delete -f recovery-pod.yaml`

3. Inspect the health of the pod:

  `kubectl get pods --all-namespaces`

  Make sure that all pods are up and running.

## Recovering a Controller Manager

1. Place a `kube-controller-manager` pod into the cluster:

    a. Run the following to get a copy of the manifest:

        kubectl --namespace=kube-system get deployment kube-controller-manager -o yaml

    b. Extract the pod `spec` by running the following command:

    ```kubectl --namespace=kube-system get deployment kube-controller-manager -ojson | jq '.spec.template.apiVersion = "v1" | .spec.template.kind = "Pod" | .spec.template.metadata = {"namespace": .metadata.namespace} | .spec.template.metadata.name = .metadata.name + "-recovery" | .spec.template'```

    c. Create an example `spec`, `recovery-pod.yaml`:

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

   d. Get the name of the master node:

        kubectl get nodes -l master=true

   e. Wrap up the pod `spec` in a pod header and specify the name of the master node in `nodeName`:

			   apiVersion: v1
			   kind: Pod
			   metadata:
			      name: recovery-controller
			   spec:
			      nodeName: <name-of-the-tectonic-master-node>
			   containers:
			    - command:
			    - ./hyperkube
			    - controller-manager
			    - --leader-elect=true
			    image: quay.io/coreos/hyperkube:v1.6.2_coreos.0
			    imagePullPolicy: IfNotPresent
			    name: kube-controller-manager

   f. Inject the pod into the cluster:

        kubectl create -f recovery-pod.yaml

   This pod acts as a temporary `kube-controller-manager`, which would convert the existing `kube-controller-manager` into pods. These pods will then be scheduled.

2. Delete the recovery pod once the existing `kube-controller-manager` has been scheduled:

  `kubectl delete -f recovery-pod.yaml`

3. Inspect the health of the pod:

  `kubectl get pods --all-namespaces`

  Make sure that all pods are up and running.

