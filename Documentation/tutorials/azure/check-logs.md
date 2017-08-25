# Checking application logs

Tectonic provides access to application logs both in Console and with `kubectl`.

This guide continues from the previous tutorial, where we [installed a simple application][first-app].

## Use Tectonic Console to check logs

Use Tectonic Console to view application logs.

1. Go to `https://my-cluster.example.com` to open the Console and log in.
2. Click *Workloads > Pods* from the left hand navigation bar.
3. In the upper right corner of the page, enter `simple`.
4. Click the *Pod Name* for any one of the three pods shown.
5. Click the Pod's *Logs* tab.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12">
    <img src="img/viewing-logs-simple-deployment.png">
    <div class="co-m-screenshot-caption">Logs for our simple Cookie Shop app</div>
  </div>
</div>

Console will display available information for the pod, including CPU, memory, and logging stats. Click *Logs* to confirm the output is the same as that shown using `kubectl logs`.

## Use kubectl to check logs

Use `kubectl get pods` to view the cluster's pods.

```sh
$ kubectl get pods
NAME                                 READY     STATUS    RESTARTS   AGE
simple-deployment-4098151155-494nl   1/1       Running   0          1m
simple-deployment-4098151155-n8bqr   1/1       Running   0          1m
simple-deployment-4098151155-p680w   1/1       Running   0          1m
```

Copy one of the pod's names, and append it to the `kubectl logs -f ` command.

```sh
kubectl logs -f simple-deployment-4098151155-n8bqr
10.2.1.1 - - [15/Aug/2017:21:30:32 +0000] "GET / HTTP/1.1" 200 576 "-" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.90 Safari/537.36" "127.0.0.1"
```

Open the app in your browser by visiting [my-cluster.example.com/simple-deployment][visit-app]. There is a one in three chance that any visit to the page will hit any one of the three available pods. Refresh the page several times to increase the likelihood that the selected pod will produce log messages.

[**NEXT:** Scaling an application with Tectonic][scale-app]


[first-app]: first-app.md
[scale-app]: scale-app.md
[visit-app]: https://my-cluster.example.com/simple-deployment
