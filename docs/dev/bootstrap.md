# Bootstrap Node

## Static-Pod Logging
Static pods running on the bootstrap node should write logs to one of the following locations:
* Control Plane Pods - `/var/log/boostrap-control-plane`
* All Others - `/var/log/bootstrap-pods`

The appropriate directory can be mounted as a `hostPath` as seen with the [API Server][api-server].


[api-server]: https://github.com/openshift/cluster-kube-apiserver-operator/blob/master/bindata/bootkube/bootstrap-manifests/kube-apiserver-pod.yaml