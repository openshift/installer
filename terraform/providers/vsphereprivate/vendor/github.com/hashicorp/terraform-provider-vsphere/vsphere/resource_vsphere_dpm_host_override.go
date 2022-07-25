package vsphere

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/clustercomputeresource"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
)

const resourceVSphereDPMHostOverrideName = "vsphere_dpm_host_override"

func resourceVSphereDPMHostOverride() *schema.Resource {
	return &schema.Resource{
		Create: resourceVSphereDPMHostOverrideCreate,
		Read:   resourceVSphereDPMHostOverrideRead,
		Update: resourceVSphereDPMHostOverrideUpdate,
		Delete: resourceVSphereDPMHostOverrideDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVSphereDPMHostOverrideImport,
		},

		Schema: map[string]*schema.Schema{
			"compute_cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the cluster.",
			},
			"host_system_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The managed object ID of the host.",
			},
			"dpm_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable DPM for this host.",
			},
			"dpm_automation_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(types.DpmBehaviorManual),
				Description:  "The automation level for power operations on this host. Can be one of manual or automated.",
				ValidateFunc: validation.StringInSlice(dpmBehaviorAllowedValues, false),
			},
		},
	}
}

func resourceVSphereDPMHostOverrideCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning create", resourceVSphereDPMHostOverrideIDString(d))

	cluster, host, err := resourceVSphereDPMHostOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDpmHostConfigInfo(d, host)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DpmHostConfigSpec: []types.ClusterDpmHostConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationAdd,
				},
				Info: info,
			},
		},
	}

	if err = clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	id, err := resourceVSphereDPMHostOverrideFlattenID(cluster, host)
	if err != nil {
		return fmt.Errorf("cannot compute ID of created resource: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] %s: Create finished successfully", resourceVSphereDPMHostOverrideIDString(d))
	return resourceVSphereDPMHostOverrideRead(d, meta)
}

func resourceVSphereDPMHostOverrideRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning read", resourceVSphereDPMHostOverrideIDString(d))

	cluster, host, err := resourceVSphereDPMHostOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := resourceVSphereDPMHostOverrideFindEntry(cluster, host)
	if err != nil {
		return err
	}

	if info == nil {
		// The configuration is missing, blank out the ID so it can be re-created.
		d.SetId("")
		return nil
	}

	// Save the compute_cluster_id and host_system_id here. These are
	// ForceNew, but we set these for completeness on import so that if the wrong
	// cluster/VM combo was used, it will be noted.
	if err = d.Set("compute_cluster_id", cluster.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"compute_cluster_id\": %s", err)
	}

	if err = d.Set("host_system_id", host.Reference().Value); err != nil {
		return fmt.Errorf("error setting attribute \"host_system_id\": %s", err)
	}

	if err = flattenClusterDpmHostConfigInfo(d, info); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Read completed successfully", resourceVSphereDPMHostOverrideIDString(d))
	return nil
}

func resourceVSphereDPMHostOverrideUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning update", resourceVSphereDPMHostOverrideIDString(d))

	cluster, host, err := resourceVSphereDPMHostOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	info, err := expandClusterDpmHostConfigInfo(d, host)
	if err != nil {
		return err
	}
	spec := &types.ClusterConfigSpecEx{
		DpmHostConfigSpec: []types.ClusterDpmHostConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationEdit,
				},
				Info: info,
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Update finished successfully", resourceVSphereDPMHostOverrideIDString(d))
	return resourceVSphereDPMHostOverrideRead(d, meta)
}

func resourceVSphereDPMHostOverrideDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] %s: Beginning delete", resourceVSphereDPMHostOverrideIDString(d))

	cluster, host, err := resourceVSphereDPMHostOverrideObjects(d, meta)
	if err != nil {
		return err
	}

	spec := &types.ClusterConfigSpecEx{
		DpmHostConfigSpec: []types.ClusterDpmHostConfigSpec{
			{
				ArrayUpdateSpec: types.ArrayUpdateSpec{
					Operation: types.ArrayUpdateOperationRemove,
					RemoveKey: host.Reference(),
				},
			},
		},
	}

	if err := clustercomputeresource.Reconfigure(cluster, spec); err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Deleted successfully", resourceVSphereDPMHostOverrideIDString(d))
	return nil
}

func resourceVSphereDPMHostOverrideImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var data map[string]string
	if err := json.Unmarshal([]byte(d.Id()), &data); err != nil {
		return nil, err
	}
	clusterPath, ok := data["compute_cluster_path"]
	if !ok {
		return nil, errors.New("missing compute_cluster_path in input data")
	}
	hostPath, ok := data["host_path"]
	if !ok {
		return nil, errors.New("missing host_path in input data")
	}

	client, err := resourceVSphereDPMHostOverrideClient(meta)
	if err != nil {
		return nil, err
	}

	cluster, err := clustercomputeresource.FromPath(client, clusterPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate cluster %q: %s", clusterPath, err)
	}

	host, err := hostsystem.SystemOrDefault(client, hostPath, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot locate host %q: %s", hostPath, err)
	}

	id, err := resourceVSphereDPMHostOverrideFlattenID(cluster, host)
	if err != nil {
		return nil, fmt.Errorf("cannot compute ID of imported resource: %s", err)
	}
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}

