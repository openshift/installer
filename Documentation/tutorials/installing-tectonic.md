# Downloading and installing Tectonic on AWS

Once your AWS account is activated, you can create your CoreOS Tectonic account, and prepare your AWS account for installation. This tutorial will cover:

* Creating your Tectonic account
* Preparing your AWS account for installation
* Installing Tectonic using AWS
* Creating a new Kubernetes cluster

## Creating your Tectonic account

To create your CoreOS Tectonic account:

1. Go to https://coreos.com/tectonic.
2. Click *Get Tectonic for Kubernetes*.
3. Create an account using either your Google account or another email address.
4. Enter your contact information, and click *Get License* for 10 nodes.
5. Agree to the license terms.

You will see notice that your order is complete, with a link to your account Overview page.  

The Overview page will list your license type, purchase date, and status, and provide links to the Tectonic documentation. It will also provide your Universal Software License, Docker Config (dockercfg) pull secret, and Kubernetes Secret.

Your Tectonic account has now been created!

## Preparing your AWS account for installation

After activating your Tectonic account, review and complete [Creating an AWS account][creating-aws] before downloading the Tectonic Installer.

Installing Tectonic requires:
* An Identity Access Management (IAM) account
* An associated Access Key
* A domain or subdomain with DNS name service at AWS Route 53
* An AWS Route 53 EC2 SSH key pair

Tectonic will create a new AWS Virtual Private Cloud (VPC), or you can select an existing VPC. To use an existing VPC, see the [existing VPC requirements][vpc-req].

## Downloading and running the Tectonic Installer

Having completed the AWS installation requirements, you are now ready to download and run the Tectonic Installer.

Make sure a current version of either the Google Chrome or Mozilla Firefox web browser is set as the default browser on the workstation where the installer will run.

Download and run the Tectonic installer by opening a new terminal, and running the following commands:

Pull down the tarball with:
```
$ curl -O https://releases.tectonic.com/tectonic-1.5.6-tectonic.1.tar.gz
```

Extract the tarball with:
```
<<<<<<< c06638378d3b1d6054f11bef2c8516cf54a25507
tar xzvf tectonic-1.5.6-tectonic.1.tar.gz
=======
$ tar xzvf tectonic-1.5.5-tectonic.3.tar.gz
>>>>>>> docs:tutorials update
```

Change to the previously untarred directory with:
```
$ cd tectonic/tectonic-installer
```

Run the Tectonic Installer:
```
$ ./$PLATFORM/installer
```
Where `$PLATFORM` is one of: `linux`, `darwin` or `windows`.

A browser window will open to begin the GUI installation process.

Use the `./$PLATFORM/installer` command to relaunch the Installer at any time. When launched, you will be given the option to *Start Over*, or to *Continue* where you left off.

![Installer Pop-up](https://coreos.com/tectonic/docs/latest/img/installer-aws.png)

Use the Tectonic Installer wizard to provision your cluster. (This should take about 10-15 minutes. If you encounter any errors, please see the [AWS: Troubleshooting Installation][aws-troubleshooting] guide.)

Once complete, access the Tectonic Console through a browser window.

## Creating a new Kubernetes cluster

With Tectonic Installer running locally, deploy the Tectonic Kubernetes distribution on a new cluster.

### Choose Cluster Type

**Platform**

Use the pulldown menu to select the platform on which the cluster will be installed.

(This installation will use Amazon Web Services as its Platform.)

### Define Cluster

Define AWS credentials and configuration options for your cluster.

**AWS Credentials**

To allow Tectonic to communicate with your AWS account, provide your AWS credentials.

Select Use a normal access key, or Use a temporary session token. 	

* *Normal access key:* Copy and paste the Access Key ID and Secret Access Keys created earlier in the AWS setup process.
* *Use a temporary session token:* CoreOS recommends that you use a temporary session token to generate temporary credentials, and protect the integrity of your main AWS credentials. Enter the Access Key Id and Secret Access Keys created during the AWS setup process.
* *Region:* Enter the EC2 region selected during your AWS setup.

Your Access Key ID is available from the AWS console. Your Secret Access Key is available from the CSV file downloaded when creating the Access Key. See [Creating an AWS account][creating-aws] for more information.

**Cluster Info**

Next, define the following attributes for your cluster:

* *CloudFormation Stack Name:* Name your cluster. This name will appear in the Tectonic Console.
* *Container Linux Update Channel:* Select the channel for your update mechanism for Container Linux (Stable, Beta or Alpha).
* *Experimental Features:* Check this box to enable operators to perform automatic updates.
* *CoreOS License and Pull Secret:* Copy and paste your Universal Software License, and your dockercfg pull secret from your Tectonic account page. (https://account.tectonic.com)
* *CloudFormation Tags:* Enter any key-value pairs you wish to add to the Cluster as CloudFormation tags.
* *DNS:* Select a pre-configured Route 53 hosted zone, and enter a subdomain.
* *Internal Cluster:* Select this checkbox to make this cluster’s EC2 instances and Elastic Load Balancers internal only.

**Certificate Authority**

Select the option to allow Tectonic to generate a Certificate Authority and key for you.

Provide a CA certificate and key in PEM format if you are managing your own PKI.

**Submit Keys**

Select your SSH Key from the pulldown list, and click *Generate a New Key* to create a KMS key pair.

Be certain to select the SSH key you submitted while setting up your AWS EC2 Network and Security keys.

**Console Login**

Enter the credentials that will be used to log in to Tectonic Console.

Click *Submit* to submit your assets and create your Kubernetes cluster. (Cluster creation may take up to 20 minutes.)

If you hit permissions errors during the creation process it is likely that your IAM account does not have sufficient privileges. Review the privileges section of our AWS: Installation Requirements to get your IAM account configured correctly.

### Boot Cluster

The final step in creating your Kubernetes cluster is to boot your cluster.

**Start Kubernetes**

The Start Kubernetes screen displays cluster creation in process.

When Starting AWS CloudFormation and Starting Kubelet are complete, download the created assets, then open a new terminal window to send the final setup commands to your cluster node.

1. Click *Download Assets* to save your cluster assets locally.
2. Open a terminal to run the commands shown on the installer's penultimate screen.
3. Run
   `ssh -i <path/keyfilename> <core@IP-ADDRESS> sudo systemctl start bootkube`
   Where `<path/keyfilename>` is the local path to your SSH key, and `<core@IP-ADDRESS>` is the IP address specified in the Tectonic Installer.
4. Then run
   `ssh <core@IP-ADDRESS> journalctl -u bootkube -f`
   Using the same `<core@IP-ADDRESS>` listed on your screen.

When Kubernetes is fully launched, click *Next Step* to continue.

Use the email address and password you used to create your Tectonic account to log in to the Console.

[**NEXT:** Deploying an application on Tectonic][first-app]


[install-req]: ../install/aws/requirements.md
[ssh-key]: ../install/aws/requirements.md#ssh-key
[vpc-req]: ../install/aws/requirements.md#using-an-existing-vpc
[trouble-shoot]: ../install/aws/troubleshooting.md
[privileges]: ../install/aws/requirements.md#privileges
[first-app]: first-app.md
[creating-aws]: creating-aws.md
[aws-troubleshooting]: ../install/aws/troubleshooting.md
