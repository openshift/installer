package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
)

func resourceComputeQuotasetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeQuotasetV2Create,
		ReadContext:   resourceComputeQuotasetV2Read,
		UpdateContext: resourceComputeQuotasetV2Update,
		Delete:        schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"fixed_ips": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"floating_ips": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"injected_file_content_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"injected_file_path_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"injected_files": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"key_pairs": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"metadata_items": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_group_rules": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"security_groups": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"cores": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"instances": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"server_groups": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"server_group_members": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceComputeQuotasetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	projectID := d.Get("project_id").(string)

	updateOpts := quotasets.UpdateOpts{}

	if v, ok := d.GetOkExists("fixed_ips"); ok {
		value := v.(int)
		updateOpts.FixedIPs = &value
	}
	if v, ok := d.GetOkExists("floating_ips"); ok {
		value := v.(int)
		updateOpts.FloatingIPs = &value
	}
	if v, ok := d.GetOkExists("injected_file_content_bytes"); ok {
		value := v.(int)
		updateOpts.InjectedFileContentBytes = &value
	}
	if v, ok := d.GetOkExists("injected_file_path_bytes"); ok {
		value := v.(int)
		updateOpts.InjectedFilePathBytes = &value
	}
	if v, ok := d.GetOkExists("injected_files"); ok {
		value := v.(int)
		updateOpts.InjectedFiles = &value
	}
	if v, ok := d.GetOkExists("key_pairs"); ok {
		value := v.(int)
		updateOpts.KeyPairs = &value
	}
	if v, ok := d.GetOkExists("metadata_items"); ok {
		value := v.(int)
		updateOpts.MetadataItems = &value
	}
	if v, ok := d.GetOkExists("ram"); ok {
		value := v.(int)
		updateOpts.RAM = &value
	}
	if v, ok := d.GetOkExists("security_group_rules"); ok {
		value := v.(int)
		updateOpts.SecurityGroupRules = &value
	}
	if v, ok := d.GetOkExists("security_groups"); ok {
		value := v.(int)
		updateOpts.SecurityGroups = &value
	}
	if v, ok := d.GetOkExists("cores"); ok {
		value := v.(int)
		updateOpts.Cores = &value
	}
	if v, ok := d.GetOkExists("instances"); ok {
		value := v.(int)
		updateOpts.Instances = &value
	}
	if v, ok := d.GetOkExists("server_groups"); ok {
		value := v.(int)
		updateOpts.ServerGroups = &value
	}
	if v, ok := d.GetOkExists("server_group_members"); ok {
		value := v.(int)
		updateOpts.ServerGroupMembers = &value
	}

	q, err := quotasets.Update(computeClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_compute_quotaset_v2: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_compute_quotaset_v2 %#v", q)

	return resourceComputeQuotasetV2Read(ctx, d, meta)
}

func resourceComputeQuotasetV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	// Depending on the provider version the resource was created, the resource id
	// can be either <project_id> or <project_id>/<region>. This parses the project_id
	// in both cases
	projectID := strings.Split(d.Id(), "/")[0]

	q, err := quotasets.Get(computeClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_compute_quotaset_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_quotaset_v2 %s: %#v", d.Id(), q)

	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("fixed_ips", q.FixedIPs)
	d.Set("floating_ips", q.FloatingIPs)
	d.Set("injected_file_content_bytes", q.InjectedFileContentBytes)
	d.Set("injected_file_path_bytes", q.InjectedFilePathBytes)
	d.Set("injected_files", q.InjectedFiles)
	d.Set("key_pairs", q.KeyPairs)
	d.Set("metadata_items", q.MetadataItems)
	d.Set("ram", q.RAM)
	d.Set("security_group_rules", q.SecurityGroupRules)
	d.Set("security_groups", q.SecurityGroups)
	d.Set("cores", q.Cores)
	d.Set("instances", q.Instances)
	d.Set("server_groups", q.ServerGroups)
	d.Set("server_group_members", q.ServerGroupMembers)

	return nil
}

func resourceComputeQuotasetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts quotasets.UpdateOpts
	)

	if d.HasChange("fixed_ips") {
		hasChange = true
		fixedIPs := d.Get("fixed_ips").(int)
		updateOpts.FixedIPs = &fixedIPs
	}

	if d.HasChange("floating_ips") {
		hasChange = true
		floatingIPs := d.Get("floating_ips").(int)
		updateOpts.FloatingIPs = &floatingIPs
	}

	if d.HasChange("injected_file_content_bytes") {
		hasChange = true
		injectedFileContentBytes := d.Get("injected_file_content_bytes").(int)
		updateOpts.InjectedFileContentBytes = &injectedFileContentBytes
	}

	if d.HasChange("injected_file_path_bytes") {
		hasChange = true
		injectedFilePathBytes := d.Get("injected_file_path_bytes").(int)
		updateOpts.InjectedFilePathBytes = &injectedFilePathBytes
	}

	if d.HasChange("injected_files") {
		hasChange = true
		injectedFiles := d.Get("injected_files").(int)
		updateOpts.InjectedFiles = &injectedFiles
	}

	if d.HasChange("key_pairs") {
		hasChange = true
		keyPairs := d.Get("key_pairs").(int)
		updateOpts.KeyPairs = &keyPairs
	}

	if d.HasChange("metadata_items") {
		hasChange = true
		metadataItems := d.Get("metadata_items").(int)
		updateOpts.MetadataItems = &metadataItems
	}

	if d.HasChange("ram") {
		hasChange = true
		ram := d.Get("ram").(int)
		updateOpts.RAM = &ram
	}

	if d.HasChange("security_group_rules") {
		hasChange = true
		securityGroupRules := d.Get("security_group_rules").(int)
		updateOpts.SecurityGroupRules = &securityGroupRules
	}

	if d.HasChange("security_groups") {
		hasChange = true
		securityGroups := d.Get("security_groups").(int)
		updateOpts.SecurityGroups = &securityGroups
	}

	if d.HasChange("cores") {
		hasChange = true
		cores := d.Get("cores").(int)
		updateOpts.Cores = &cores
	}

	if d.HasChange("instances") {
		hasChange = true
		instances := d.Get("instances").(int)
		updateOpts.Instances = &instances
	}

	if d.HasChange("server_groups") {
		hasChange = true
		serverGroups := d.Get("server_groups").(int)
		updateOpts.ServerGroups = &serverGroups
	}

	if d.HasChange("server_group_members") {
		hasChange = true
		serverGroupMembers := d.Get("server_group_members").(int)
		updateOpts.ServerGroupMembers = &serverGroupMembers
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_compute_quotaset_v2 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)
		_, err := quotasets.Update(computeClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_compute_quotaset_v2: %s", err)
		}
	}

	return resourceComputeQuotasetV2Read(ctx, d, meta)
}