// expandClusterDrsVMConfigInfo reads certain ResourceData keys and returns a
// ClusterDpmHostConfigInfo.
func expandClusterDpmHostConfigInfo(d *schema.ResourceData, host *object.HostSystem) (*types.ClusterDpmHostConfigInfo, error) {
	obj := &types.ClusterDpmHostConfigInfo{
		Behavior: types.DpmBehavior(d.Get("dpm_automation_level").(string)),
		Enabled:  structure.GetBool(d, "dpm_enabled"),
		Key:      host.Reference(),
	}

	return obj, nil
}

// flattenClusterDpmHostConfigInfo saves a ClusterDpmHostConfigInfo into the
// supplied ResourceData.
func flattenClusterDpmHostConfigInfo(d *schema.ResourceData, obj *types.ClusterDpmHostConfigInfo) error {
	return structure.SetBatch(d, map[string]interface{}{
		"dpm_automation_level": obj.Behavior,
		"dpm_enabled":          obj.Enabled,
	})
}

// resourceVSphereDPMHostOverrideIDString prints a friendly string for the
// vsphere_dpm_host_override resource.
func resourceVSphereDPMHostOverrideIDString(d structure.ResourceIDStringer) string {
	return structure.ResourceIDString(d, resourceVSphereDPMHostOverrideName)
}

// resourceVSphereDPMHostOverrideFlattenID makes an ID for the
// vsphere_dpm_host_override resource.
func resourceVSphereDPMHostOverrideFlattenID(cluster *object.ClusterComputeResource, host *object.HostSystem) (string, error) {
	return strings.Join([]string{cluster.Reference().Value, host.Reference().Value}, ":"), nil
}

// resourceVSphereDPMHostOverrideParseID parses an ID for the
// vsphere_dpm_host_override and outputs its parts.
func resourceVSphereDPMHostOverrideParseID(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 3)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("bad ID %q", id)
	}
	return parts[0], parts[1], nil
}

// resourceVSphereDPMHostOverrideFindEntry attempts to locate an existing DRS VM
// config in a cluster's configuration. It's used by the resource's read
// functionality and tests. nil is returned if the entry cannot be found.
func resourceVSphereDPMHostOverrideFindEntry(
	cluster *object.ClusterComputeResource,
	host *object.HostSystem,
) (*types.ClusterDpmHostConfigInfo, error) {
	props, err := clustercomputeresource.Properties(cluster)
	if err != nil {
		return nil, fmt.Errorf("error fetching cluster properties: %s", err)
	}

	for _, info := range props.ConfigurationEx.(*types.ClusterConfigInfoEx).DpmHostConfig {
		if info.Key == host.Reference() {
			log.Printf("[DEBUG] Found DPM config info for host %q in cluster %q", host.Name(), cluster.Name())
			return &info, nil
		}
	}

	log.Printf("[DEBUG] No DPM config info found for host %q in cluster %q", host.Name(), cluster.Name())
	return nil, nil
}

// resourceVSphereDPMHostOverrideObjects handles the fetching of the cluster
// and host depending on what attributes are available:
// * If the resource ID is available, the data is derived from the ID.
// * If not, it's derived from the compute_cluster_id and host_system_id
// attributes.
func resourceVSphereDPMHostOverrideObjects(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.HostSystem, error) {
	if d.Id() != "" {
		return resourceVSphereDPMHostOverrideObjectsFromID(d, meta)
	}
	return resourceVSphereDPMHostOverrideObjectsFromAttributes(d, meta)
}

func resourceVSphereDPMHostOverrideObjectsFromAttributes(
	d *schema.ResourceData,
	meta interface{},
) (*object.ClusterComputeResource, *object.HostSystem, error) {
	return resourceVSphereDPMHostOverrideFetchObjects(
		meta,
		d.Get("compute_cluster_id").(string),
		d.Get("host_system_id").(string),
	)
}

func resourceVSphereDPMHostOverrideObjectsFromID(
	d structure.ResourceIDStringer,
	meta interface{},
) (*object.ClusterComputeResource, *object.HostSystem, error) {
	// Note that this function uses structure.ResourceIDStringer to satisfy
	// interfacer. Adding exceptions in the comments does not seem to work.
	// Change this back to ResourceData if it's needed in the future.
	clusterID, hostID, err := resourceVSphereDPMHostOverrideParseID(d.Id())
	if err != nil {
		return nil, nil, err
	}

	return resourceVSphereDPMHostOverrideFetchObjects(meta, clusterID, hostID)
}

func resourceVSphereDPMHostOverrideFetchObjects(
	meta interface{},
	clusterID string,
	hostID string,
) (*object.ClusterComputeResource, *object.HostSystem, error) {
	client, err := resourceVSphereDPMHostOverrideClient(meta)
	if err != nil {
		return nil, nil, err
	}

	cluster, err := clustercomputeresource.FromID(client, clusterID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate cluster: %s", err)
	}

	host, err := hostsystem.FromID(client, hostID)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot locate virtual machine: %s", err)
	}

	return cluster, host, nil
}

func resourceVSphereDPMHostOverrideClient(meta interface{}) (*govmomi.Client, error) {
	client := meta.(*Client).vimClient
	if err := viapi.ValidateVirtualCenter(client); err != nil {
		return nil, err
	}
	return client, nil
}
