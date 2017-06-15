# Failure domains of Tectonic Identity

## Architectural Overview:
![][dex-architecture]

Each of the components on the left, depend on Tectonic Identity to fulfil its functionality and the components on the right provide services that Tectonic Identity needs. The Tectonic Console, kubectl, and API server all talk to one or more Tectonic Identity instances through a load balancer. This load balancer use a stable DNS name and must be available to users of kubectl to refresh their identity tokens.

Below is a brief summary of how each module interacts with Tectonic Identity:

* **Tectonic Console**: It depends on Tectonic Identity to login users and create a browser session cookie via the graphical login interface.
* **Kubectl**: Uses an identity token to authenticate to the API server and every few hours uses a refresh token to renew it.
               Known issue: kubectl makes a roundtrip on every request, this will be fixed in kubectl v1.7.0+ and potentially a patch release to v1.6. Issue documented [here][kubernetes-issue]
* **API server**: It uses the `/.well-known` and `/keys` endpoints of Tectonic Identity. It does periodic refreshes from the `/keys` Tectonic Identity endpoint, to authenticate users making requests to the API server. The API server does not require access to Tectonic Identity on every request and caches these public keys allowing the API to be resilient to temporary outages of Tectonic Identity, the load balancer, or the network.
* **NGINX Ingress Controller**: It allows inbound connections to reach Tectonic Identity. It performs health checks at the `/identity/healthz` endpoint to ensure Tectonic Identity is up and running. It can also load balance across several instances of Tectonic Identity for high-availability of the service to protect against single process failures and machine failures.
* **Authentication providers**: Tectonic Identity allows federating user login to a remote identity service e.g. Github, Google, LDAP. To authenticate a user via the login interface.
* **Backend storage**: Tectonic Identity requires persisting state to perform various tasks such as track refresh tokens, preventing replays, periodically rotating keys, etc. Otherwise, the service is primarily stateless. In Tectonic this persistent state is provided by Kubernetes third party resources which serves as the backend storage.

## Failures

This section highlights the scenarios in which a failure can arise due to one of the architectural components malfunctioning or because different components cannot interact.

### Tectonic Identity Instance(s) Fails or Becomes Unreachable

This scenario includes failures of the load balancer fronting Tectonic Identity or the backend instances of Tectonic Identity behind that load balancer. Generally, Tectonic Identity is run as a deployment on top of Kubernetes; meaning if a Tectonic Identity instance stops responding to the /identity/healthz endpoint it will be rescheduled. Also note that it is possible to run multiple instances of Tectonic Identity to avoid a single point of failure.

**Mean time to recovery (MTTR)** for a Tectonic Identity instance that has failed will be 5s-60s to launch a new container or download the Tectonic Identity container and launch.

**Downtime** of all Tectonic Identity API endpoints will incur if ALL instances of Tectonic Identity are down OR the load balancer is down/unreachable by clients.
   * Low risk to clients that cache keys, like the Kubernetes API server, from seconds to minutes of downtime.
   * Higher Risk to clients that require access to Tectonic Identity APIs to work. e.g. Starting a Kubernetes API server.
   * Higher Risk if Tectonic Identity doesn't recover that keys won't be rotated before expiration. This would require downtime of hours.

### Authentication provider

If a user has chosen an external authentication provider, say LDAP for example, to be the login mechanism for Tectonic and the LDAP server is down, users will be unable to login to cluster. If the LDAP server goes down after the user has logged in, Tectonic Identity will not able to request a new access token when the current access token expires. The token expiry time is configurable.

**Mean time to recovery (MTTR)** for Tectonic Identity from an authentication provider being down is 0 seconds as long as the user has logged in. Authentication providers are only required when a user attempts to login.

**Downtime** of all Tectonic Identity API endpoints will incur if a new user has to be authenticated or an existing user's credentials need to be refreshed.
   * Risk to users relying on the authentication provider.

### Backend storage malfunctions

Tectonic Identity relies on the Kubernetes third party resources for any sort of persistent state information. Any task that depends on reading or writing to the storage will throw a “Database Error”. This affects automatic access token refreshing, signing key rotation, etc.

**Mean time to recovery (MTTR)** for Tectonic Identity from a backend malfunction is 0 seconds as Tectonic Identity APIs, which require storage, will begin returning 200 erros.

**Downtime** of all Tectonic Identity many API endpoints will incur IF the backend storage malfunctions.
   * Low risk to clients that cache keys, like the Kubernetes API server, from seconds to minutes of downtime. 
   * Risk to clients that require access to Tectonic Identity APIs to work. e.g. Starting a Kubernetes API server


## Simulating Downtime

### Tectonic Identity Instance(s) Fails or Becomes Unreachable

The easiest ways for testing outage is to scale the Tectonic Identity deployment on Kubernetes down to 0. For example for clusters installed using the the Tectonic Installer simply run:
```
$ kubectl scale deployment tectonic-identity -n tectonic-system --replicas=0
```

### Backend storage malfunctions

The easiest way to test backend failure is to revoke the service account roles to Tectonic Identity.

## FAQ
Q: How does a user recover or reconfigure a cluster in case Tectonic Identity fails?

A: If Tectonic Identity is down and the user is unable to login via the Tectonic console they can make use of the kubeconfig generated by the installer which can be found in the [assets folder](assets-zip.html). This will allow the user to access the kubernetes API directly.

Q: Is there a circular dependency between the API server and Tectonic Identity? For example Tectonic Identity is used to authenticate requests and Tectonic Identity also uses the Kubernetes API server to persist data.

A: The Kubernetes API server can authenticate requests in three ways: using OIDC via Tectonic Identity for users, using TLS client certificate for disaster recovery by users, and [service accounts](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/) for on-cluster services. The Tectonic Identity service uses a service account to talk to the API server which resolves the dependency.


[dex-architecture]: ../img/failure-domains-identity.png
[kubernetes-issue]: https://github.com/kubernetes/kubernetes/issues/42654#issuecomment-297832539

