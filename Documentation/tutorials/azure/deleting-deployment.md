# Deleting an application on Tectonic

This tutorial demonstrates how to delete an application, and return a Tectonic cluster to its original state using `kubectl` and Tectonic Console.

## Deleting an application with Tectonic Console

To delete an application using the Tectonic Console:

1. Go to *Workloads > Deployments*.
2. Click the gear icon next to the deployment, select *Delete Deployment...* and confirm.
3. Go to *Routing > Services*.
4. Click the gear icon next to the Service, select *Delete Service...* and confirm.

Go to the URL created while [deploying an app][deployed-app] in the first tutorial to confirm that the app has been deleted.

## Deleting an application with kubectl

First, use `kubectl get all` to list all applications and services on the cluster.

```sh
$ kubectl get all
NAME                                    READY     STATUS    RESTARTS   AGE
po/simple-deployment-3220091887-070cw   1/1       Running   0          7m
po/simple-deployment-3220091887-h293v   1/1       Running   0          7m
po/simple-deployment-3220091887-s4szh   1/1       Running   0          7m

NAME                 CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
svc/kubernetes       10.3.0.1      <none>        443/TCP        1h
svc/simple-service   10.3.217.22   <none>        80:32739/TCP   37m

NAME                       DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/simple-deployment   3         3         3            3           37m

NAME                              DESIRED   CURRENT   READY     AGE
rs/simple-deployment-3126793206   0         0         0         37m
rs/simple-deployment-3220091887   3         3         3         7m
```

Next, delete the `simple` app's Deployment and Service using `kubectl delete`:

```sh
$ kubectl delete deploy/simple-deployment svc/simple-service
deployment "simple-deployment" deleted
service "simple-service" deleted
```

This will delete the Deployment (including Pods and Replica Sets) and the Service associated with the `simple` app created in the [Deploying an application on Tectonic][first-app] tutorial.

As of version 1.7, Ingress resources no longer appear in the `all` group. Use the same sequence of `kubectl get ingress` and `kubectl delete` to delete the Ingress resource.

```sh
$ kubectl get ingress
NAME             HOSTS                         ADDRESS        PORTS     AGE
simple-ingress   my-cluster.example.com   172.17.4.201   80        44m
$ kubectl delete ing/simple-ingress
```


[first-app]: first-app.md
[deployed-app]: first-app.md#deploying-an-app-with-tectonic-console
