# Pull Secrets

Pull secrets are used to pull private container images from registries like Quay and Docker Hub.

Quay offers pre-populated Kubernetes pull secrets for each of your repositories. Check out the Quay docs for [more information][quay-registry-doc].

## Add a pull secret with kubectl

To add a Pull Secret via the command-line, write a Kubernetes `Secret` resource file. This includes a name, namespace, and base-64 encoded docker `config.json` file.

This can be created with Kubernetes by copying the output of the following command:

```
$ kubectl create secret docker-registry staging-secret \
    --docker-server=quay.io \
    --docker-username=giffeeLover93 \
    --docker-password='my secret passphrase' \
    --docker-email='giffeeLover93@example.com' \
    --dry-run -o yaml
```

Which outputs:

```yaml
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  name: staging-secret
data:
  .dockercfg: eyJxdWF5LmlvIjp7InVzZXJuYW1lIjoiZ2lmZmVlTG92ZXI5MyIsInBhc3N3b3JkIjoibXkgc2VjcmV0IHBhc3NwaHJhc2UiLCJlbWFpb CI6ImdpZmZlZUxvdmVyOTNAZXhhbXBsZS5jb20iLCJhdXRoIjoiWjJsbVptVmxURzkyWlhJNU16cHRlU0J6WldOeVpYUWdjR0Z6YzNCb2NtRnpaUT09In19
type: kubernetes.io/dockercfg
```

The Pull Secret is registered with Kubernetes by passing the path to the above YAML file to `kubectl create`:

```
$ kubectl create -f my-custom-staging-secret.yaml
secret "staging-secret" created
```

Default Pull Secrets for a given Namespace can also be added via the Tectonic interface.

![Namespace secret][namespace-secret]

Check out the [Kubernetes Pull Secret user guide][k8s-pull-secret-ug] for more information about Kubernetes Pull Secrets.

[namespace-secret]: ../img/walkthrough/namespace-secret.png
[k8s-pull-secret-ug]: https://kubernetes.io/docs/user-guide/images/#specifying-imagepullsecrets-on-a-pod
[quay-registry-doc]: https://coreos.com/os/docs/latest/registry-authentication.html
