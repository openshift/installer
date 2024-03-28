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

These properties are used to control how the script works.  All of the properties will dictate how much or little work the script will do.

| Property            | Description                                                                                                                                                                                                                                                                                       |
|---------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| createInstallConfig | Enable script to create install config based on configuration of variables.ps1.  If set to false, the script will not generate install-config.yaml.  Note, install-config.yaml will be needed if you want this script to generate the ignition files for vm creation.                             |
| downloadInstaller   | Enable script to download installer to be used.  If not downloading the installer, the installer must be placed in the same directory as this script.  The installer is needed to determine what template to upload to the cluster for vm cloning.                                                |
| uploadTemplateOva   | Enable script to upload OVA template to be used for all VM being created.                                                                                                                                                                                                                         |
| generateIgnitions   | Enable script to generate ignition configs.  This is normally used when install-config.yaml is provided to script, but need script to generate the ignition configs for VMs.                                                                                                                      |
| waitForComplete     | This option has the script step through the process of waiting for installation complete.  Most of this functionality is provided by `openshift-install wait-for`.  The script will will check for when api is ready, bootstrap complete, accept CSRs and then for all COs to be done installing. |
| delayVMStart        | This option has the script delay the start of the VMs after their creation.                                                                                                                                                                                                                       |


### vCenter Configuration

These properties are related to vCenter.  These will be used by script to create the install-config.yaml as well as define a single failure domain if no failure domains are configured in the variables.ps1 file.

| Property    | Description                                                                                                                                    |
|-------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| vcenter     | This property specifies the vCenter URL                                                                                                        |
| username    | Specifies the username used to connect to vCenter.  This is only used in generating the install-config.yaml information for use by the cluster |
| password    | The password to be used to connect to vCenter. This is only used in generating the install-config.yaml information for use by the cluster      |
| portgroup   | The default port group to use.  This property will be used when no failure domains are defined.                                                |
| datastore   | The default data store to use.  This property will be used when no failure domains are defined.                                                |
| datacenter  | The default data center to use.  This property will be used when no failure domains are defined.                                               |


### Cluster Configuration

These properties are used to configure the new cluster.  These properties will also be used to determine how to create the VMs.

| Property        | Description                                                                                                                                           |
|-----------------|-------------------------------------------------------------------------------------------------------------------------------------------------------|
| clustername     | The name of the cluster to create.                                                                                                                    |
| basedomain      | The base domain to use for the cluster.                                                                                                               |
| sshkeypath      | The path of the ssh key to use when generating ignition configs.                                                                                      |
| failure_domains | Contains the failure domains for the cluster.  This needs to be a JSON formatted configuration following what is provided in the example below table. |

Example of a failure domain configuration:

```powershell
$failure_domains = @"
[
    {
        "datacenter": "IBMCloud",
        "cluster": "vcs-mdcnc-workload-1",
        "datastore": "mdcnc-ds-1",
        "network": "ocp-ci-seg-14"
    },{
        "datacenter": "IBMCloud",
        "cluster": "vcs-mdcnc-workload-2",
        "datastore": "mdcnc-ds-2",
        "network": "ocp-ci-seg-14"
    },{
        "datacenter": "IBMCloud",
        "cluster": "vcs-mdcnc-workload-3",
        "datastore": "mdcnc-ds-3",
        "network": "ocp-ci-seg-14"
    },{
        "datacenter": "datacenter-2",
        "cluster": "vcs-mdcnc-workload-4",
        "datastore": "mdcnc-ds-4",
        "network": "ocp-ci-seg-14"
    }
]
"@
```

### VM Configuration

| Property                   | Description                                                                                                                                    |
|----------------------------|------------------------------------------------------------------------------------------------------------------------------------------------|
| vm_template                | Name of the existing vm template in vCenter to use.  Use this option to prevent script from uploading a template and use an existing template. |
| dns                        | DNS server for network configuration of all VMs.                                                                                               |
| gateway                    | Gateway for network configuration of all VMs.                                                                                                  |
| netwask                    | The network mask for all VMs.                                                                                                                  |
| lb_ip_address              | The IP address to use for the load balancer VM.                                                                                                |
| bootstrap_ip_address       | The IP address to use for the bootstrap VM.                                                                                                    |
| control_plane_memory       | Amount of memory to assign each control plane VM.  This value is in MB.                                                                        |
| control_plane_num_cpus     | Number of CPUs to assign to each control plane VM.                                                                                             |
| control_plane_count        | Number of control plane VMs.                                                                                                                   |
| control_plane_ip_addresses | The IP addresses to assign all control plane VMs.                                                                                              |
| control_plane_hostnames    | The host names to assign to all control plane VMs.                                                                                             |
| compute_memory             | Amount of memory to assign each compute VM.  This value is in MB.                                                                              |
| compute_num_cpus           | Number of CPUs to assign to each compute VM.                                                                                                   |
| compute_count              | Number of compute VMs.                                                                                                                         |
| compute_ip_addresses       | The IP addresses to assign all compute VMs.                                                                                                    |
| compute_hostnames          | The host names to assign to all compute VMs.                                                                                                   |


## Build a Cluster

### Manual Method

1. Create an install-config.yaml.  The machine CIDR for the dev cluster is 139.178.89.192/26.

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

2. Run `openshift-install create manifest`.

3. Update any configurations you need before generating ignition configs.  It is recommended to remove the master machine CR and the worker machineset CR.  This can be accomplished by running:

```
rm -f openshift/99_openshift-cluster-api_master-machines-*.yaml openshift/99_openshift-cluster-api_worker-machineset-*.yaml
```

3. Run `openshift-install create ignition-configs`.

4. Create and configure a variables.ps1 file.
   There is an example variables.ps1 file in this directory named variables.ps1.example. The example file contains all properties that are available to be configured. At a minimum, you need to set values for the following variables.
* clustername
* basedomain
* username
* password
* vcenter
* vcentercredpath

The bootstrap ignition config must be placed in a location that will be accessible by the bootstrap machine. For example, you could store the bootstrap ignition config in a gist.

Even if declaring static IPs a DHCP server is still required early in the boot process to download the ignition files.

5. Run `pwsh -f upi.ps1`.

The script will now attempt to create all VMs based on the config provided.  The script will start the VMs and will say `Install Comlete` when all VMs are cloned and started.

6. Run `openshift-install wait-for install-complete`. Wait for the cluster install to finish.

7. Enjoy your new OpenShift cluster.

8. Run `pwsh -f upi-destroy.ps1` to clean up all folder, VMs, and tags.

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
