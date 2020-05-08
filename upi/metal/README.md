# UPI Bare Metal

## Pre-requisites

### Matchbox

Setup an instance of matchbox using this [guide][coreos-matchbox-getting-started]. The matchbox instance must be reachable from public internet.

Store the tls assets that will be used for client authentication and the CA certificate that allows client to trust the matchbox server on the host that will run the example terraform scripts.

### AWS

Create a public route53 zone using this [guide][aws-create-public-route53-zone].

Setup `default` AWS cli profile on the host that will run the example terraform scripts using this [guide][aws-cli-configure-creds]

### Packet

Setup a Project in Packet.net that will be used to deploy servers, for example using this [guide][packet-deploy-server]

Setup API keys for your Project in Packet.net using this [guide][packet-api-keys]

Store the API keys in `PACKET_AUTH_TOKEN` so that `terraform-provide-packet` can use it to deploy servers in the project. For more info see [this][terraform-provider-packet-auth]

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
[packet-api-keys]: https://www.packet.com/developers/changelog/project-only-api-keys/
[packet-deploy-server]: https://support.packet.com/kb/articles/deploy-a-server
[terraform-examples]: https://github.com/s-urbaniak/terraform-examples#terraform-examples
[terraform-getting-started]: https://learn.hashicorp.com/terraform/getting-started/install.html
[terraform-provider-packet-auth]: https://www.terraform.io/docs/providers/packet/index.html#auth_token
