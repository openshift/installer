# AWS: Installation requirements

The following tools and access rights are required to use Tectonic Installer with and Amazon Web Services (AWS) account.

* Tectonic License and Pull Secret
* A public AWS Route 53 Hosted Zone identifier. Public Route 53 DNS resolution is a requirement for controller-worker TLS communication. Choose a domain or subdomain and [configure it for name service at Route 53][aws-r53-doc]. Tectonic will create 2 subdomains in this Hosted Zone during provisioning.
* An EC2 Region and Availability Zone
* An EC2 SSH Key pair in that region
* An AWS Access Key and Secret (or a temporary Access Key, Secret, and Session Token)
* Two (1 Controller, 1 Worker) t2.medium nodes (minimum).
* Access to a minimum of 30GB of storage for each node.
* A current version of the Google Chrome or Mozilla Firefox web browser to run Tectonic Installer.

## Privileges

The AWS credentials you provide require access to the following AWS services:

* Route 53
* EC2
* ELB
* S3
* Security Groups
* VPC

An importable AWS policy containing the minimum privileges needed to run the Tectonic installer can be found [here][tectonic-installer-aws-policy].

## Temporary credentials

The following steps demonstrate how to generate and use temporary AWS credentials in conjunction with the Tectonic Installer:

1. Ensure the AWS CLI tool is installed by following the instructions on the [AWS CLI documentation][aws-cli-doc]. On Fedora, this can be done with `dnf install`:
```bash
$ sudo dnf install awscli
```

2. Ensure the AWS CLI is configured to use your access key ID and secret access key:
```bash
$ aws configure
```

3. Create a `tectonic-installer` role in AWS with the trust policy detailed [here][aws-trust-policy]. The trust relationship policy grants an entity permission to assume the role.
```bash
$ aws iam create-role --role-name tectonic-installer --assume-role-policy-document file://Documentation/files/aws-sts-trust-policy.json
```

The `file://` prefix is required before the filepath.

4. Add an inline AWS policy document to the `tectonic-installer` role containing the minimum privileges needed to run the Tectonic installer. The policy is available [here][tectonic-installer-aws-policy].
```bash
$ aws iam put-role-policy --role-name tectonic-installer --policy-name TectonicInstallerPolicy --policy-document file://Documentation/files/aws-policy.json
```

5. Add your user's ARN, found on the IAM user detail page, to the trusted entities for the `tectonic-installer` role. To do so, click on the `Trust Relationships` tab and then on the `Edit Trust Relationship` button to bring up the trusted entities JSON editor. You'll then need to add a new section for your user's ARN.

