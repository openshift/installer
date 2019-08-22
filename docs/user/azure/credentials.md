# Service Principal

Before proceeding with the OpenShift install, you should create a service principal with administrative rights for your subscription following the steps
outlined here:

[Azure: Creating an Service Principal][sp-create]

## Step 1: Create a Service Principal

You can create a Service Principal using the Azure [portal][sp-create-portal] or the Azure [cli][sp-create-cli]

## Step 2: Request permissions for the Service Principal from Tenant Administrator

In order to properly mint credentials for components in the cluster, your service principal needs to request for the following Application [permissions][ad-permissions] before you can deploy OpenShift on Azure: `Azure Active Directory Graph -> Application.ReadWrite.OwnedBy`

You can request permissions using the Azure portal or the Azure cli.

### Requesting permissions using the Azure cli

Find the AppId for your service principal by using,

```console
$ az ad sp list --show-mine -otable
AccountEnabled    AppDisplayName     AppId                                 AppOwnerTenantId                      AppRoleAssignmentRequired    DisplayName        Homepage                   ObjectId                              ObjectType        Odata.type                                    PublisherName    ServicePrincipalType    SignInAudience
----------------  -----------------  ------------------------------------  ------------------------------------  ---------------------------  -----------------  -------------------------  ------------------------------------  ----------------  --------------------------------------------  ---------------  ----------------------  ----------------
...
```

Use can request `Application.ReadWrite.OwnedBy` permission by using,

```sh
az ad app permission add --id <AppId> --api 00000002-0000-0000-c000-000000000000 --api-permissions 824c81eb-e3f8-4ee6-8f6d-de7f50d565b7=Role
```

NOTE: `Application.ReadWrite.OwnedBy` permission is granted to the the application only after it is provided an [`Admin Consent`][ad-admin-consent] by the Tenant Administrator.

## Step 3: Attach Administrative Role

Azure installer creates new identities for the cluster and therefore requires access to create new roles, and role assignments. Therefore, you will require the service principal to have at least `Contributor` and `User Access Administrator` [roles][built-in-roles] assigned in your subscription.

You can create role assignments for your service principal using the Azure [portal][sp-assign-portal] or the Azure [cli][sp-assign-cli]

## Step 4: Acquire Client Secret

You need to save the client secret values to configure your local machine to run the installer. This step is your opportunity to collect those values, and additional credentials can be added to the service principal in the Azure portal if you didn't capture them.

You can get client secret for your service principal using the Azure [portal][sp-creds-portal] or the Azure [cli][sp-creds-cli]

The default location of the service principal file is in ${HOME}/.azure/osServicePrincipal.json. An alternative location can be specified by setting the AZURE_AUTH_LOCATION environmental variable. You can generate a service principal file by running the following commands from a shell:

````console
$ az ad sp create-for-rbac --role Owner --name team-installer | jq --arg sub_id "$(az account show | jq -r '.id')" '{subscriptionId:$sub_id,clientId:.appId, clientSecret:.password,tenantId:.tenant}' > ~/.azure/osServicePrincipal.json
````

Once a credentials file has been generated, and the proper permissions have been set for your account, you can [install an OpenShift cluster](install.md).

[ad-admin-consent]: https://docs.microsoft.com/en-us/azure/active-directory/develop/v1-permissions-and-consent#types-of-consent
[ad-permissions]: https://docs.microsoft.com/en-us/azure/active-directory/develop/v1-permissions-and-consent
[sp-create]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals
[sp-create-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#create-service-principal-for-azure-ad
[sp-create-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#create-a-service-principal
[built-in-roles]: https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
[sp-assign-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#assign-the-service-principal-to-a-role
[sp-assign-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#manage-service-principal-roles
[sp-creds-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#get-credentials
[sp-creds-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#reset-credentials
