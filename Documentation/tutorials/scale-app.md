# Scaling an application on Tectonic

Kubernetes can scale applications either statically, based on an explicit count of replicas, or dynamically, adding and removing replicas based on application load.

This tutorial outlines static application scaling in Tectonic:

* Scaling an application from 3 replicas up to 5 using `kubectl` and Tectonic
* Monitoring this change in the Tectonic Console

It is assumed that you have a functioning Tectonic cluster to try these changes.

## Scaling a deployment with kubectl

Start by copying the following YAML into a file named `cookies.yaml`:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: cookies
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

Deploy this file with the following `kubectl create` command:

```sh
$ kubectl create -f cookies.yaml
deployment "cookies" created
$ kubectl get deployments/cookies
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies   3         3         3            3           18s
```

Next, use the text editor of your choice to change the `replicas` value to 5:

```sh
spec:
  replicas: 5   # bump from 3 to 5
```

Apply the changes with `kubectl apply`:

```sh
$ kubectl apply -f cookies.yaml
deployment "cookies" configured
$ kubectl get deploy/cookies
NAME      DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies   5         5         5            5           1m
```

Notice the pod count increases from 3 to 5.

### Scaling a deployment with the kubectl scale subcommand

Deployments can be scaled directly on the `kubectl` command line with the `scale` subcommand. The following command scales the `cookies` Deployment back down to 3 replicas.

```sh
$ kubectl scale deployments/cookies --replicas=3
deployment "cookies" scaled
$ kubectl get deployments/cookies
NAME                 DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
cookies              3         3         3            3           6m
```

## Scaling a deployment with Tectonic Console

To control application scale in Tectonic Console, first copy the `cookies` Deployment's YAML Manifest given above into the Console's YAML editor to create a new deployment:

1. Click *Deployments* under *Workloads*.
2. Click the *Create Deployment* button. The YAML Manifest editor is shown, with a skeletal Manifest file. Delete the example Manifest YAML.
3. Paste the `cookies` Deployment's Manifest into the editor.
4. Click *Create* to save and instantiate the `cookies` Deployment.

### Scaling a Deployment with the Console UI

You may also use the Tectonic Console to scale a deployment graphically.

1. Click *Workloads > Deployments* in the menu at left.
2. Click the Deployment to be scaled. The Deployment detail is shown.
3. Click the current number of pods beneath the Desired Count heading. The *Modify Desired Count* dialog is shown.
4. Change the number of Replicas with the plus and minus buttons, or by typing the number directly. Click *Save Desired Count*.

### Scaling by editing the Manifest in Tectonic Console

To edit a Deployment's Manifest file in the Console, click *Deployments* beneath *Workloads* in the menu at left. Click the Deployment to be edited, named *cookies* in this example. The Deployment detail is shown. Click the *YAML* tab. The Deployment Manifest is shown.

Change the `replicas` count to the desired number of replicas of the Deployment. In this example, the `cookies` Deployment is scaled from 3 replicas to 5.

```sh
spec:
  replicas: 5   # bump from 3 to 5
```

Click *Save Changes* to scale the Deployment to the new number of replicas.

## Monitoring Deployment scale in Tectonic Console

Use the *Pods* tab of a given deployment to track the number of Replicas, and the health of those Pods.

[**NEXT:** Versioning and updating an application on Tectonic][versioning-app]


[versioning-app]: rolling-deployments.md