The example Trust Relationship below has been edited to add a user's (named tectonic) ARN:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com",
        "AWS": "arn:aws:iam::477645798577:user/tectonic"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```

6. Assume the `tectonic-installer` role with your AWS user using the AWS CLI tool as follows:
```bash
$ aws sts assume-role --role-arn=<TECTONIC_INSTALLER_ROLE_ARN> --role-session-name=tectonic-installer --role-session-name=<DESIRED_USER_NAME>
```
   The returned response will look like:
```json
{
  "Credentials": {
    "SecretAccessKey": "<SECRET_ACCESS_KEY>",
    "AccessKeyId": "<ACCESS_KEY_ID>",
    "Expiration": "2016-12-14T02:21:37Z",
    "SessionToken": "<SESSION_TOKEN>"
  },...
}
```

Use the `SECRET_ACCESS_KEY`, `ACCESS_KEY_ID`, and `SESSION_TOKEN` to authenticate in the installer.

If building the Tectonic cluster using the CLI directly, then you can configure `terraform` to perform the STS `assume-role` operation automatically on every run. It will automatically retrieve and use the temporary credentials every time so you don't have to refresh them manually when they expire.

To enable Terraform to perform the `assume-role` operation edit the file `platforms/aws/main.tf` and change the `provider "aws" { ... }` block to include the following configuration:
```hcl
provider "aws" {
  region = "${var.tectonic_aws_region}"
  assume_role {
    role_arn = "<tectonic-installer-ROLE-ARN>"
    session_name = "terraform"
  }
}
```
You can then run Terraform using an unpriviledged user that only has permissions to assume the `tectonic-installer` role.

## SSH key

The final step of the Tectonic install requires an SSH key and access to standard utilities like `ssh` and `scp`. Setting up a new key on AWS should take less than 5 minutes.

Tectonic uses AWS S3 to store all credentials, using server-side AES encryption for storage, and TLS encryption for upload/download. Any pod run in the system can query the AWS metadata, get node AWS credentials, and pull down cluster credentials from AWS S3. CoreOS plans to address this issue in a later release.

First, create a key.

1. Open a new terminal. Check if you already have a key by running `ls ~/.ssh/`. If you've previously created a key, you may see a file like `id_rsa.pub`. If you'd like to use this key, skip to *upload the key to AWS* below.
2. Type `ssh-keygen --help` to validate you have the openssh utilities installed. If you cannot find the binaries on your system, please consult your distro's documentation.
3. Type `ssh-keygen -t rsa -b 4096 -C "aws tectonic for alice@example.com"`. The content after `-C` is a comment. Replace alice@example.com with an appropriate AWS email or IAM account.
4. Follow the prompts on screen to finish creating your keypair. If you chose the default file name and location, your key should be in `$HOME/.ssh/id_rsa.pub`. Otherwise, the key-pair is in your current directory.

Next, upload the key to AWS.

1. Sign in using your IAM user or temporary credentials.
2. Go to *Services > Compute > EC2*.
3. Use the pulldown menu to select the same region as that selected for Tectonic installation.
4. On the left navigation under *Network & Security*, click *Key Pairs*.
5. Click *Import Key Pair*. Follow the displayed instructions to import your public key file, whose name should end in `.pub`.

For additional information about AWS and SSH keys consult the [official AWS guide][aws-ssh-key-doc].

## Access

In order to access the cluster two ELB backed services are exposed. Both are accessible over the standard TLS port (443).

## Install Tectonic

With temporary credentials and an SSH key, you'll be ready to install Tectonic. Head over to the [install doc][install-aws] to get started.

## Subnet/VPC requirements

The following table includes the high level networking features required to install Tectonic into new or existing VPCs, with or without public access to cluster services.

|                  | **Public facing cluster**          | **Internal cluster**                         |
-------------------|------------------------------------|----------------------------------------------|
| **New VPC**      | Installer creates public subnets | Select 'internal' in Tectonic installer      |
| **Existing VPC** | 2 subnets, connected to an IGW | Create 2 subnets, Establish a VPN |

### Configuring a *public* cluster

* Subnets for Controllers must have an attached and routed Internet Gateway.
* Subnets for Workers must be able to route requests to the Controller subnets and must have an associated route table that specifies a default gateway.
* The route tables should be explicitly attached to their subnets.

### Configuring an *internal* cluster

* Subnets for Controllers and Workers must be able to route requests to each other and must have an associated route table that specifies a default gateway.
* The route tables should be explicitly attached to their subnets.
* You must have VPN access to the subnet as it is does not offer an inbound connection to the Internet.

### Using an existing VPC

By default, Tectonic Installer creates a new AWS Virtual Private Cloud (VPC) for each cluster. Advanced users can choose to use an existing VPC instead. An existing VPC must have an [Internet Gateway][aws-vpc-inet-gateway]. Tectonic Installer will not create an Internet Gateway in an existing VPC.

An existing VPC for a public cluster must have a public subnet for controllers, and a private subnet for workers. An existing VPC for an internal cluster must have 2 private subnets, one each for controllers and workers.

Public subnets have a default route to the Internet Gateway and should auto-assign IP addresses. Private subnets have a default route to a default gateway, such as a NAT Gateway or a Virtual Private Gateway.

*DHCP Options Set* attached to the VPC must have an AWS [private domain name][aws-vpc-dns-hostnames]. In us-east-1 region, an AWS private domain name is ec2.internal whereas other regions use region.compute.internal.

When using an existing VPC, tag AWS VPC subnets with the `kubernetes.io/cluster/my-cluster-name = shared` tag. `shared` is used to tag resources shared between multiple clusters, which should not be destroyed if any individual cluster is destroyed. If this tag is not specified, AWS ELB integration with Tectonic may not be able to use VPC subnets.


[aws-cli-doc]: http://docs.aws.amazon.com/cli/latest/userguide/installing.html
[aws-r53-doc]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/creating-migrating.html
[aws-ssh-key-doc]: http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/get-set-up-for-amazon-ec2.html
[aws-trust-policy]: ../../files/aws-sts-trust-policy.json
[aws-vpc-inet-gateway]: https://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_Internet_Gateway.html
[install-aws]: index.md
[tectonic-installer-aws-policy]: ../../files/aws-policy.json
[aws-vpc-dns-hostnames]: http://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/vpc-dns.html#vpc-dns-hostnames
