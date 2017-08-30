# Deploying an application on Tectonic

When [installation][installing] is complete, log in to your Tectonic Console to set up cluster credentials, and deploy a simple application.

In this tutorial you will:

* Set up credentials to work with a new Tectonic Kubernetes cluster.
* Deploy a simple application with `kubectl`.
* Deploy a simple application with the Tectonic Console.

Applications are usually deployed on Tectonic clusters by passing a YAML manifest file to the `kubectl create` CLI tool. Once the application is deployed, it can be monitored and scaled using the Tectonic web interface.

## Configuring credentials

To configure credentials for your cluster, populate a `kubeconfig` file with valid authentication credentials, then configure `kubectl` to use them to connect to a Tectonic Cluster.

First, log in to Tectonic Console to authenticate.

1. Log into the Tectonic Console.
2. Click *username > My Account* on the bottom left of the page.
3. Click KUBECTL: *Download Configuration*, and follow the onscreen instructions to authenticate.
4. When the *Set Up kubectl* window opens, click *Verify Identity*.
5. Enter your username and password, and click *Login*.
6. From the Login Successful screen, copy the alphanumeric string.
7. Switch back to the Tectonic Console tab, and enter the verification string in the field provided.

Then, download the `kubectl-config` and `kubectl` files.

1. Click *Generate Configuration* to open the *Download kubectl Configuration* window, and follow the instructions provided.
2. Click *Download Configuration* to download your `kubectl-config` file.
3. Click *Mac* or *Linux* to download the `kubectl` binary for your operating system.

Click *I’m Done* to exit the download window, and return to the console.

Next, move the downloaded `kubectl` file to `/usr/local/bin` (or any other directory in your PATH).

```sh
$ chmod +x kubectl
$ mv kubectl /usr/local/bin/kubectl
```

Make the downloaded `kubectl-config` file kubectl’s default by copying it to a `.kube` directory on your machine.

```sh
$ mkdir -p ~/.kube/ # create the directory
$ cp path/to/file/kubectl-config $HOME/.kube/config # rename the file and copy it into the directory
```

Once you've downloaded and copied the `kubectl-config` file to the `.kube` directory, you’re ready to start using `kubectl`. Protect `.kube/config` with appropriate file permissions as it contains cluster access credentials.

To invoke `kubectl` with a different `kubeconfig` than the default, name the desired `kubeconfig` in an environment variable or on the `kubectl` command line:
  * Set the environment variable `KUBECONFIG` to the location of the selected `kubectl-config` file.
    `$ export KUBECONFIG=/path/to/kubectl-config`
  * Or, use `kubectl` with the `–kubeconfig option`.
    `$ kubectl --kubeconfig=/path/to/kubectl-config get pods`

Once `kubectl` is properly configured, it can be used to explore Kubernetes entities:

```sh
$ kubectl get nodes
NAME         LABELS                              STATUS
10.0.0.197   kubernetes.io/hostname=10.0.0.197   Ready
10.0.0.198   kubernetes.io/hostname=10.0.0.198   Ready
10.0.0.199   kubernetes.io/hostname=10.0.0.199   Ready
```

Review our [kubectl documentation][kubedoc] for more setup help.

## Deploying a simple application

As an example application, we will deploy a simple, stateless website for a local bakery to the cluster: The Cookie Shop.

This example allows us to explore two useful Kubernetes concepts, `deployments` and `services`. Both are top-level Kubernetes objects, just like `nodes`.

[**Deployments:**][k8s-deployment] Run multiple copies of a container across multiple nodes

[**Services:**][k8s-service] Endpoint that load balances traffic to containers run by a deployment

Copy the following YAML into a file named `simple-deployment.yaml`.

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
```

The parameter `replicas: 3`, will create 3 running copies. `Image: quay.io/coreos/example-app:v1.0` defines the container image to run, hosted on [Quay.io][quay-repo].

Then, copy the following YAML into a file named `simple-service.yaml`.

```yaml
kind: Service
apiVersion: v1
metadata:
  name: simple-service
  namespace: default
spec:
  selector:
    k8s-app: simple
  ports:
  - protocol: TCP
    port: 80
  type: LoadBalancer
