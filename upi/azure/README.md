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
1. Download the installer and place it in the arm directory
2. Make sure python3 is installed 
3. Execute: pip install dotmap 


./setup_azarm.sh UniqueResourceGroupName  
read -p "Press [Enter] to start deploy"  
./deploy_azarm.sh UniqueResourceGroupNaem  

