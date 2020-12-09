# Required Privileges & Permissions
In order to install an OpenShift cluster to a vCenter, the user provided to the installer needs privileges to read and create the necessary resources. The easiest way to achieve this level of permission and ensure success is to install with a user who has global administrative privileges.

If the provided user has global admin privileges, no further action for permissions is required. Otherwise, the rest of this document can be used as a resource to create a user with more fine-grained privileges.

## Create new roles with the appropriate privileges

In the tables below describe the absolute minimal set of privileges to install and run OpenShift including Machine management and the vSphere Storage provider.

### Precreated virtual machine folder

If there is a pre-existing virtual machine folder that OpenShift RHCOS guests will be created within the privilege set can be used below.

Role Name | vSphere object | Privilege Set
--- | --- | ---
openshift-vcenter-level | vSphere vCenter | Cns.Searchable<br/>InventoryService.Tagging.AttachTag<br/>InventoryService.Tagging.CreateCategory<br/>InventoryService.Tagging.CreateTag<br/>InventoryService.Tagging.DeleteCategory<br/>InventoryService.Tagging.DeleteTag<br/>InventoryService.Tagging.EditCategory<br/>InventoryService.Tagging.EditTag<br/>Sessions.ValidateSession<br/>StorageProfile.View
openshift-cluster-level | vSphere vCenter Cluster | Host.Config.Storage<br/>Resource.AssignVMToPool<br/>VApp.AssignResourcePool<br/>VApp.Import<br/>VirtualMachine.Config.AddNewDisk
openshift-datastore-level| vSphere Datastore | Datastore.AllocateSpace<br/>Datastore.Browse<br/>Datastore.FileManagement
openshift-portgroup-level | vSphere Port Group | Network.Assign
openshift-folder-level| Virtual Machine Folder | Resource.AssignVMToPool<br/>VApp.Import<br/>VirtualMachine.Config.AddExistingDisk<br/>VirtualMachine.Config.AddNewDisk<br/>VirtualMachine.Config.AddRemoveDevice<br/>VirtualMachine.Config.AdvancedConfig<br/>VirtualMachine.Config.Annotation<br/>VirtualMachine.Config.CPUCount<br/>VirtualMachine.Config.DiskExtend<br/>VirtualMachine.Config.DiskLease<br/>VirtualMachine.Config.EditDevice<br/>VirtualMachine.Config.Memory<br/>VirtualMachine.Config.RemoveDisk<br/>VirtualMachine.Config.Rename<br/>VirtualMachine.Config.ResetGuestInfo<br/>VirtualMachine.Config.Resource<br/>VirtualMachine.Config.Settings<br/>VirtualMachine.Config.UpgradeVirtualHardware<br/>VirtualMachine.Interact.GuestControl<br/>VirtualMachine.Interact.PowerOff<br/>VirtualMachine.Interact.PowerOn<br/>VirtualMachine.Interact.Reset<br/>VirtualMachine.Inventory.Create<br/>VirtualMachine.Inventory.CreateFromExisting<br/>VirtualMachine.Inventory.Delete<br/>VirtualMachine.Provisioning.Clone


### Installer created virtual machine folder

Including the role-set above one additional role needs to be created if the installer is to create a vSphere virtual machine folder.
Since the datacenter's top-level virtual machine folder is hidden the only way to support installation that creates a vm folder for the OpenShift cluster is to create a new datacenter role and propagate. Once installation is complete the `openshift-folder-level` role could be applied to the folder that the installer created.

Role Name | vSphere object | Privilege Set
--- | --- | ---
openshift-datacenter-level| vSphere vCenter Datacenter | Resource.AssignVMToPool<br/>VApp.Import<br/>VirtualMachine.Config.AddExistingDisk<br/>VirtualMachine.Config.AddNewDisk<br/>VirtualMachine.Config.AddRemoveDevice<br/>VirtualMachine.Config.AdvancedConfig<br/>VirtualMachine.Config.Annotation<br/>VirtualMachine.Config.CPUCount<br/>VirtualMachine.Config.DiskExtend<br/>VirtualMachine.Config.DiskLease<br/>VirtualMachine.Config.EditDevice<br/>VirtualMachine.Config.Memory<br/>VirtualMachine.Config.RemoveDisk<br/>VirtualMachine.Config.Rename<br/>VirtualMachine.Config.ResetGuestInfo<br/>VirtualMachine.Config.Resource<br/>VirtualMachine.Config.Settings<br/>VirtualMachine.Config.UpgradeVirtualHardware<br/>VirtualMachine.Interact.GuestControl<br/>VirtualMachine.Interact.PowerOff<br/>VirtualMachine.Interact.PowerOn<br/>VirtualMachine.Interact.Reset<br/>VirtualMachine.Inventory.Create<br/>VirtualMachine.Inventory.CreateFromExisting<br/>VirtualMachine.Inventory.Delete<br/>VirtualMachine.Provisioning.Clone<br/>Folder.Create<br/>Folder.Delete

## Permission assignments

The easiest way to ensure proper permissions is to grant Global Permissions to the user with the privileges above. Otherwise, it is necessary to ensure that the user with the listed privileges has permissions granted on all necessary entities in the vCenter.

For more information, consult [vSphere Permissions and User Management Tasks][vsphere-perms]

### Precreated virtual machine folder

Role Name | Propagate | Entity
--- | --- | ---
openshift-vcenter-level | False | vSphere vCenter
ReadOnly | False | vSphere vCenter Datacenter
openshift-cluster-level | True | vSphere vCenter Cluster
openshift-datastore-level | False | vSphere vCenter Datastore
ReadOnly | False | vSphere Switch
openshift-portgroup-level | False | vSphere Port Group
openshift-folder-level | True | vSphere vCenter Virtual Machine folder


### Installer created virtual machine folder
Role Name | Propagate | Entity
--- | --- | ---
openshift-vcenter-level | False | vSphere vCenter
openshift-datacenter-level | True | vSphere vCenter Datacenter
openshift-cluster-level | True | vSphere vCenter Cluster
openshift-datastore-level | False | vSphere vCenter Datastore
ReadOnly | False | vSphere Switch
openshift-portgroup-level | False | vSphere Port Group


## Walkthrough: Creating and Assigning Global Roles
The following is a visual walkthrough of creating and assigning global roles in the vSphere 6 web client. Roles can be similarly created for specific clusters. For more information, refer to the [vSphere docs][vsphere-docs].

### Creating a new role
Roles can be created and edited in __Administration > Access Control > Roles__.

When creating a new role, first assign permissions (using the list above for guidance):
![Select privileges](images/select-privileges.png)

Once you save your role, the new privileges will be visible:
![View privileges](images/view-privileges.png)

### Assigning a role
Roles can be assigned in __Administration > Access Control > Global Permissions__.
The newly created role can be assigned to a group or directly to a user.

To assign the newly created role, click the `+` for Add Permission:
![Assign role](images/assign-role.png)

[vsphere-docs]: https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.security.doc/GUID-5372F580-5C23-4E9C-8A4E-EF1B4DD9033E.html
[vsphere-perms]: https://docs.vmware.com/en/VMware-vSphere/6.7/com.vmware.vsphere.security.doc/GUID-5372F580-5C23-4E9C-8A4E-EF1B4DD9033E.html
