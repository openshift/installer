# Azure Authentication using Azure CLI (az login)

The installer can use an Azure CLI session (`az login`) for its own Azure API calls, in addition to
client secret, client certificate, and managed identity.

### Pitfalls
Although the installer can now use `az login` to authenticate, CCO does not support this and hence the installer
should create the cluster in manual credentials mode only.

`openshift-install create cluster` still needs a service principal or managed identity.
Use `az login` for installer API calls (validation, metadata, `create install-config`).

### Prerequisites
- Azure CLI is installed. More information is in [1].
- An active session via `az login`.
- `credentialsMode: Manual` in the install-config when using this path.

## Steps
After [2], the installer can authenticate without a service principal file when Azure CLI is logged in.

1. Log in to Azure

`az login`

2. (Optional) Set the subscription to use

`az account set --subscription <subscription-id>`

3. Set `credentialsMode: Manual` in `install-config.yaml` and run installer commands that only need installer API access.

The installer reads the default subscription and tenant from `~/.azure/azureProfile.json` and uses `azidentity.NewAzureCLICredential` [3].

No `osServicePrincipal.json` file or `AZURE_AUTH_LOCATION` environment variable is needed for those installer API calls.

If `~/.azure/osServicePrincipal.json` (or the path in `AZURE_AUTH_LOCATION`) is present, the installer uses that file instead of Azure CLI.

Certificate, client secret, and managed identity continue to work as before. For certificates, see [Azure Authentication using Client certificates](azure_client_certs_auth.md).

To create a cluster, provide a service principal (`~/.azure/osServicePrincipal.json`) or managed identity.

### References
[1] - Install Azure CLI : https://learn.microsoft.com/en-us/cli/azure/install-azure-cli
[2] - Azure CLI authentication : https://learn.microsoft.com/en-us/cli/azure/authenticate-azure-cli
[3] - Azure CLI credential : https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity#AzureCLICredential
