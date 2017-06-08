# Internal Cluster

This directory contains Terraform configuration that provisions a VPC with a VPN connection in AWS. This setup is designed to emulate a customer-like deployment in order to end-to-end test deploying Tectonic as a private "Internal Cluster" to an "Existing VPC.

This Terraform configuration provisions the following AWS resources by default:
* 1 VPC with name configured by `TF_VAR_vpc_name`
* 4 subnets in the VPC with count configured by `TF_VAR_subnet_count`
* 1 public subnet containing an internet gateway and NAT gateway
* 1 private Route 53 zone for `tectonic.dev.coreos.systems` configured by `TF_VAR_base_domain`
* 1 t2.micro EC2 instance in the public subnet running [OpenVPN Access Server](https://aws.amazon.com/marketplace/pp/B00MI40CAE/ref=mkt_wir_openvpn_byol)
* 1 VPN gateway and VPN connection

## Usage

### Install Terraform

[Download the Terraform binary](https://www.terraform.io/downloads.html) and install it.


### Configure Credentials

Any existing credentials available in the `~/.aws/credentials` file will automatically be used. Otherwise, make the AWS credentials available by exporting the following environment variables:

```
export AWS_ACCESS_KEY_ID=<aws-key-id>
export AWS_SECRET_ACCESS_KEY=<aws-key-secret>
```

### Additional Variables

Terraform will prompt for any unset required variables. These variables can be manually entered at every run, exported as environment variables, or configured with a [terraform.tfvars](https://www.terraform.io/docs/configuration/variables.html#variable-files) file that will be ignored by git and used for every run. Simply create a `terraform.tfvars` file and set any required variables or overrides

### Running

Validate the configuration and plan the run with:

```
terraform plan
```

Provision the infrastructure with:

```
terraform apply
```

### Connect to the VPN

Once the infrastructure is ready, Terraform will output an `ovpn_url` variable containing the URL of the OpenVPN Access Server. In order to connect to the VPN, take the following steps:

1. Navigate to `ovpn_url` and login to the Access Server with the username `openvpn` and the password provided when running Terraform.
2. Download the OpenVPN configuration file from the Access Server.
3. Follow the instructions for the appropriate OS to setup a VPN connection using the configuration file.
4. When establishing the VPN connection, use the same credentials used when connecting to the Access Server. If prompted, do not provide a private key password.

### Manual DNS configuration

Terraform does not support changing SOA TTLs in Route 53. As a result, in order to make use of the private zone immediately, the TTLs must be manually modified in the AWS console.

### Tectonic Installation

Once all the infrastructure is provisioned and the VPN connection is available, a Tectonic clusyou can be installed the VPC. When running the Tectonic installer, be sure to:

* Select the provisioned private DNS Zone using the GUI or by setting the `TF_VAR_tectonic_aws_external_private_zone` environment variable.
* Install Tectonic in the provisioned VPC by selecting the "Existing VPC" option and selecting the appropriate VPC ID in the GUI or by setting the `TF_VAR_tectonic_aws_external_vpc_id` environment variable.


### Tear Down

To tear down the infrastructure or to restart the process, run:

```
terraform destroy
```

### Troubleshooting

If `terraform apply` fails, Terraform will not automatically roll back the created resources. Before attempting to create the infrastructure again, the resources must be destroyed manually by running:

```
terraform destroy
```