```

To connect the Service to the containers run by the Deployment, the Deployment `containerPort` and the Service `port` must match.

Instantiate the cluster objects specified in the `simple-deployment.yaml` and `simple-service.yaml` manifest by passing the filepaths to `kubectl create`. Check that they were created successfully by listing out the objects afterwards:

```sh
$ kubectl create -f simple-deployment.yaml
deployment "simple-deployment" created
$ kubectl get deployments
NAME                          DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/simple-deployment      3         3         3            3           7m
```

```sh
$ kubectl create -f simple-service.yaml
service "simple-service" created
$ kubectl get services -o wide
NAME                    CLUSTER-IP   EXTERNAL-IP                                                               PORT(S)        AGE       SELECTOR
svc/simple-service     10.3.0.204   a9b5de374e28611e6945f02c590b59c5-2010998492.us-west-2.elb.amazonaws.com   80:32567/TCP   7m        app=simple
```

The manifest specifies the pods required for a replicated Deployment, and connects them to a Kubernetes [external load balancer service][k8s-svc-lb]. On AWS, this service connects to an [Elastic Load Balancer (ELB)][aws-elb] through which it is exposed to the internet.

The `EXTERNAL-IP` column gives the DNS name of the externally routable port for the service.

Check your setup by navigating to the URL listed. (It may take a few minutes for AWS to set up the ELB, and for the URL to resolve.)

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12">
    <img src="img/cookie-v1.png">
    <div class="co-m-screenshot-caption">Our simple Cookie Shop application up and running on Tectonic</div>
  </div>
</div>

Use the Tectonic Console to monitor your app’s public IP, Service, Deployment, and related Pods.

Go to *Routing > Services* to monitor the site’s services.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/nginx-service.png">
    <div class="co-m-screenshot-caption">Viewing the Service in the Console</div>
  </div>
</div>

Go to *Workloads > Deployments* and click on the deployment name to monitor the deployment’s Pods.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="img/nginx-deployment.png">
    <div class="co-m-screenshot-caption">Viewing the Deployment in the Console</div>
  </div>
</div>

### Deploying an app with Tectonic Console

You can also deploy an application using the Tectonic Console by navigating to the Deployments page, creating a new deployment or service, and copying the above YAML contents.

This example will use the simple app deployed above. To create an identical app using the Tectonic Console, first delete the existing app using `kubectl delete`.

```sh
$ kubectl delete deploy/simple-deployment svc/simple-service
deployment "simple-deployment" deleted
service "simple-service" deleted
```

To deploy using the Tectonic console, use the content of the two YAML files already created.

First, deploy the sample app.
1. In the console, go to *Workloads > Deployments*, and click *Create Deployment*.
2. A pane will open, showing a default YAML deployment file.
4. Copy the contents of `simple-deployment.yaml`, and paste into the YAML pane, replacing its contents.
5. Click *Create*.

The Console will create your deployment, and display its Overview window.

Then, add the service.
1. Go to *Routing > Services*, and click *Create Service*.
2. Copy the contents of `simple-service.yaml` into the pane, replacing the default content.
4. Click *Create*.

The Console will create your service, and display its Overview window. Copy the *External Load Balancer* URL displayed into a browser to check your work. (It may take several minutes for AWS to update their ELB.)

## Using your own container images

The examples above used container images that have been shared publicly. To generate and host your own container images, we suggest using [Quay.io][quay-io] or [Quay Enterprise][QE]. The Quay container registry offers sophisticated access controls, easy automated builds, and automated security scanning, free for public projects.

Substitute your custom image and version (known as a "tag") in the Deployment above, by changing:

```yaml
      containers:
        - name: nginx
          image: quay.io/coreos/example-app:v1.0
```

[**NEXT:** Scaling an application with Tectonic][scale-app]


[QE]: https://coreos.com/quay-enterprise/
[assets-zip]: ../../admin/assets-zip.md
[aws-elb]: https://aws.amazon.com/elasticloadbalancing/
[edit-service]: img/edit-service.png
[k8s-deployment]: https://kubernetes.io/docs/user-guide/deployments/
[k8s-service]: https://kubernetes.io/docs/user-guide/services/
[k8s-svc-lb]: https://kubernetes.io/docs/user-guide/load-balancer/
[quay-io]: https://quay.io
[registry-auth]: https://coreos.com/os/docs/latest/registry-authentication.html
[quay-repo]: https://quay.io/repository/coreos/example-app
[kubedoc]: https://coreos.com/kubernetes/docs/latest/configure-kubectl.html
[installing]: installing-tectonic.md
[rolling-deployments]: rolling-deployments.md
[second-app]: second-app.md
[scale-app]: scale-app.md
[pull-secrets]: manage-pull-secrets.md
[namespaces]: manage-namespaces.md
