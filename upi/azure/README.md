# Red Hat OpenShift Container 4.x UPI Platform on Azure

## Overview
This provides a Azure ARM Template to launch user provided infrastructure implemenation of OpenShift 4.x on Azure.

This creates:
 
    A Resource Group
    3 Masters
    3-16 Workers
    A API Loadbalancer
    A Application Loadbalancer
    2 Availablity Groups

### To Use:

### Prerequisites:
1. Clone the installer repo: git clone https://github.com/openshift/installer
1. Download the installer and place it in the upi/azure directory
2. Make sure python3 is installed 
3. Execute: pip install dotmap 
4. Using the provided example(install-config.yaml.example), 
   create a install-config.yaml, including your  
   ssh key, and the pull secret. 
5. Choose a unique name for the resource group, and use it in the commands below, instead of [UniqueResourceGroupName]


./setup_azarm.sh [UniqueResourceGroupName] 
./deploy_azarm.sh [UniqueResourceGroupName]   
 

### Customizatiion:
1. To change the version of OpenShift: 
	a. Change the RHCOS image url in setup_azarm.sh
        b. Download the installer for the desired version (Note 4.3 is the minimum version supported for Azure UPI ARM)
2. Change the VM size via the azuredeploy.parameters.json 
3. To Fork the repo, for customizations: The azuredeploy.json contains the user and project name, as well as the branch. Please change these
items if you fork this project. 

