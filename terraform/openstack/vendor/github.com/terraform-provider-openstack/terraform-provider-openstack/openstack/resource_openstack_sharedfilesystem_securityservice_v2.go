package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
)

func resourceSharedFilesystemSecurityServiceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemSecurityServiceV2Create,
		ReadContext:   resourceSharedFilesystemSecurityServiceV2Read,
		UpdateContext: resourceSharedFilesystemSecurityServiceV2Update,
		DeleteContext: resourceSharedFilesystemSecurityServiceV2Delete,

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
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active_directory", "kerberos", "ldap",
				}, true),
			},

			"dns_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"server": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSharedFilesystemSecurityServiceV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	createOpts := securityservices.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        securityservices.SecurityServiceType(d.Get("type").(string)),
		DNSIP:       d.Get("dns_ip").(string),
		User:        d.Get("user").(string),
		Domain:      d.Get("domain").(string),
		Server:      d.Get("server").(string),
	}

	if v, ok := d.GetOkExists("ou"); ok {
		createOpts.OU = v.(string)

		sfsClient.Microversion = sharedFilesystemV2SecurityServiceOUMicroversion
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_securityservice_v2 create options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)
	securityservice, err := securityservices.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_sharedfilesystem_securityservice_v2: %s", err)
	}

	d.SetId(securityservice.ID)

	return resourceSharedFilesystemSecurityServiceV2Read(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	// Select microversion to use.
	sfsClient.Microversion = sharedFilesystemV2MinMicroversion
	if _, ok := d.GetOkExists("ou"); ok {
		sfsClient.Microversion = sharedFilesystemV2SecurityServiceOUMicroversion
	}

	securityservice, err := securityservices.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_sharedfilesystem_securityservice_v2"))
	}

	// Workaround for resource import.
	if securityservice.OU == "" {
		sfsClient.Microversion = sharedFilesystemV2SecurityServiceOUMicroversion
		securityserviceOU, err := securityservices.Get(sfsClient, d.Id()).Extract()
		if err == nil {
			d.Set("ou", securityserviceOU.OU)
		}
	}

	nopassword := securityservice
	nopassword.Password = ""
	log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_securityservice_v2 %s: %#v", d.Id(), nopassword)

	d.Set("name", securityservice.Name)
	d.Set("description", securityservice.Description)
	d.Set("type", securityservice.Type)
	d.Set("domain", securityservice.Domain)
	d.Set("dns_ip", securityservice.DNSIP)
	d.Set("user", securityservice.User)
	d.Set("server", securityservice.Server)

	// Computed.
	d.Set("project_id", securityservice.ProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceSharedFilesystemSecurityServiceV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemV2MinMicroversion

	var updateOpts securityservices.UpdateOpts

	// Name should always be sent, otherwise it is vanished by manila backend.
	name := d.Get("name").(string)
	updateOpts.Name = &name

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("type") {
		updateOpts.Type = d.Get("type").(string)
	}

	if d.HasChange("dns_ip") {
		dnsIP := d.Get("dns_ip").(string)
		updateOpts.DNSIP = &dnsIP
	}

	if d.HasChange("ou") {
		ou := d.Get("ou").(string)
		updateOpts.OU = &ou

		sfsClient.Microversion = sharedFilesystemV2SecurityServiceOUMicroversion
	}

	if d.HasChange("user") {
		user := d.Get("user").(string)
		updateOpts.User = &user
	}

	if d.HasChange("domain") {
		domain := d.Get("domain").(string)
		updateOpts.Domain = &domain
	}

	if d.HasChange("server") {
		server := d.Get("server").(string)
		updateOpts.Server = &server
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_securityservice_v2 %s update options: %#v", d.Id(), updateOpts)

	if d.HasChange("password") {
		password := d.Get("password").(string)
		updateOpts.Password = &password
	}

	_, err = securityservices.Update(sfsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_sharedfilesystem_securityservice_v2 %s: %s", d.Id(), err)
	}

	return resourceSharedFilesystemSecurityServiceV2Read(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	if err := securityservices.Delete(sfsClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_securityservice_v2"))
	}

	return nil
}
