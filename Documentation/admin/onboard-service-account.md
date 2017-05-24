# Adding a service account to Tectonic cluster

## About service accounts

Service accounts are API credentials stored in the Kubernetes API and mounted into pods at well known paths, giving the pod an identity which can be access controlled. Pods use service accounts to authenticate against Kubernetes API from within the cluster. If an app uses kubectl or the official Kubernetes Go client within a pod to talk to the API, these credentials are loaded automatically.

Since RBAC denies all requests unless explicitly allowed, service accounts, and the pods that use them, must be granted access through RBAC rules.
Kubernetes automatically creates a "default" service account in every namespace. If pods don't explicitly request a service account, they're assigned to this "default" one.

A service account, `default`, by default is created by Tectonic. However, creating an additional service account is permitted.

## Creating an additional service accounts

To create an additional service account, for example, an ingress role, either create a yaml file as follows or use the Tectonic console to create one. Given is an example service account for ingress.

### Using an YAML file

1. Define the role, `ingress.yaml`, which gives administrative privileges to the service account within the default namespace:


      apiVersion: rbac.authorization.k8s.io/v1alpha1
      kind: ClusterRoleBinding
      metadata:
        name: public-ingress
      roleRef:
        apiGroup: rbac.authorization.k8s.io
        kind: ClusterRole
        name: admin
      subjects:
      - kind: ServiceAccount
        name: default
        namespace: public

2. Run the following:

  `kubectl create serviceaccount `ingress.yaml`
   serviceaccount "ingress" created

If multiple pods running in the same namespace require different levels of access, create a unique service account for each. The newly created service account can be mounted onto the pod by specifying the service account name in the pod spec.

    apiVersion: extensions/v1beta1
    kind: Deployment
    metadata:
      name: nginx-deployment
    spec:
      replicas: 3
      template:
        metadata:
          labels:
            k8s-app: nginx
        spec:
          containers:
          - name: nginx
            image: nginx:1.7.9
          serviceAccountName: ingress # Specify the custom service account

### Using the Tectonic console

## Granting access rights to service accounts

Use the Tectonic console to provide access rights to a service account.

### Namespace

1. Log in to the Tectonic UI
2. Navigate to *Role Bindings* under *Administration*.
3. Click *Namespace Role Binding*.
4. Select a desired namespace from the drop-down.
5. Select a Role Name.
   See [Default Roles in Tectonic][identity-management].
6. Select *Service Account* from subject kind.
7. Specify a a name to identify subject.
8. Click *Create Binding*.


### Cluster-wide

1. Log in to the Tectonic UI
2. Navigate to *Role Bindings* under *Administration*.
3. Click *Cluster-wide Role Binding*.
4. Select a desired namespace from the drop-down.
5. Select a Role Name.
   See [Default Roles in Tectonic][identity-management].
6. Select *Service Acccount* from subject kind.
7. Specify a a name to identify subject.
8. Click *Create Binding*.

[identity-management]: identity-management.md
