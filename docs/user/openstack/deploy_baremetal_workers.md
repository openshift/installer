# Deploying OpenShift bare-metal workers on OpenStack cloud provider


## Table of Contents
- [Common prerequisites](#common-prerequisites)
- [Considerations whe deploying bare-metal workers](#considerations-whe-deploying-bare-metal-workers)
- [Deploying cluster with BM workers on tenant network deployed by the installer](#deploying-cluster-with-bm-workers-on-tenant-network-deployed-by-the-installer)
- [Deploying cluster with BM workers on preexisting network](#deploying-cluster-with-bm-workers-on-preexisting-network)
- [Deploy cluster with vm workers only and add BM workers after](#deploy-cluster-with-vm-workers-only-and-add-bm-workers-after)


## Common prerequisites

* [Ironic bare metal service][2] is enabled and accessible through the [OpenStack Compute API][3].
* Bare-metal machines are available as an [OpenStack flavor][1]
* At the time bare-metal machines are booted, the image used to boot them is available from the [Glance OpenStack image service][4].
If given a url, the installer will fetch that image from the location specified and upload it to Glance. See relevant
[documentation][5] and [clusterOSImage][6]. 


## Considerations whe deploying bare-metal workers

* Long boot times for bare-metal machines may result in deployment timeout

    The first time a bare-metal server is booted the boot time could be significantly longer that
    the first time a vm server is booted. This longer boot time could exceed the timeout settings of 
    the installer. If that occurs, the deployment will fail with a timeout error.  The deployment maybe
    restarted and completed by re-running  the installer with the appropriate *wait-for*
    command. For example:
        
        ./openshift-install wait-for install-complete --log-level debug 
          

## Deploying cluster with BM workers on tenant network deployed by the installer

The initial cluster is deployed using bare-metal workers only. Bare-metal workers and control-plane vms are all
attached to the tenant network provisioned by the installer.

- Requirements: 
    - By default, OpenStack networks support attaching both vms and bare-metal machines to them. 
    * [Ironic bare metal service][2] can listen for, and PXE-boot machines in tenant networks

- Create install-config.yaml:

    - Set compute.[worker].platform.openstack.type to the bare-metal server flavor.
    - Set control-plane.platform.openstack.type to the vm flavor which will be used by the control-plane masters.

        For example:
 
                controlPlane:
                   platform:
                     openstack:
                       type: <vmComputeFlavorForMasters>
                 
                 ... other settings
                 
                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <baremetalComputeFlavor>
                   replicas: 3
         
                 ... other settings
                 


- Run the openshift installer:

      ./openshift-install create cluster --log-level debug
        
- wait for the installer to complete. 

    If the installer times out waiting for the bare-metal workers to complete booting,
    restart the installation using the appropriate *wait-for* command:
    
        ./openshit-install watier-for install-complete --log-level debug

## Deploying cluster with BM workers on preexisting network

Initial cluster is deployed using bare-metal workers only. Bare-metal workers are attached to
a preexisting network. 

- Requirements:
    
    - An OpenStack subnet has been pre-provisioned. The subnet supports attaching both vms and bare-metal severs to it. 

- Create install-config.yaml:

    - Set compute.[worker].platform.openstack.type to the bare-metal server flavor.

    - Set control-plane.platform.openstack.type to the vm flavor which will be used by the control-plane masters.

    - Set platform.openstack.machinesSubnet to the UUID of the pre-provisioned subnet. 

        For example:
 
                controlPlane:
                   platform:
                     openstack:
                       type: <vmComputeFlavorForMasters>
                 
                 ... other settings
                 
                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <baremetalComputeFlavor>
                   replicas: 3
         
                 ... other settings
                 
                 platform:
                   openstack:
                     machinesSubnet: <uuidOfPreprovisionedSubnet>

- Run the openshift installer:

        ./openshift-install create cluster --log-level debug
        
- wait for the installer to complete. 

    If the installer times out waiting for the bare-metal workers to complete booting,
    restart the installation using the appropriate *wait-for* command:
    
        ./openshit-install watier-for install-complete --log-level debug

## Deploy cluster with vm workers only and add BM workers after 
Initial cluster deployment is done with only 
vm workers attached to the installer-provisioned network. The bare-metal
workers are added after initial installation. Bare-metal workers are
attached to a pre-existing network. Traffic between masters and workers is routed between subnets.

- Requirements: 
    - Cloud provider is configure to route traffic between vm subnets and bare-metal subnets.   
    
- Create install-config.yaml:
    - Set compute.[worker].platform.openstack.type to the vm flavor which will be used by worker nodes.
    - Set control-plane.platform.openstack.type to the vm flavor which will be used by the control-plane masters.
    
        For example:
 
                controlPlane:
                   platform:
                     openstack:
                       type: <vmComputeFlavorForMasters>
                 
                 ... other settings
                 
                 compute:
                 - architecture: amd64
                   hyperthreading: Enabled
                   name: worker
                   platform:
                     openstack:
                       type: <vmComputeFlavorForWorkers>
                   replicas: 3
         
                 ... other settings
                 
                 platform:
                   openstack:
                     computeFlavor: <vmComputeFlavorForMasters>

- Run the openshift installer:

        ./openshift-install create cluster --log-level debug
        
- wait for the installer to complete. 

- Once the installation is complete, enable routing between the subnet created by the installer and the preexisting subnet
where the bare-metal machines will be attached. 

- Deploy the bare-metal workers using machinesets and the machine-api-operator
    - Create bare-metal worker machineset yaml file: baremetalMachineSet.yaml  
        
    - Create bare-metal worker machineset resource
        
          oc create -v baremetalMachineSet.yaml



[1]: <https://docs.openstack.org/nova/latest/user/flavors.html#:~:text=Flavors%C2%B6,server%20that%20can%20be%20launched> "In OpenStack, flavors define the compute, memory, and storage capacity of nova computing instances"
[2]: <https://docs.openstack.org/ironic/latest/>
[3]: <https://docs.openstack.org/api-ref/compute/>
[4]: <https://docs.openstack.org/glance/latest/>
[5]: <https://github.com/openshift/installer/blob/master/docs/dev/openstack/customization.md>
[6]: <https://github.com/openshift/installer/blob/master/docs/user/openstack/customization.md>