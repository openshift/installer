# User Management in Tectonic

User management in Tectonic is performed in two stages. The first stage involves authenticating a user and the second involves authorizing the user to perform a given set of tasks associated with a role. Authenticating users in Tectonic is managed by Tectonic Identity, whereas authorizing users is controlled by Role Based Access Control (RBAC). All Tectonic clusters enable RBAC which uses the user information produced by Tectonic Identity to enforce permissions. This document describes managing users and access control in Tectonic.

## Tectonic Identity

 Tectonic Identity is an authentication service for both Tectonic Console and `kubectl` and allows these components to talk to the API server on an end user's behalf. Users are either defined in the Tectonic `ConfigMap` or  integrated by using an external Identity Provider (IdP).  For more information see:

* [Static user management][user-management]
* [LDAP user management][ldap-user-management]
* [SAML user management][saml-user-management]

### Tectonic authentication through Dex

Tectonic Identity is built on top of [Dex][dex], an open-source OpenID Connect server. Dex binds with the backing LDAP server by using the end user's plain text password. Though some LDAP implementations allow passing hashed passwords, Dex doesn't support hashing and instead strongly recommends that all administrators use TLS. This is achieved by configuring port 636 instead of 389 in the Dex's `configMap`. Choosing port 389 can potentially cause leaking passwords.

In case of LDAP, Tectonic Console logs the user in to a LDAP server. LDAP through Dex issues a ID Token. Tectonic on behalf of the user call kubectl with the ID Token. Kubectl connects to the API Server and verify the validity of the ID Token and authenticity of the user. Once authorized API server returns a response to kubectl.

### Overview of Dex

Dex runs natively on top of Tectonic clusters by using [third-party resources][third-party], and drives API server authentication through the OpenID Connect plugin. Clients, such as the Tectonic Console and kubectl, act on behalf users who can log in to the Tectonic cluster through an identity provider, such as LDAP, that both Tectonic and Dex support.

Dex server issues short-lived, signed tokens on behalf of users. This token response, called ID Token, is a signed JSON web token. ID Token contains names, emails, unique identifiers, and a set of groups that can be used to identify a user. Dex publishes public keys, and Tectonic API server uses these to verify ID Tokens. The username and group information of a user is used in conjunction with RBAC to enforce authorization policy.

### Components of Tectonic Identity

The three major components of Tectonic Identity are the API server, Tectonic Console, and kubectl.

#### Tectonic API server

The Tectonic API Server is expected to enable it's OpenID Connect plugin, deferring to Dex for authentication. The API server is not a Dex client.

#### Kubectl

For Dex, kubectl is a public client. All kubectl instances share a `client_id` and `client_secret`, and the `client_secret` isn't considered private. kubectl communicates only with the API server.

#### Tectonic Console

Tectonic Console communicates with both Dex and the API server. Therefore, Tectonic Console is considered to be an admin client for Dex. However to be trusted by both Kubernetes and Dex, ID Tokens are issued to both Console and kubectl. When a user logs in to a Tectonic Console, Dex creates an ID Token for both Console and kubectl allowing Console to both authenticate with Kubernetes and the Dex APIs.

## RBAC in Tectonic

Authorization in  Tectonic is controlled by a set of permissions called Roles. Role Bindings grant the permissions associated with a Role to a subject. Subject is defined as a type of account used to access the Tectonic clusters.

There are three types of subjects in Tectonic:

* User
* Group
* Service Account

There are two types of roles in Tectonic:

* Roles: Scope of Roles is restricted to a namespace
* Cluster Roles: Cluster Roles are restricted to a physical cluster

Access is granted to a role based on how it's bound. For the same purpose, Tectonic has two types of role binding:

 * Namespace Role Binding: Defines permission at namespace level for users or a group of users
 * Cluster-wide Role Binding: Defines permission at cluster level for users or a group of users

### Default Roles in Tectonic

Tectonic inherits most of the roles from Kubernetes upstream.  There are Cluster-wide Roles, Namespace Roles, and System Roles. ingress-controller is a Cluster-wide role to handle Tectonic ingress traffic.

The default Cluster-wide roles in tectonic are:

| Cluster Roles | Permissions   |
| ------------- |:-------------|
| cluster-admin | Full control over all the objects in a cluster.|
| admin         | Full control over all objects in a namespace. Bind this role into a namespace to give administrative control to a user or group.|
| user          | Access to all common objects, either within a namespace or cluster-wide, but is prevented from changing the RBAC policies. |
| readonly      | Read only view for all objects. Can be used cluster-wide, or just within a specific namespace.|
| ingress-controller |Designed to be used by an ingress controller to configure routing.|


[user-management]: user-management.md
[ldap-user-management]: ldap-user-management.md
[saml-user-management]: saml-user-management.md
[dex]: https://github.com/coreos/dex
[third-party]: https://github.com/coreos/dex/blob/master/Documentation/storage.md#Kubernetes-third-party-resources
