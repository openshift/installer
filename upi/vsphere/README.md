This project shows two ways to install an UPI cluster.  We will discuss how to install using one of these two techniques:
- Terraform
- PowerShell

# Table of Contents
- [PowerShell](#PowerShell)
  - [Pre-Requisites](#pre-requisites)
  - [PowerShell Setup](#powershell-setup)
    - [VMware.PowerCLI](#vmwarepowercli)
    - [EPS](#eps)
  - [Script Configuration]
  - [OpenShift Installation Configuration]()
- [Terraform](#Terraform)
  - [Pre-Requisites](#pre-requisites-1)
  - [Build a Cluster](#build-a-cluster-1)

# PowerShell
This section will describe the process to generate the vSphere VMs using PowerShell and the supplied scripts in this module.

## Pre-requisites
* PowerShell
* PowerShell VMware.PowerCLI Module
* PowerShell EPS Module

## PowerShell Setup

PowerShell will need to have a couple of plugin installed in order for our script to work.  The plugins we need to install are VMware.PowerCLI and EPS.

### VMware.PowerCLI

To install the VMware.PowerCLI, you can run the following command:

```shell
pwsh -Command 'Install-Module VMware.PowerCLI -Force -Scope CurrentUser'
```

### EPS

To install the EPS module, you can run the following command:

```shell
pwsh -Command 'Install-Module -Name EPS -RequiredVersion 1.0 -Force -Scope CurrentUser'
```

### Generating CLI Credentials

The PowerShell scripts require that a credentials file be generated with the credentials to be used for generating the vSphere resources.  This does not have to be the credentials used by the OpenShift cluster via the install-config.yaml, but must have all permissions to create folders, tags, templates, and vms.  To generate the credentials files, run:

```shell
pwsh -command "\$User='<username>';\$Password=ConvertTo-SecureString -String '<password>' -AsPlainText -Force;\$Credential = New-Object -TypeName System.Management.Automation.PSCredential -ArgumentList \$User, \$Password;\$Credential | Export-Clixml secrets/vcenter-creds.xml"
```

Be sure to modify `<username>` to be the username for vCenter and `<password>` to the your password.  The output of this needs to go into `secrets/vcenter-creds.xml`.  Make sure the secrets directory exists before running the credentials generation command above.

## Script Configuration

The PowerShell script provided by this project provides examples on how to do several aspects to creating a UPI cluster environment.  It is configurable to do as much or as little as you need.  For the CI build process, we will handle all install-config.yaml configuration, uploading of templates, and monitoring of installation progress.  This project can handle doing all that as well if configured appropriately.

### Behavioral Configurations

| Property            | Description                                                                                                                                                                                                                                                                                        |
|---------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| createInstallConfig | Enable script to create install config based on configuration of variables.ps1                                                                                                                                                                                                                     |
| downloadInstaller   | Enable script to download installer to be used.  If not downloading the installer, the installer must be placed in the same directory as this script.                                                                                                                                              |
| uploadTemplateOva   | Enable script to upload OVA template to be used for all VM being created.                                                                                                                                                                                                                          |
| generateIgnitions   | Enable script to generate ignition configs.  This is normally used when install-config.yaml is provided to script, but need script to generate the ignition configs for VMs.                                                                                                                       |
| waitForComplete     | This option has the script step through the process of waiting for installation complete.  Most of this functionality is provided by `openshift-install wait-for`.  The script will will check for when api is ready, bootstrap complete, accept CSRs and then for all COs to be done installing.  |
| delayVMStart        | This option has the script delay the start of the VMs after their creation.                                                                                                                                                                                                                        |

## Build a Cluster

# Terraform
This section will walk you through generating a cluster using Terraform.

<a id="terraform-pre-requisites"></a>
## Pre-Requisites

* terraform
* jq

## Build a Cluster

1. Create an install-config.yaml.
The machine CIDR for the dev cluster is 139.178.89.192/26.

```
apiVersion: v1
baseDomain: devcluster.openshift.com
metadata:
  name: mstaeble
networking:
  machineNetwork:
  - cidr: "139.178.89.192/26"
platform:
  vsphere:
    vCenter: vcsa.vmware.devcluster.openshift.com
    username: YOUR_VSPHERE_USER
    password: YOUR_VSPHERE_PASSWORD
    datacenter: dc1
    defaultDatastore: nvme-ds1
pullSecret: YOUR_PULL_SECRET
sshKey: YOUR_SSH_KEY
```

2. Run `openshift-install create ignition-configs`.

3. Fill out a terraform.tfvars file with the ignition configs generated.
There is an example terraform.tfvars file in this directory named terraform.tfvars.example. The example file is set up for use with the dev cluster running at vcsa.vmware.devcluster.openshift.com. At a minimum, you need to set values for the following variables.
* cluster_id
* cluster_domain
* vsphere_user
* vsphere_password
* ipam_token OR bootstrap_ip, control_plane_ips, and compute_ips

The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.

Even if declaring static IPs a DHCP server is still required early in the boot process to download the ignition files. 

4. Run `terraform init`.

5. Ensure that you have you AWS profile set and a region specified. The installation will use create AWS route53 resources for routing to the OpenShift cluster.

6. Run `terraform apply -auto-approve`.
This will reserve IP addresses for the VMs.

7. Run `openshift-install wait-for bootstrap-complete`. Wait for the bootstrapping to complete.

8. Run `terraform apply -auto-approve -var 'bootstrap_complete=true'`.
This will destroy the bootstrap VM.

9. Run `openshift-install wait-for install-complete`. Wait for the cluster install to finish.

10. Enjoy your new OpenShift cluster.

11. Run `terraform destroy -auto-approve`.
