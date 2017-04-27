# User Management through Tectonic Identity

## Overview

Tectonic Identity is an authentication service for both Tectonic Console and `kubectl` and allows these components to talk to the API server on an end user's behalf. All Tectonic clusters also enable Role Based Access Control (RBAC) which uses the user information produced by Identity to enforce permissions.

This document describes managing users and access control in Tectonic.

## Identity Configuration

Tectonic Identity pulls all its configuration options from a config file stored in a `ConfigMap`, which admins can view and edit using `kubectl`. As a precaution, it's recommended to use the administrative kubeconfig in their [downloaded `assets.zip`][assets-zip] when editing Identity's config in case of misconfiguration.

First, backup the existing config using `kubectl`:

```
kubectl get configmaps tectonic-identity --namespace=tectonic-system -o yaml > identity-config.yaml.bak
```

Edit the current `ConfigMap` with the desired changes:

```
kubectl edit configmaps tectonic-identity --namespace=tectonic-system
```

Trigger a rolling update using `kubectl`. Identity's deployment is intended to be resilient against invalid config files, but admins should verify the new state and restore the `ConfigMap` backup if Identity enters a crash loop. The following command will cause an update:

```
kubectl patch deployment tectonic-identity \
    --patch "{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"`date +'%s'`\"}}}}}" \
    --namespace tectonic-system
```

The update's success can then be inspecting by watching the pods in the `tectonic-system` namespace.

```
kubectl get pods --namespace=tectonic-system
```

### Add static user

Static users are those defined directly in the Identity `ConfigMap`. Static users are intended to be used for initial setup, and potentially for troubleshooting and recovery. A static user acts as a stand-in, authenticating users without a connection to a backend Identity provider. To add a new static user, update the tectonic-identity `ConfigMap` with a new `staticPasswords` entry.

```yaml
    staticPasswords:
    # All the following fields are required.
    - email: "test1@example.com"
      # Bcrypt hash for string "password"
      hash: "$2a$10$2b2cU8CPhOTaGrs1HRQuAueS7JTT5ZHsHSzYiFPm1leZck7Mc8T4W"
      # username to display. NOT used during login.
      username: "test1"
      # Randomly generated using uuidgen.
      userID: "1d55c7c4-a76d-4d74-a257-31170f2c4845"
```

A Bcrypt encoded hash of the user's password can be generated using the using the [coreos/bcrypt-tool](https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0).

To ensure the static user has been added successfully try and log in with the new user from the Tectonic console.

### Change Password for Static User

