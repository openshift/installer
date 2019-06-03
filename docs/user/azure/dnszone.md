# Public DNS Zone

Azure's DNS Zone service is used by the OpenShift installer to configure cluster DNS resolution and provide name lookup
for the cluster to the outside world. To use OpenShift, you must have created a public DNS zone in Azure in
the same subscription as your OpenShift cluster. You must also ensure the zone is "authoritative" for the domain. There are
two ways to do this outlined below: root domain and subdomain. A root domain is `example.com`. A subdomain is of
the form `clusters.example.com`.

The below sections identify how to ensure your hosted zone is authoritative for a domain.

## Step 1: Acquire/Identify Domain

You may skip this step if using an existing domain and registrar. You will move the authoritative DNS to Azure or
submit a delegation request for a subdomain in a later step.

Azure can also purchase domains for you and act as a registrar. If you allow Azure to purchase a new domain for you,
you can skip the remainder of these steps (the domain is created and the hosted zone is created correctly for you)!

[Documentation][buy-domain-from-azure] on buying domain for Azure

## Step 2: Create Public Hosted Zone

Whether using a root domain or a subdomain, you must create a public, hosted zone.

[Azure: Creating a Public DNS Zone][create-hosted-zone]

To use the root domain, you'd create the hosted zone with the value `example.com`. To use a subdomain, you'd
create a hosted zone with the value `clusters.example.com`. (Use appropriate domain values for your situation.)

[buy-domain-from-azure]: https://docs.microsoft.com/en-us/azure/app-service/manage-custom-dns-buy-domain
[create-hosted-zone]: https://docs.microsoft.com/en-us/azure/dns/dns-delegate-domain-azure-dns
