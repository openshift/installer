package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/apiversions"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func resourceSharedFilesystemShareAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareAccessV2Create,
		ReadContext:   resourceSharedFilesystemShareAccessV2Read,
		DeleteContext: resourceSharedFilesystemShareAccessV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSharedFilesystemShareAccessV2Import,
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

			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "user", "cert", "cephx",
				}, false),
			},

			"access_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rw", "ro",
				}, false),
			},

			"access_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSharedFilesystemShareAccessV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion
	accessType := d.Get("access_type").(string)
	if accessType == "cephx" {
		sfsClient.Microversion = sharedFilesystemV2SharedAccessCephXMicroversion
	}

	shareID := d.Get("share_id").(string)

	grantOpts := shares.GrantAccessOpts{
		AccessType:  accessType,
		AccessTo:    d.Get("access_to").(string),
		AccessLevel: d.Get("access_level").(string),
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_share_access_v2 create options: %#v", grantOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	var access *shares.AccessRight
	err = resource.Retry(timeout, func() *resource.RetryError {
		access, err = shares.GrantAccess(sfsClient, shareID, grantOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		detailedErr := errors.ErrorDetails{}
		e := errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating openstack_sharedfilesystem_share_access_v2: %s: %s", err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error creating openstack_sharedfilesystem_share_access_v2: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_sharedfilesystem_share_access_v2 %s to become available.", access.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"active"},
		Pending:    []string{"new", "queued_to_apply", "applying"},
		Refresh:    sharedFilesystemShareAccessV2StateRefreshFunc(sfsClient, shareID, access.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 %s to become available: %s", access.ID, err)
	}

	d.SetId(access.ID)

	return resourceSharedFilesystemShareAccessV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareAccessV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	// Set the client to the minimum supported microversion.
	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	// Now check and see if the OpenStack environment supports microversion 2.21.
	// If so, use that for the API request for access_key support.
	apiInfo, err := apiversions.Get(sfsClient, "v2").Extract()
	if err != nil {
		return diag.Errorf("Unable to query API endpoint for openstack_sharedfilesystem_share_access_v2: %s", err)
	}

	compatible, err := compatibleMicroversion("min", "2.21", apiInfo.Version)
	if err != nil {
		return diag.Errorf("Error comparing microversions for openstack_sharedfilesystem_share_access_v2 %s: %s", d.Id(), err)
	}

	if compatible {
		sfsClient.Microversion = sharedFilesystemV2SharedAccessMinMicroversion
	}

	shareID := d.Get("share_id").(string)
	access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_sharedfilesystem_share_access_v2"))
	}

	for _, v := range access {
		if v.ID == d.Id() {
			log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_share_access_v2 %s: %#v", d.Id(), v)

			d.Set("access_type", v.AccessType)
			d.Set("access_to", v.AccessTo)
			d.Set("access_level", v.AccessLevel)
			d.Set("region", GetRegion(d, config))

			// This will only be set if the Shared Filesystem environment supports
			// microversion 2.21.
			d.Set("access_key", v.AccessKey)

			return nil
		}
	}

	log.Printf("[DEBUG] Unable to find openstack_sharedfilesystem_share_access_v2 %s", d.Id())
	d.SetId("")

	return nil
}

func resourceSharedFilesystemShareAccessV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	shareID := d.Get("share_id").(string)

	revokeOpts := shares.RevokeAccessOpts{AccessID: d.Id()}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete openstack_sharedfilesystem_share_access_v2 %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = shares.RevokeAccess(sfsClient, shareID, revokeOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		e := CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_share_access_v2")
		if e == nil {
			return nil
		}
		detailedErr := errors.ErrorDetails{}
		e = errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 on %s to be removed: %s: %s", shareID, err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 on %s to be removed: %s (%d): %s", shareID, k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for openstack_sharedfilesystem_share_access_v2 %s to become denied.", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"denied"},
		Pending:    []string{"active", "new", "queued_to_deny", "denying"},
		Refresh:    sharedFilesystemShareAccessV2StateRefreshFunc(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return nil
		}
		return diag.Errorf("Error waiting for openstack_sharedfilesystem_share_access_v2 %s to become denied: %s", d.Id(), err)
	}

	return nil
}

func resourceSharedFilesystemShareAccessV2Import(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for openstack_sharedfilesystem_share_access_v2. Format must be <share id>/<ACL id>")
		return nil, err
	}

	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	shareID := parts[0]
	accessID := parts[1]

	access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Unable to get %s openstack_sharedfilesystem_share_v2: %s", shareID, err)
	}

	for _, v := range access {
		if v.ID == accessID {
			log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_share_access_v2 %s: %#v", accessID, v)

			d.SetId(accessID)
			d.Set("share_id", shareID)
			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("[DEBUG] Unable to find openstack_sharedfilesystem_share_access_v2 %s", accessID)
}
