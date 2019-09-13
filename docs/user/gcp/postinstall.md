# Service Account: Post Install

If the service account used to create the cluster was given the Owner role or included the Service Account Key Admin role, 
the service account no longer requires elevated permissions after install. You may change its role to Viewer or 
remove all roles bound to it. You can perform these steps by revisiting the service account role binding step you performed
earlier.

[GCP: Assign service account roles][sa-assign]

[sa-assign]: https://cloud.google.com/iam/docs/granting-roles-to-service-accounts#granting_access_to_a_service_account_for_a_resource
