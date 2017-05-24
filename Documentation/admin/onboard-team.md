# Adding a team to a Tectonic cluster

A group can be assigned to a cluster or a specific namespace within a cluster. Use the *Role Binding* option in the Tectonic cluster to do so.

##  Prerequisites and guidelines

Before proceeding, ensure that the prerequisites given in the respective Identity Provider (IdP) section are met. Depending on the IdP used in the deployment, see either of the following:

* [Static user management][user-management]
* [LDAP user management][ldap-user-management]
* [SAML user management][saml-user-management]

## Setting up a namespace administrator

Tectonic configures three default namespaces:  default, kube-system, and tectonic-system. The namespace administrator role will have full permission to the objects in a namespace. All Kubernetes clusters have two categories of users: service accounts managed by Kubernetes, and normal users. Service accounts are managed by the API server and can be created by using API calls. Normal user accounts are externally created and managed, such as by using an LDAP server or Google account. Kubernetes does not have corresponding objects representing normal user accounts.

## Granting access rights

Access rights are granted by using  roles. In order to grant access to a namespace:

### Cluster

1. Log in to the Tectonic UI.
2. Navigate to *Role Bindings* under *Administration*.
3. Click *Cluster-wide Role Binding*.
4. Select a desired namespace from the drop-down.
5. Select a Role Name.
   See [Default Roles in Tectonic][identity-management].
6. Select *Group* from subject kind.
7. Specify a a name to identify subject.
8. Click *Create Binding*.

### Namespace

1. Log in to the Tectonic UI.
2. Navigate to *Role Bindings* under *Administration*.
3. Click *Namespace Role Binding*.
4. Select a desired namespace from the drop-down.
5. Select a Role Name.
   See [Default Roles in Tectonic][identity-management].
6. Select *Group* from subject kind.
7. Specify a name to identify subject.
8. Click *Create Binding*.


## Users as part of multiple groups

Individual users can be part of multiple groups. The individual LDAP users or groups aren't viewable on the Tectonic console. However, the roles and role bindings attached to users and groups are displayed on the individual Roles page. Editing the YAML file associated with individual role is permitted to the role with necessary rights.  Creating a rule or role binding is allowed from the role detail page.

## Managing removed LDAP users and groups

When removed from LDAP, users and groups are cached.

[user-management]: user-management.md
[ldap-user-management]: ldap-user-management.md
[saml-user-management]: saml-user-management.md
[identity-management]: identity-management.md
