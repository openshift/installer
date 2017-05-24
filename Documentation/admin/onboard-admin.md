# Adding an administrator to a Tectonic cluster

An administrator can be assigned to a cluster or a specific namespace within a cluster. Use the *Role Binding* option in the Tectonic cluster to do so.

##  Prerequisites and guidelines

Before proceeding, ensure that the prerequisites given in the respective Identity Provider (IdP) section are met. Depending on the IdP used in the deployment, see one of the following:

* [Static user management][user-management]
* [LDAP user management][ldap-user-management]
* [SAML user management][saml-user-management]

## Granting access rights

Access rights are granted to an administrator user by using role binding.

### Setting up a Cluster administrator

Use the default role, *cluster-admin*, for cluster administration operation. Bind the role appropriately to the privileges.

1. Log in to the Tectonic UI.
2. Navigate to *Role Bindings* under *Administration*.
3. Click *Cluster-wide Role Binding*.
4. Select a desired namespace from the drop-down.
5. Select a Role Name.
   See [Default Roles in Tectonic][identity-management].
6. Select an appropriate subject from subject kind.
7. Specify a name to help identify the subject.
8. Click *Create Binding*.

### Setting up a Namespace administrator

1. Log in to the Tectonic UI
2. Navigate to *Role Bindings* under *Administration*.
3. Select a desired namespace from the drop-down.
   All the available Role Names based on Namespace Roles are populated in Role Name drop-down.
4. Select a Role Name from the drop-down.
4. Select the desired subject kind.
7. Specify a name to identify subject.
8. Click *Create Binding*.

[user-management]: user-management.md
[ldap-user-management]: ldap-user-management.md
[saml-user-management]: saml-user-management.md
[identity-management]: identity-management.md
