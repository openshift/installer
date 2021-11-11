# Equinix Metal UPI

## Pre-requisites

### Matchbox

Setup an instance of matchbox using this [guide][coreos-matchbox-getting-started]. The matchbox instance must be reachable from public internet.

Store the tls assets that will be used for client authentication and the CA certificate that allows client to trust the matchbox server on the host that will run the example terraform scripts.

### AWS

Create a public route53 zone using this [guide][aws-create-public-route53-zone].

Setup `default` AWS cli profile on the host that will run the example terraform scripts using this [guide][aws-cli-configure-creds]

### Equinix Metal

Some portions of this guide refer to "Packet" variables and tools. [Packet was
acquired by Equinix][acquisition] and the features became available as [Equinix
Metal][introducing-equinix-metal] in 2020.

Setup a Project in [Equinix Metal][equinix-metal] that will be used to deploy servers, for example using this [guide][metal-deploy-server].

Setup API keys for your project in Equinix Metal using this [guide][metal-api-keys]

Store the API keys in `PACKET_AUTH_TOKEN` so that `terraform-provider-packet`
can use it to deploy servers in the project. For more info see
[this][terraform-provider-packet-auth]

#### Terraform

Install Terraform on the host that will run the example terraform scripts using the Getting Started [Guide][terraform-getting-started]

## Generating the terraform variable example

Install [terraform-examples][terraform-examples] on your host.

```sh
go get -u github.com/s-urbaniak/terraform-examples
```

Generate the `terraform.tfvars.example` using

```sh
terraform-examples config.tf > terraform.tfvars.example
```

[aws-cli-configure-creds]: https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html
[aws-create-public-route53-zone]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/CreatingHostedZone.html
[coreos-matchbox-getting-started]: https://matchbox.psdn.io/getting-started/
[metal-api-keys]: https://metal.equinix.com/developers/docs/accounts/users/#api-keys
[metal-deploy-server]: https://metal.equinix.com/developers/docs/deploy/on-demand/
[terraform-examples]: https://github.com/s-urbaniak/terraform-examples#terraform-examples
[terraform-getting-started]: https://learn.hashicorp.com/terraform/getting-started/install.html
[terraform-provider-packet-auth]: https://registry.terraform.io/providers/packethost/packet/latest/docs#auth_token
[acquisition]: https://www.equinix.com/newsroom/press-releases/2020/03/equinix-completes-acquisition-of-bare-metal-leader-packet
[introducing-equinix-metal]: https://blog.equinix.com/blog/2020/10/06/equinix-metal-metal-and-more/
[equinix-metal]: https://metal.equinix.com