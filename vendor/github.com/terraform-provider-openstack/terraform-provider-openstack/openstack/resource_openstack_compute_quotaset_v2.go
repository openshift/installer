package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceComputeQuotasetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeQuotasetV2Create,
		Read:   resourceComputeQuotasetV2Read,
		Update: resourceComputeQuotasetV2Update,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceComputeQuotasetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	projectID := d.Get("project_id").(string)
	fixedIPs := d.Get("fixed_ips").(int)
	floatingIPs := d.Get("floating_ips").(int)
	injectedFileContentBytes := d.Get("injected_file_content_bytes").(int)
	injectedFilePathBytes := d.Get("injected_file_path_bytes").(int)
	injectedFiles := d.Get("injected_files").(int)
	keyPairs := d.Get("key_pairs").(int)
	metadataItems := d.Get("metadata_items").(int)
	ram := d.Get("ram").(int)
	securityGroupRules := d.Get("security_group_rules").(int)
	securityGroups := d.Get("security_groups").(int)
	cores := d.Get("cores").(int)
	instances := d.Get("instances").(int)
	serverGroups := d.Get("server_groups").(int)
	serverGroupMembers := d.Get("server_group_members").(int)

	updateOpts := quotasets.UpdateOpts{
		FixedIPs:                 &fixedIPs,
		FloatingIPs:              &floatingIPs,
		InjectedFileContentBytes: &injectedFileContentBytes,
		InjectedFilePathBytes:    &injectedFilePathBytes,
		InjectedFiles:            &injectedFiles,
		KeyPairs:                 &keyPairs,
		MetadataItems:            &metadataItems,
		RAM:                      &ram,
		SecurityGroupRules:       &securityGroupRules,
		SecurityGroups:           &securityGroups,
		Cores:                    &cores,
		Instances:                &instances,
		ServerGroups:             &serverGroups,
		ServerGroupMembers:       &serverGroupMembers,
	}

	q, err := quotasets.Update(computeClient, projectID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_compute_quotaset_v2: %s", err)
	}

	d.SetId(projectID)

	log.Printf("[DEBUG] Created openstack_compute_quotaset_v2 %#v", q)

	return resourceComputeQuotasetV2Read(d, meta)
}

func resourceComputeQuotasetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	q, err := quotasets.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_compute_quotaset_v2")
	}

	log.Printf("[DEBUG] Retrieved openstack_compute_quotaset_v2 %s: %#v", d.Id(), q)

	d.Set("project_id", d.Id())
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

func resourceComputeQuotasetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
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
		_, err := quotasets.Update(computeClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_compute_quotaset_v2: %s", err)
		}
	}

	return resourceComputeQuotasetV2Read(d, meta)
}
