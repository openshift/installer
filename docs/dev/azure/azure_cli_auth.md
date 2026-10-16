# Azure Authentication using Azure CLI (az login)

The installer can use an Azure CLI session (`az login`) for its own Azure API calls, in addition to
client secret, client certificate, and managed identity.

> **Scope:** This procedure enables Azure CLI authentication for installer API operations such as `create install-config`, validation, and metadata generation. It does not, by itself, provide all credentials needed to complete `openshift-install create cluster`.
>
> For a full cluster installation, use this procedure only for the installer API operations described above. You must separately prepare the CCO `CredentialsRequest` Secrets required by `credentialsMode: Manual` and configure CAPZ with a supported Azure identity: a service principal with a client secret, a client certificate, or a managed identity.
>
> For this procedure, `credentialsMode: Manual` controls how CCO satisfies `CredentialsRequest` objects; it does not select the installer authentication method or provide the local `az login` identity to CAPZ.

## Pitfalls
When using Azure CLI authentication, set `credentialsMode: Manual` (step 5) because CCO cannot use the local Azure CLI session.

The Azure CLI cloud must match install-config `cloudName`.

## Prerequisites
- Azure CLI is installed. More information is in [1].
- An Azure account with access to the target tenant and subscription.

## Steps
1. (Optional) If you are using Azure Government, China, or another non-public cloud, set the CLI cloud so it matches install-config `cloudName`

`az cloud set --name <matching-cloud>`

2. Log in to Azure ([2])

`az login`

3. (Optional) Set the subscription to use

`az account set --subscription <subscription-id>`

4. If you do not already have an `install-config.yaml`, run `openshift-install create install-config`.
5. Edit `install-config.yaml` and set `credentialsMode: Manual` before running any subsequent installer commands that use this authentication path.
6. Run the installer commands that only need installer API access.

When no service principal file is present, the installer can authenticate with the Azure CLI session after `az login`.
The installer reads the default subscription and tenant from `~/.azure/azureProfile.json` and uses `azidentity.NewAzureCLICredential` [3].

No `osServicePrincipal.json` file or `AZURE_AUTH_LOCATION` environment variable is needed for those installer API calls.

If `~/.azure/osServicePrincipal.json` (or the path in `AZURE_AUTH_LOCATION`) is present, the installer uses that file instead of Azure CLI.

Certificate, client secret, and managed identity continue to work as before. For certificates, see [Azure Authentication using Client certificates](azure_client_certs_auth.md).

## References
[1] - Install Azure CLI : https://learn.microsoft.com/en-us/cli/azure/install-azure-cli
[2] - Azure CLI authentication : https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli
[3] - Azure CLI credential : https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#AzureCLICredential
