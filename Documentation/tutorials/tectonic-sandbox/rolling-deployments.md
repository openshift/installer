# Versioning and updating an application on Tectonic

Kubernetes uses rolling deployments to update an application between releases. When a deployment configuration is updated, Tectonic automatically starts to execute a smooth, rolling update of the application. If the update fails it can roll back to the previous instance.

This tutorial will:
* Create a service.
* Add ingress to the service.
* Deploy a sample "Cookie Shop" application v1.0.
* Update the sample application from v1.0 to v2.0.
* Perform a health check.
* Roll back to v1.0.

## Updating an app using Tectonic Console

Use Tectonic Sandbox and Console to update an app, and demonstrate the power of rolling updates.

### Create a service

If the service created using `simple-service.yaml` and `simple-ingress.yaml` from the [Deploying an application tutorial][first-app] is running, it may be reused for this tutorial, and its deployment updated. If not, create the service before continuing with this tutorial.

To create a service using Tectonic Console:

* Go to *Routing > Services* and click *Create Service*.
* Copy and paste the contents of `simple-service.yaml` described in [Deploying an application from YAML][first-yaml] in to the editor, to replace the YAML file displayed.
* Click *Create* to create the service, and open the *Service Overview* window.

### Add ingress

Then, create ingress to the service:

* Go to *Routing > Ingress* and click *Create Ingress*.
* Copy and paste the contents of `simple-ingress.yaml` described in [Deploying an application from YAML][first-yaml] in to the editor, to replace the YAML file displayed.
* Click *Create* to add ingress, and open the ingress *Overview* window.

### Deploy v1.0

When executing a rolling update, the rate at which updates will occur, and health checks to determine when an app has launched correctly may be defined.

First, edit the `simple-deployment.yaml` file used in [Deploying an application tutorial][first-app] to create a service which more explicitly defines how Kubernetes will run the app.

To edit the service:

* Set the update model to rolling update.
* Add a `readiness` and `liveness` probe which instruct Kubernetes when and how to report that the pods are alive and well.
* Specify the `restartPolicy`, `dnsPolicy`, and `TerminationGracePeriod`.

Go to *Workloads > Deployments*, and click *Create Deployment*.

Enter the following YAML into the window, and click *Create* to save changes and deploy the app.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: simple-deployment
  namespace: default
  labels:
    k8s-app: simple
spec:
  replicas: 3
  revisionHistoryLimit: 2
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  template:
    metadata:
      labels:
        k8s-app: simple
    spec:
      containers:
        - name: nginx
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

Load the service address to see the application:

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12">
    <img src="img/sandbox-cookie-v1.png">
  </div>
</div>

### Deploy v2.0

Next, update the app to change the bakery's branding on the site. Build a new container image, `v2.0`, with the updated logo and background color.

To update the app:

* Go to *Workloads > Deployments*, and click the *simple-deployment* created above.
* Click the *YAML* tab.
* Edit the file displayed to change
  `image: quay.io/coreos/example-app:v1.0`
to
  `image: quay.io/coreos/example-app:v2.0`

The edited file will look like this:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: simple-deployment
  namespace: default
  labels:
    k8s-app: simple
spec:
  replicas: 3
  revisionHistoryLimit: 2
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  template:
    metadata:
      labels:
        k8s-app: simple
    spec:
      containers:
        - name: nginx
          image: quay.io/coreos/example-app:v2.0
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

Click *Save Changes* to update the app, and begin live updates.

Click the *Pods* tab to see existing pods, and watch the live update in progress. Tectonic Console provides live updates as newer pods appear, and the corresponding older pods disappear.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/sandbox-simple-deploy-updating.png">
    <div class="co-m-screenshot-caption">Pods being created and terminated by Tectonic</div>
  </div>
</div>

Reload the Cookie Shop page as the deployment updates to see the changes.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12">
    <img src="img/sandbox-cookie-v2.png">
  </div>
</div>

### Perform a health check

The deployment specifies a readiness probe and a liveness probe to determine when it is safe to send traffic to each instance of the application.

If the **readiness** probe fails, the rolling update will be halted automatically. Traffic will never be sent to this pod.

If the **liveness** probe fails at any time over the life of the pod, traffic will be shifted away from the pod.

Between these two mechanisms, Tectonic is always informed of the state of the application and can act accordingly.

### Roll back to v1.0

Use the Console to make quick changes to the deployment. Click the *YAML* tab for *simple-deployment* to edit the manifest, and change `v2.0` back to `v1.0` to roll back the change.

After saving the file, the pods will execute a rolling update back to version 1.0.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/sandbox-simple-deploy-yml.png">
    <div class="co-m-screenshot-caption">Using the YAML editor to edit the deployment</div>
  </div>
</div>

## Updating an app using kubectl

The same process may be followed using kubectl and the command line.

### Create a service and add ingress using kubectl

```sh
$ kubectl create -f simple-service.yaml
service "simple-service" created
$ kubectl get services -o wide
NAME              CLUSTER-IP   EXTERNAL-IP    PORT(S)        AGE    SELECTOR
simple-service    10.3.113.190 <pending>      80:30657/TCP   1d     k8s-app=simple
```

Then, create the ingress on the cluster:
```sh
$ kubectl create -f simple-ingress.yaml
ingress "simple-deployment" created
$ kubectl get ingress
NAME                HOSTS                               ADDRESS   PORTS     AGE
simple-deployment   console.tectonicsandbox.com             80        24s
```

### Deploy v1.0 using kubectl

First, edit the `simple-deployment.yaml` file used in [Deploying an application tutorial][first-app], as described above in [Deploy v1.0](#deploy-v10) with Tectonic Console.

Then, use `kubectl apply` to deploy the app. Kubernetes will manage the migration from the old version to the new version automatically.

```sh
$ kubectl apply -f simple-deployment.yaml
deployment "simple-deployment" configured
```

Load the Service address to see the application:

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12">
    <img src="img/sandbox-cookie-v1.png">
  </div>
</div>

### Deploy v2.0 using kubectl

Next, update the app to change the bakery's branding on the site. Build a new container image, `v2.0`, with the updated logo and background color.

To update the app, first edit the `simple-deployment.yaml`file, as described above, to change
  `image: quay.io/coreos/example-app:v1.0`
to
  `image: quay.io/coreos/example-app:v2.0`

Then, open the deployment in Tectonic Console. The Console shows live updates, refreshing as listed pods are created and destroyed by the deployment. Click the *Pods* tab to see existing pods.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/sandbox-simple-service.png">
    <div class="co-m-screenshot-caption">Viewing the pods of your deployment</div>
  </div>
</div>

Next, switch back to a terminal and apply the change:

```sh
$ kubectl apply -f simple-deployment.yaml
deployment "simple-deployment" configured
```

Tectonic Console provides live updates as newer pods appear, and the corresponding older pods disappear.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/sandbox-simple-deploy-updating.png">
    <div class="co-m-screenshot-caption">Pods being created and terminated by Tectonic</div>
  </div>
</div>

Reload the Cookie Shop page as the deployment updates to see changes.

Once the app is updated, use Tectonic Console to [perform a health check](#perform-a-health-check) and [roll back to v1.0](#roll-back-to-v10), as described above.


[**NEXT:** Advanced metrics monitoring with Prometheus][prometheus]


[first-app]: first-app.md#deploy-a-simple-app
[first-yaml]: first-app.md#deploying-an-application-from-yaml
[k8s-deployments]: https://kubernetes.io/docs/user-guide/deployments/
[k8s-rolling-updates]: https://kubernetes.io/docs/user-guide/rolling-updates/
[prometheus]: monitoring.md
