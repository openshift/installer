# Service Principal

Before proceeding with the OpenShift install, you should create a service principal with administrative rights for your subscription following the steps
outlined here:

[Azure: Creating an Service Principal][sp-create]

## Step 1: Create a Service Principal

You can create a Service Principal using the Azure [portal][sp-create-portal] or the Azure [cli][sp-create-cli]

## Step 2: Attach Administrative Role

Azure installer creates new identities for the cluster and therefore requires access to create new roles, and role assignments. Therefore, you will require the service principal to have at least `Contributor` and `User Access Administrator` [roles][built-in-roles] assigned in your subscription.

You can create role assignments for your service principal using the Azure [portal][sp-assign-portal] or the Azure [cli][sp-assign-cli]

## Step 3: Acquire Client Secret

You need to save the client secret values to configure your local machine to run the installer. This step is your opportunity to collect those values, and additional credentials can be added to the service principal in the Azure portal if you didn't capture them.

You can get client secret for your service principal using the Azure [portal][sp-creds-portal] or the Azure [cli][sp-creds-cli]

[sp-create]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals
[sp-create-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#create-service-principal-for-azure-ad
[sp-create-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#create-a-service-principal
[built-in-roles]: https://docs.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
[sp-assign-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#assign-the-service-principal-to-a-role
[sp-assign-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#manage-service-principal-roles
[sp-creds-portal]: https://docs.microsoft.com/en-us/azure-stack/user/azure-stack-create-service-principals#get-credentials
[sp-creds-cli]: https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest#reset-credentials