# Scaling an application on Tectonic

Kubernetes can scale applications either statically, based on an explicit count of desired replicas, or [dynamically, adding and removing replicas based on application load][horizontal-autoscale].

This document outlines static application scaling in Tectonic.

* Scaling an application from 3 replicas up to 5 using `kubectl` and Tectonic Console.
* Monitoring this change in the Tectonic Console.

This guide continues the previous tutorial, in which an [app was installed][first-app] on a Tectonic cluster.

## Scaling a deployment with Tectonic Console

To control application scale in Tectonic Console, first create a new `cookies` Deployment that defines a number of `replicas`for the deployment.

1. Click *Workloads > Deployments*.
2. Click *Create Deployment* to open the YAML manifest editor, with a skeletal manifest file. Delete the example YAML, and paste the the `cookies` Deployment manifest below into the editor.
3. Click *Create* to save and instantiate the `cookies` Deployment.

### cookies Deployment manifest:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: cookies
  namespace: default
  labels:
    k8s-app: simple
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  template:
    metadata:
      labels:
        k8s-app: cookies
    spec:
      containers:
        - name: cookies-container
          image: quay.io/coreos/example-app:v1.0
          ports:
            - name: http
              containerPort: 80
          readinessProbe:
           httpGet:
             path: /
             port: 80
             scheme: HTTP
          livenessProbe:
           initialDelaySeconds: 30
           timeoutSeconds: 1
           httpGet:
             path: /
             port: 80
             scheme: HTTP
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      terminationGracePeriodSeconds: 30
```

Then, use the Console to scale the deployment.

1. Click *Workloads > Deployments*.
2. Click the *cookies* deployment to open its *Overview* page.
3. Under *Desired Count*, click the current number of pods to open the *Modify Desired Count* dialog.
4. Enter a new number, and click *Save Desired Count* to change the number of healthy pods the deployment will maintain.

Use Tectonic Console to edit the manifest and scale the deployment.

1. Click *Workloads > Deployments*.
2. Click the *cookies* deployment to open its *Overview* page.
3. Click the *YAML* tab to open an editable deployment manifest.

Change the `replicas` count to the desired number of pods in the deployment. In this example, the `cookies` deployment is scaled from 3 pods to 5.

```sh
spec:
  replicas: 5   # bump from 3 to 5
```

Click *Save Changes* to scale the deployment to the new number of pods.

## Monitoring deployment scale in Tectonic Console

Use the *Pods* tab of a given deployment to track the number of replicas and their health.

## Scaling a deployment with kubectl

First, copy the `cookies` Deployment manifest above into a file named `cookies.yaml`.

If you compare this Deployment with the `simple-deployment` in the [first app tutorial][first-app], you'll notice a few additions. These additional lines set up health checks and rolling update features.

Next, use `kubectl create` to deploy `cookies.yaml`:

```sh
$ kubectl create -f cookies.yaml
deployment "cookies" created
$ kubectl get deployments/cookies
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies   3         3         3            3           18s
```

Then, open `cookies.yaml` in a text editor, and change the `replicas` value to 5:

```sh
spec:
  replicas: 5   # bump from 3 to 5
```

Use `kubectl apply` to apply the changes:

```sh
$ kubectl apply -f cookies.yaml
deployment "cookies" configured
$ kubectl get deploy/cookies
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies   5         5         5            5           1m
```

Notice the pod count increases from 3 to 5.

### Scaling a deployment with the kubectl scale subcommand

Deployments can be also be scaled directly on the `kubectl` command line with the `scale` subcommand. The following command scales the `cookies` deployment back down to three replicas.

```sh
$ kubectl scale deployments/cookies --replicas=3
deployment "cookies" scaled
$ kubectl get deployments/cookies
NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies              3         3         3            3           6m
```

[**NEXT:** Versioning and updating an application on Tectonic][versioning-app]


[first-app]: first-app.md
[horizontal-autoscale]: https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/
[installing]: install.md
[versioning-app]: rolling-deployments.md
