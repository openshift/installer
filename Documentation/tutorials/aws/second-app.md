# Deploying another application on Tectonic

This guide outlines deploying a more complicated application on our Tectonic cluster, the Kubernetes Guestbook application, which has three moving parts: a frontend, and two Redis deployments. The guide outlines deploying the app using `kubectl` and the Tectonic Console, and monitoring the app with the Tectonic Console.

This guide expects you to be familiar with the material covered in [Deploying an application on Tectonic][first-app] and to have a working Tectonic cluster on AWS.

## Using kubectl

First, we'll deploy the [Guestbook app][guestbook-upstream], using the `kubectl` CLI and the Tectonic Console.

We suggest using [Quay.io][quay-io] or [Quay Enterprise][QE] to host custom container images. The Quay container registry offers sophisticated access controls, easy automated builds, and automated security scanning, free for public projects!

First create a directory called `guestbook`, which may exist anywhere on your system. Then, copy the following YAML files to that directory.

```sh
$ mkdir guestbook
```

Create the following six YAML files, saved to the `guestbook` directory.

`guestbook/redis-master-deployment.yaml`

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: redis-master
  namespace: default
spec:
  replicas: 1
  template:
    metadata:
      labels:
        k8s-app: redis
        role: master
        tier: backend
    spec:
      containers:
      - name: master
        image: redis
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        ports:
        - containerPort: 6379
```

`guestbook/redis-master-service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-master
  namespace: default
  labels:
    k8s-app: redis
    tier: backend
    role: master
spec:
  ports:
  - port: 6379
    targetPort: 6379
  selector:
    k8s-app: redis
    tier: backend
    role: master
```

`guestbook/redis-slave-deployment.yaml`:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: redis-slave
  namespace: default
spec:
  replicas: 2
  template:
    metadata:
      labels:
        k8s-app: redis
        role: slave
        tier: backend
    spec:
      containers:
      - name: slave
        image: gcr.io/google_samples/gb-redisslave:v1
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        # Take note of this env section!
        env:
        - name: GET_HOSTS_FROM
          value: dns
        ports:
        - containerPort: 6379
```

`guestbook/redis-slave-service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-slave
  namespace: default
  labels:
    k8s-app: redis
    tier: backend
    role: slave
spec:
  ports:
  - port: 6379
  selector:
    k8s-app: redis
    tier: backend
    role: slave
```

`guestbook/frontend-deployment.yaml`:

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: frontend
  namespace: default
spec:
  replicas: 3
  template:
    metadata:
      labels:
        k8s-app: guestbook
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: gcr.io/google-samples/gb-frontend:v4
        resources:
          requests:
            cpu: 100m
            memory: 100Mi
        env:
        - name: GET_HOSTS_FROM
          value: dns
        ports:
        - containerPort: 80
```

`guestbook/frontend-service.yaml`:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend
  namespace: default
  labels:
    k8s-app: guestbook
    tier: frontend
spec:
  type: LoadBalancer
  ports:
  - port: 80
  selector:
    k8s-app: guestbook
    tier: frontend
```

Launch command:

```sh
$ kubectl create -f guestbook/
service "frontend" created
deployment "frontend" created
service "redis-master" created
deployment "redis-master" created
service "redis-slave" created
deployment "redis-slave" created
$ kubectl get deploy/frontend svc/frontend -o wide
NAME           CLUSTER-IP   EXTERNAL-IP                                                             PORT(S)        AGE       SELECTOR
svc/frontend   10.3.0.175   aaebd8247ef2311e6a045021d1620193-54019671.us-west-2.elb.amazonaws.com   80:31020/TCP   1m        k8s-app=guestbook,tier=frontend

NAME              DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/frontend   3         3         3            3           1m
```

In a browser, navigate to the `EXTERNAL-IP` address:

![Frontend Works!][frontend]

## Using the Tectonic Console

The Guestbook app can also be deployed using the Tectonic Console.

If you have already deployed the Guestbook app, either delete the deployment using the `kubectl` CLI, or copy your files, and rename them.

Delete your app from the terminal:

```sh
$ kubectl delete deploy/frontend svc/frontend deploy/redis-master svc/redis-master deploy/redis-slave svc/redis-slave
service "frontend" deleted
deployment "frontend" deleted
service "redis-master" deleted
deployment "redis-master" deleted
service "redis-slave" deleted
deployment "redis-slave" deleted
```

First, create the redis-master deployment by copying and pasting your YAML files into the Tectonic Console.

Go to *Workloads > Deployments*, and click *Create Deployment* to copy and paste your three `deployment.yaml` files into the Console: `frontend-deployment.yaml`, `redis-master-deployment.yaml`, and `redis-slave-deployment.yaml`.

Then, go to *Routing > Services* and add the Service files: `frontend-service.yaml`, `redis-master-service.yaml`, and `redis-slave-service.yaml`.

Finally, go to the URL listed on the `frontend` service's Tectonic Console page to see your Guestbook.

![Frontend Works!][frontend]

## Further reading

Now that you have deployed an application on your Tectonic cluster you may find these guides useful:

* [Deploying a simple Nginx in Tectonic][first-app]
* [Introduction to Rolling Deployments][rolling-deployments]
* [Managing Namespaces in Tectonic][namespaces]
* [Managing Pull Secrets in Tectonic][pull-secrets]

[frontend]: img/frontend.png
[quay-io]: https://quay.io
[QE]: https://coreos.com/quay-enterprise/
[first-app]: first-app.md
[namespaces]: ../../admin/manage-namespaces.md
[pull-secrets]: ../../admin/manage-pull-secrets.md
[guestbook-upstream]: https://github.com/kubernetes/kubernetes/tree/master/examples/guestbook#guestbook-example
[rolling-deployments]: rolling-deployments.md