To change the password of an exisiting user, generate a bcrypt hash for the new password (using [coreos/bcrypt-tool](https://github.com/coreos/bcrypt-tool/releases/tag/v1.0.0)) and plug in this value into the tectonic-identity `ConfigMap` for the selected user.

```yaml
    staticPasswords:
    # Existing user entry.
    - email: "test1@example.com"
      # Newly generated Bcrypt hash
      hash: "$2a$10$TcWtvcw0N8.xK8nKdBw80uzYij6cJwuQhtfYnEf/hEC9bRTzlWdIq"
      username: "test1"
      userID: "1d55c7c4-a76d-4d74-a257-31170f2c4845"
```

After the config changes are applied, the user can log in to the console using the new password.

### Add ClusterRoleBindings with Role Based Access Control (RBAC)

`ClusterRoles` grant access to types of objects in any namespace in the cluster. Tectonic comes preloaded with three `ClusterRoles`:

1. user
2. readonly
3. admin

`ClusterRoles` are applied to a `User`, `Group` or `ServiceAccount` via a `ClusterRoleBinding`. A `ClusterRoleBinding` can be used to grant permissions to users in all namespaces across the entire cluster, where as a `RoleBinding` is used to grant namespace specific permissions. The following `ClusterRoleBinding` resource definition adds an exisiting user to the admin role.

```yaml
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1alpha1
metadata:
  name: admin-test
# subjects holds references to the objects the role applies to.
subjects:
  # May be "User", "Group" or "ServiceAccount".
  - kind: User
    # Preexisting user's email
    name: test1@example.com
# roleRef contains information about the role being used.
# It can only reference a ClusterRole in the global namespace.
roleRef:
  kind: ClusterRole
  # name of an existing ClusterRole, either "readonly", "user", "admin",
  # or a custom defined role.
  name: admin
  apiGroup: rbac.authorization.k8s.io
```

The above `ClusterRoleBinding` RBAC resource definition can be applied through `kubectl`.

```
kubectl create -f admin-test.yaml
```

The new `ClusterRoleBinding` can viewed on the Tectonic Console under the Administration tab.

The `ClusterRoleBinding` can be deleted to revoke users' permissions.

```
kubectl delete -f admin-test.yaml
```

For additional details see the [Kubernetes RBAC documentation][k8s-rbac].

### Managing in-cluster API access

Pods use service accounts to authenticate against Kubernetes API from within the cluster. Service accounts are API credentials stored in the Kubernetes API and mounted into pods at well known paths, giving the pod an identity which can be access controlled. If an app uses `kubectl` or the official Kubernetes Go client within a pod to talk to the API, these credentials are loaded automatically.

Since RBAC denys all requests unless explicitly allowed, service accounts, and the pods that use them, must be granted access through RBAC rules.

Kubernetes automatically creates a "default" service account in every namespace. If pods don't explicitly request a service account, they're assigned to this "default" one.

```
$ kubectl get serviceaccounts
NAME               SECRETS   AGE
default            1         1h
$ kubectl create deployment nginx --image=nginx
deployment "nginx" created
$ kubectl get pods
NAME                     READY     STATUS    RESTARTS   AGE
nginx-3121059884-x7btf   1/1       Running   0          20s
```

If we inspect the `spec` of the pod, we'll see that the pod of the deployment has been assigned the "default" service account:

```
$ kubectl get pod nginx-3121059884-x7btf -o yaml
# ...
spec:
  containers:
  - image: nginx
    imagePullPolicy: Always
    name: nginx
    volumeMounts:
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: default-token-twmyd
      readOnly: true
  serviceAccountName: default
  volumes:
  - name: default-token-twmyd
    secret:
      defaultMode: 420
      secretName: default-token-twmyd
# ...
```

To allow the pod to talk to the API server, create a `Role` for the account, then use a RoleBinding to grant the service account the role's powers. For example, if the pod needs to be able to read `ingress` resources:

```yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1alpha1
metadata:
  name: default-service-account
  namespace: default
rules:
  - apiGroups: ["extensions"]
    resources: ["ingress"]
    verbs: ["get", "watch", "list"]
    nonResourceURLs: []
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1alpha1
metadata:
  name: default-service-account
  namespace: default
subjects:
  # The subject is the target service account
  - kind: ServiceAccount
    name: default
    namespace: default
roleRef:
  # The roleRef specifies the role to give to the
  # service account.
  kind: Role
  namespace: default
  name: default-service-account # Tectonic also provides "readonly", "user", and "admin" cluster roles.
  apiGroup: rbac.authorization.k8s.io
```

If multiple pods running in the same namespace require different levels of access, create a unique service account for each.

```
$ kubectl create serviceaccount my-robot-account
serviceaccount "my-robot-account" created
```

The newly created service account can be mounted into the pod by specifying the service account's name in the pod spec.

```yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
      serviceAccountName: my-robot-account # Specify the custom service account
```

The `RoleBinding` would then reference the custom service account name instead of "default".

Note that because service account credentials are stored in secrets, any clients with the ability to read secrets can extract the bearer token and act on behalf of that service account. It's recommended to be cautious when giving service accounts powers or clients the ability to read secrets.

[assets-zip]: assets-zip.md
[k8s-rbac]: http://kubernetes.io/docs/admin/authorization/#rbac-mode
