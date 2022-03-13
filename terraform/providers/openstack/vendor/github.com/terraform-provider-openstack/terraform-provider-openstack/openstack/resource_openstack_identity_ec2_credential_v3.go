package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/ec2credentials"
)

func resourceIdentityEc2CredentialV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityEc2CredentialV3Create,
		ReadContext:   resourceIdentityEc2CredentialV3Read,
		DeleteContext: resourceIdentityEc2CredentialV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"access": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
				ForceNew:  true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"trust_id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityEc2CredentialV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	tokenInfo, err := getTokenInfo(identityClient)
	if err != nil {
		return diag.Errorf("Error getting token info: %s", err)
	}

	var tenantID string
	// Set tenant to user value if provided, get one from token otherwise
	if definedProject, ok := d.GetOk("project_id"); ok {
		tenantID = definedProject.(string)
	} else {
		tenantID = tokenInfo.projectID
	}

	createOpts := ec2credentials.CreateOpts{
		TenantID: tenantID,
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = tokenInfo.userID
	}

	log.Printf("[DEBUG] openstack_identity_ec2_credential_v3 create options: %#v", createOpts)

	ec2Credential, err := ec2credentials.Create(identityClient, userID, createOpts).Extract()

	if err != nil {
		if v, ok := err.(gophercloud.ErrDefault404); ok {
			return diag.Errorf("Error creating openstack_identity_ec2_credential_v3: %s", v.ErrUnexpectedResponseCode.Body)
		}
		return diag.Errorf("Error creating openstack_identity_ec2_credential_v3: %s", err)
	}

	d.SetId(ec2Credential.Access)

	d.Set("access", ec2Credential.Access)
	d.Set("secret", ec2Credential.Secret)
	d.Set("user_id", ec2Credential.UserID)
	d.Set("project_id", ec2Credential.TenantID)
	d.Set("trust_id", ec2Credential.TrustID)

	return resourceIdentityEc2CredentialV3Read(ctx, d, meta)
}

func resourceIdentityEc2CredentialV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	tokenInfo, err := getTokenInfo(identityClient)
	if err != nil {
		return diag.Errorf("Error getting token info: %s", err)
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = tokenInfo.userID
	}

	ec2Credential, err := ec2credentials.Get(identityClient, userID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_ec2_credential_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_ec2_credential_v3 %s: %#v", d.Id(), ec2Credential)

	d.Set("secret", ec2Credential.Secret)
	d.Set("user_id", ec2Credential.UserID)
	d.Set("project_id", ec2Credential.TenantID)
	d.Set("access", ec2Credential.Access)
	d.Set("trust_id", ec2Credential.TrustID)

	return nil
}

func resourceIdentityEc2CredentialV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	tokenInfo, err := getTokenInfo(identityClient)
	if err != nil {
		return diag.Errorf("Error getting token info: %s", err)
	}

	var userID string
	// Set userid to defined tf value, use one from token otherwise
	if definedUser, ok := d.GetOk("user_id"); ok {
		userID = definedUser.(string)
	} else {
		userID = tokenInfo.userID
	}

	err = ec2credentials.Delete(identityClient, userID, d.Id()).ExtractErr()
	if err != nil {
		err = CheckDeleted(d, err, "Error deleting openstack_identity_ec2_credential_v3")
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
