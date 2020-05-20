package vsphere

import (
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// defaultAPITimeout is a default timeout value that is passed to functions
// requiring contexts, and other various waiters.
const defaultAPITimeout = time.Minute * 5

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USER", nil),
				Description: "The user name for vSphere API operations.",
			},

			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
				Description: "The user password for vSphere API operations.",
			},

			"vsphere_server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_SERVER", nil),
				Description: "The vSphere Server name for vSphere API operations.",
			},
			"allow_unverified_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_ALLOW_UNVERIFIED_SSL", false),
				Description: "If set, VMware vSphere client will permit unverifiable SSL certificates.",
			},
			"vcenter_server": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_VCENTER", nil),
				Deprecated:  "This field has been renamed to vsphere_server.",
			},
			"client_debug": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG", false),
				Description: "govmomi debug",
			},
			"client_debug_path_run": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG_PATH_RUN", ""),
				Description: "govmomi debug path for a single run",
			},
			"client_debug_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_CLIENT_DEBUG_PATH", ""),
				Description: "govmomi debug path for debug",
			},
			"persist_session": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PERSIST_SESSION", false),
				Description: "Persist vSphere client sessions to disk",
			},
			"vim_session_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_VIM_SESSION_PATH", filepath.Join(os.Getenv("HOME"), ".govmomi", "sessions")),
				Description: "The directory to save vSphere SOAP API sessions to",
			},
			"rest_session_path": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_REST_SESSION_PATH", filepath.Join(os.Getenv("HOME"), ".govmomi", "rest_sessions")),
				Description: "The directory to save vSphere REST API sessions to",
			},
			"vim_keep_alive": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_VIM_KEEP_ALIVE", 10),
				Description: "Keep alive interval for the VIM session in minutes",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vsphere_compute_cluster":                         resourceVSphereComputeCluster(),
			"vsphere_compute_cluster_host_group":              resourceVSphereComputeClusterHostGroup(),
			"vsphere_compute_cluster_vm_affinity_rule":        resourceVSphereComputeClusterVMAffinityRule(),
			"vsphere_compute_cluster_vm_anti_affinity_rule":   resourceVSphereComputeClusterVMAntiAffinityRule(),
			"vsphere_compute_cluster_vm_dependency_rule":      resourceVSphereComputeClusterVMDependencyRule(),
			"vsphere_compute_cluster_vm_group":                resourceVSphereComputeClusterVMGroup(),
			"vsphere_compute_cluster_vm_host_rule":            resourceVSphereComputeClusterVMHostRule(),
			"vsphere_content_library":                         resourceVSphereContentLibrary(),
			"vsphere_content_library_item":                    resourceVSphereContentLibraryItem(),
			"vsphere_custom_attribute":                        resourceVSphereCustomAttribute(),
			"vsphere_datacenter":                              resourceVSphereDatacenter(),
			"vsphere_datastore_cluster":                       resourceVSphereDatastoreCluster(),
			"vsphere_datastore_cluster_vm_anti_affinity_rule": resourceVSphereDatastoreClusterVMAntiAffinityRule(),
			"vsphere_distributed_port_group":                  resourceVSphereDistributedPortGroup(),
			"vsphere_distributed_virtual_switch":              resourceVSphereDistributedVirtualSwitch(),
			"vsphere_drs_vm_override":                         resourceVSphereDRSVMOverride(),
			"vsphere_dpm_host_override":                       resourceVSphereDPMHostOverride(),
			"vsphere_file":                                    resourceVSphereFile(),
			"vsphere_folder":                                  resourceVSphereFolder(),
			"vsphere_ha_vm_override":                          resourceVSphereHAVMOverride(),
			"vsphere_host_port_group":                         resourceVSphereHostPortGroup(),
			"vsphere_host_virtual_switch":                     resourceVSphereHostVirtualSwitch(),
			"vsphere_license":                                 resourceVSphereLicense(),
			"vsphere_resource_pool":                           resourceVSphereResourcePool(),
			"vsphere_tag":                                     resourceVSphereTag(),
			"vsphere_tag_category":                            resourceVSphereTagCategory(),
			"vsphere_virtual_disk":                            resourceVSphereVirtualDisk(),
			"vsphere_virtual_machine":                         resourceVSphereVirtualMachine(),
			"vsphere_nas_datastore":                           resourceVSphereNasDatastore(),
			"vsphere_storage_drs_vm_override":                 resourceVSphereStorageDrsVMOverride(),
			"vsphere_vapp_container":                          resourceVSphereVAppContainer(),
			"vsphere_vapp_entity":                             resourceVSphereVAppEntity(),
			"vsphere_vmfs_datastore":                          resourceVSphereVmfsDatastore(),
			"vsphere_virtual_machine_snapshot":                resourceVSphereVirtualMachineSnapshot(),
			"vsphere_host":                                    resourceVsphereHost(),
			"vsphere_vnic":                                    resourceVsphereNic(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vsphere_compute_cluster":            dataSourceVSphereComputeCluster(),
			"vsphere_content_library":            dataSourceVSphereContentLibrary(),
			"vsphere_content_library_item":       dataSourceVSphereContentLibraryItem(),
			"vsphere_custom_attribute":           dataSourceVSphereCustomAttribute(),
			"vsphere_datacenter":                 dataSourceVSphereDatacenter(),
			"vsphere_datastore":                  dataSourceVSphereDatastore(),
			"vsphere_datastore_cluster":          dataSourceVSphereDatastoreCluster(),
			"vsphere_distributed_virtual_switch": dataSourceVSphereDistributedVirtualSwitch(),
			"vsphere_folder":                     dataSourceVSphereFolder(),
			"vsphere_host":                       dataSourceVSphereHost(),
			"vsphere_network":                    dataSourceVSphereNetwork(),
			"vsphere_resource_pool":              dataSourceVSphereResourcePool(),
			"vsphere_storage_policy":             dataSourceVSphereStoragePolicy(),
			"vsphere_tag":                        dataSourceVSphereTag(),
			"vsphere_tag_category":               dataSourceVSphereTagCategory(),
			"vsphere_vapp_container":             dataSourceVSphereVAppContainer(),
			"vsphere_virtual_machine":            dataSourceVSphereVirtualMachine(),
			"vsphere_vmfs_disks":                 dataSourceVSphereVmfsDisks(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	c, err := NewConfig(d)
	if err != nil {
		return nil, err
	}
	return c.Client()
}
