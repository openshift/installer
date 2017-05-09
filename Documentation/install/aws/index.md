# AWS: Installation

## Preparation

Check the [requirements doc][install-aws-requirements] to see what's needed. The short version:

* an IAM account
* an [associated SSH key][ssh-key-req]
* a domain or subdomain with [DNS name service at AWS Route53][aws-r53-doc].
* Tectonic will create a new AWS Virtual Private Cloud (VPC), or you can select an existing VPC. To use an existing VPC, see the [existing VPC requirements][install-aws-requirements-evpc].

## Step 1: Download and run Tectonic Installer

Make sure a current version of either the Google Chrome or Mozilla Firefox web browser is set as the default browser on the workstation where Installer will run.

Download the [Tectonic installer][latest-tectonic-release].

```bash
wget https://releases.tectonic.com/tectonic-1.6.2-tectonic.1.tar.gz
tar xzvf tectonic-1.6.2-tectonic.1.tar.gz
cd tectonic
```

Run the Tectonic Installer.

```bash
./tectonic-installer/$PLATFORM/installer
```

Where $PLATFORM is one of: `linux`, `darwin` or `windows`.

A browser window should open to begin the GUI installation process.

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="../../img/installer-aws.png">
    <div class="co-m-screenshot-caption">Entering AWS credentials in the Tectonic Installer</div>
  </div>
</div>

## Step 2: Install

Be sure to read the [installation requirements][install-aws-requirements], which include a section on [privileges for your AWS credentials][install-aws-requirements-creds], as well as the [known issues section in the Troubleshooting guide][install-aws-troubleshooting] before you install.

Follow the on-screen instructions to provision your cluster. This process should take about 10-15 minutes.

Once complete you should be able to access your Tectonic Console.

If you encounter any errors check the [troubleshooting][install-aws-troubleshooting] documentation.

## Step 3: Tectonic Console

<div class="row">
  <div class="col-lg-8 col-lg-offset-2 col-md-10 col-md-offset-1 col-sm-12 col-xs-12 co-m-screenshot">
    <img src="../../img/prometheus-monitoring.png">
    <div class="co-m-screenshot-caption">Viewing a node in the Tectonic Console</div>
  </div>
</div>

Now you are ready to access Tectonic Console, configure `kubectl`, and deploy your first application to the cluster. A `kubeconfig` with the appropriate configuration for `kubectl` is available for download in the Tectonic Console.

For those new to Tectonic and Kubernetes, the [usage guide][usage] contains tutorials and guides demonstrating how to get started.


[aws-r53-doc]: https://docs.aws.amazon.com/Route53/latest/DeveloperGuide/creating-migrating.html
[ssh-key-req]: requirements.md#ssh-key
[install-aws-requirements]: requirements.md
[install-aws-requirements-creds]: requirements.md#privileges
[install-aws-requirements-evpc]: requirements.md#using-an-existing-vpc
[usage]: ../../tutorials/index.md
[latest-tectonic-release]: https://releases.tectonic.com/tectonic-1.6.2-tectonic.1.tar.gz
[install-aws-troubleshooting]: ../../troubleshooting/faq.md
