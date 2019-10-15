# GCP Service APIs
To install OpenShift to your GCP project, the installer requires following [Service APIs][service-apis-summary] to be enabled for your project.

## Enable API Services needed by your cluster

You will need to enable the following API services in your project:

- Compute Engine API (`compute.googleapis.com`)
- Google Cloud APIs (`cloudapis.googleapis.com`)
- Cloud Resource Manager API (`cloudresourcemanager.googleapis.com`)
- Google DNS API (`dns.googleapis.com`)
- Identity and Access Management (IAM) API (`iam.googleapis.com`)
- IAM Service Account Credentials API (`iamcredentials.googleapis.com`)
- Service Management API (`servicemanagement.googleapis.com`)
- Service Usage API (`serviceusage.googleapis.com`)
- Google Cloud Storage JSON API (`storage-api.googleapis.com`)
- Cloud Storage (`storage-component.googleapis.com`)

You can enable these services using the console or the CLI (console service names in parentheses)

More information:
- [GCP: Enable Services][enable-svc]


[enable-svc]: https://cloud.google.com/service-usage/docs/enable-disable#enabling
[service-apis-summary]: https://cloud.google.com/terms/services
