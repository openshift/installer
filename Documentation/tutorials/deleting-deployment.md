# Deleting an application on Tectonic

In this tutorial you will learn how to delete a deployment using `kubectl` and the Tectonic Console.

## Deleting an application with kubectl

First, use `kubectl get all` to list all applications and services on your cluster.

```sh
$ kubectl get all
NAME                                    READY     STATUS    RESTARTS   AGE
po/nginx-deployment-1447934386-45snm    1/1       Running   1          3d
po/nginx-deployment-1447934386-4q0ms    1/1       Running   1          3d
po/nginx-deployment-1447934386-mjdd5    1/1       Running   1          3d
po/simple-deployment-305980935-1knc5    1/1       Running   0          19h
po/simple-deployment-305980935-d37lv    1/1       Running   0          19h
po/simple-deployment-305980935-lr326    1/1       Running   0          19h

NAME                     CLUSTER-IP   EXTERNAL-IP        PORT(S)        AGE
svc/kubernetes           10.3.0.1     <none>             443/TCP        4d
svc/simple-service       10.3.0.178   a1c717248163e...   80:30984/TCP   3d
svc/simple-service-two   10.3.0.56    a16b97a7b1644...   80:32185/TCP   3d

NAME                           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/cookies                 5         5         5            5           3d
deploy/nginx-deployment        3         3         3            3           3d
deploy/simple-deployment       3         3         3            3           3d
deploy/simple-deployment-two   3         0         0            0           3d

NAME                              DESIRED   CURRENT   READY     AGE
rs/nginx-deployment-1447934386    3         3         3         3d
rs/simple-deployment-169731590    0         0         0         19h
rs/simple-deployment-256231244    0         0         0         3d
rs/simple-deployment-305980935    3         3         3         19h
```

Next, delete the `simple` app's Deployment and Service using `kubectl delete`:

```sh
$ kubectl delete deployment/simple-deployment svc/simple-service
deployment "simple-deployment" deleted
service "simple-service" deleted
```

This will delete the Deployment (including Pods and Replica Sets) and the Service associated with the `simple` app created earlier.

```sh
$ kubectl get all
NAME                                       READY     STATUS    RESTARTS   AGE
po/nginx-deployment-1447934386-45snm       1/1       Running   1          3d
po/nginx-deployment-1447934386-4q0ms       1/1       Running   1          3d
po/nginx-deployment-1447934386-mjdd5       1/1       Running   1          3d
po/simple-deployment-two-256231244-1j2nv   1/1       Running   0          1m
po/simple-deployment-two-256231244-mvcl9   1/1       Running   0          1m
po/simple-deployment-two-256231244-z142j   1/1       Running   0          1m

NAME                     CLUSTER-IP   EXTERNAL-IP        PORT(S)        AGE
svc/kubernetes           10.3.0.1     <none>             443/TCP        4d
svc/simple-service-two   10.3.0.56    a16b97a7b1644...   80:32185/TCP   3d

NAME                           DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/cookies                 5         5         5            5           3d
deploy/nginx-deployment        3         3         3            3           3d
deploy/simple-deployment-two   3         3         3            3           3d

NAME                                 DESIRED   CURRENT   READY     AGE
rs/nginx-deployment-1447934386       3         3         3         3d
rs/simple-deployment-two-256231244   3         3         3         1m
```

## Deleting an application with the Tectonic Console

To delete an application using the Tectonic Console:

1. Go to *Workloads > Deployments*.
2. Click the gear icon next to the deployment, select *Delete Deployment...* and confirm.
3. Go to *Routing > Services*.
4. Click the gear icon next to the Service, select *Delete Service...* and confirm.
