# GCP Service Account

Before proceeding with the OpenShift install, you should create a secondary IAM service account following the steps
outlined here:


## Step 1: Create a Service Account

Create a GCP service account using the console or the CLI

[GCP: Creating a service account][sa-create]

## Step 2: Assign Project Roles to your Service Account

You need to assign the proper roles to the newly created service account so it can be used to create an OpenShift cluster.
The minimum set of [roles][gcp-roles] you can assign the service account are the following:
- Compute Admin
- DNS Administrator
- Security Admin
- Service Account Admin
- Service Account User
- Storage Admin

If you want the OpenShift cluster to create new limited credentials for its own operators, you also need the following
role:
- Service Account Key Admin

Optionally, you can just assign the "Owner" role to your service account.

To assign roles to your service account you may use the console or the CLI:

[GCP: Assign service account roles][sa-assign]

## Step 3: Create and save a Service Account Key

You will need to create and save a service account key for your service account so you can use it with the OpenShift
Installer.

[GCP: Creating a service account key][sa-key]


## Step 4: Enable API Services needed by your cluster

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

[GCP: Enable Services][enable-svc]


[sa-create]: https://cloud.google.com/iam/docs/creating-managing-service-accounts#creating_a_service_account
[gcp-roles]: https://cloud.google.com/iam/docs/understanding-roles#predefined_roles
[sa-assign]: https://cloud.google.com/iam/docs/granting-roles-to-service-accounts#granting_access_to_a_service_account_for_a_resource
[sa-key]: https://cloud.google.com/iam/docs/creating-managing-service-account-keys#creating_service_account_keys
[enable-svc]: https://cloud.google.com/service-usage/docs/enable-disable#enabling
