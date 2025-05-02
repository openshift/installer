package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMAppIDCloudDirectoryUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage AppID Cloud Directory user",
		CreateContext: resourceIBMAppIDCloudDirectoryUserCreate,
		ReadContext:   resourceIBMAppIDCloudDirectoryUserRead,
		DeleteContext: resourceIBMAppIDCloudDirectoryUserDelete,
		UpdateContext: resourceIBMAppIDCloudDirectoryUserUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"active": {
				Description: "Determines if the user account is active or not",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"create_profile": {
				Description: "A boolean indication if a profile should be created for the Cloud Directory user",
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Default:     true,
			},
			"locked_until": {
				Description: "Integer (epoch time in milliseconds), determines till when the user account will be locked",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"user_id": {
				Description: "Cloud Directory user ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"subject": {
				Description: "The user's identifier ('subject' in identity token)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_name": {
				Description: "Cloud Directory user display name",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"user_name": {
				Description: "Optional username",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"password": {
				Description: "User password",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"status": {
				Description:  "Accepted values `PENDING` or `CONFIRMED`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PENDING",
				ValidateFunc: validation.StringInSlice([]string{"PENDING", "CONFIRMED"}, false),
			},
			"email": {
				Description: "A set of user emails",
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Description: "User email",
							Type:        schema.TypeString,
							Required:    true,
						},
						"primary": {
							Description: "`true` if this is primary email",
							Type:        schema.TypeBool,
							Optional:    true,
						},
					},
				},
			},
			"meta": {
				Description: "Cloud Directory user metadata",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created": {
							Description: "User creation date",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"last_modified": {
							Description: "Last user modification date",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func resourceIBMAppIDCloudDirectoryUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/userID", id)
	}

	tenantID := idParts[0]
	userID := idParts[1]

	user, resp, err := appIDClient.GetCloudDirectoryUserWithContext(ctx, &appid.GetCloudDirectoryUserOptions{
		TenantID: &tenantID,
		UserID:   &userID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Cloud Directory user from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID Cloud Directory user: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("user_id", userID)

	if user.DisplayName != nil {
		d.Set("display_name", *user.DisplayName)
	}

	if user.UserName != nil {
		d.Set("user_name", *user.UserName)
	}

	if user.Status != nil {
		d.Set("status", *user.Status)
	}

	if user.Active != nil {
		d.Set("active", *user.Active)
	}

	if user.LockedUntil != nil {
		d.Set("locked_until", *user.LockedUntil)
	}

	if user.Emails != nil {
		if err := d.Set("email", flattenAppIDUserEmails(user.Emails)); err != nil {
			return diag.Errorf("Error setting AppID user emails: %s", err)
		}
	}

	if user.Meta != nil {
		if err := d.Set("meta", flattenAppIDUserMetadata(user.Meta)); err != nil {
			return diag.Errorf("Error setting AppID user metadata: %s", err)
		}
	}

	attr, resp, err := appIDClient.CloudDirectoryGetUserinfoWithContext(ctx, &appid.CloudDirectoryGetUserinfoOptions{
		TenantID: &tenantID,
		UserID:   &userID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID user attributes: %s\n%s", err, resp)
	}

	if attr.Sub != nil {
		d.Set("subject", *attr.Sub)
	}

	return nil
}

func resourceIBMAppIDCloudDirectoryUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	password := d.Get("password").(string)
	status := d.Get("status").(string)
	active := d.Get("active").(bool)
	createProfile := d.Get("create_profile").(bool)
	emails := d.Get("email").(*schema.Set)

	input := &appid.StartSignUpOptions{
		TenantID:            &tenantID,
		Active:              &active,
		Emails:              expandAppIDUserEmails(emails.List()),
		Password:            &password,
		Status:              &status,
		ShouldCreateProfile: &createProfile,
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		input.DisplayName = helpers.String(displayName.(string))
	}

	if lockedUntil, ok := d.GetOk("locked_until"); ok {
		input.LockedUntil = core.Int64Ptr(int64(lockedUntil.(int)))
	}

	if userName, ok := d.GetOk("user_name"); ok {
		input.UserName = helpers.String(userName.(string))
	}

	user, resp, err := appIDClient.StartSignUpWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error creating AppID Cloud Directory user: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *user.ID))
	return resourceIBMAppIDCloudDirectoryUserRead(ctx, d, meta)
}

func resourceIBMAppIDCloudDirectoryUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	userID := d.Get("user_id").(string)

	resp, err := appIDClient.DeleteCloudDirectoryUserWithContext(ctx, &appid.DeleteCloudDirectoryUserOptions{
		TenantID: &tenantID,
		UserID:   &userID,
	})

	if err != nil {
		return diag.Errorf("Error deleting AppID Cloud Directory user: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDCloudDirectoryUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	password := d.Get("password").(string)
	status := d.Get("status").(string)
	active := d.Get("active").(bool)
	emails := d.Get("email").(*schema.Set)
	userID := d.Get("user_id").(string)

	input := &appid.UpdateCloudDirectoryUserOptions{
		TenantID: &tenantID,
		Active:   &active,
		Emails:   expandAppIDUserEmails(emails.List()),
		Password: &password,
		Status:   &status,
		UserID:   &userID,
	}

	if displayName, ok := d.GetOk("display_name"); ok {
		input.DisplayName = helpers.String(displayName.(string))
	}

	if lockedUntil, ok := d.GetOk("locked_until"); ok {
		input.LockedUntil = core.Int64Ptr(int64(lockedUntil.(int)))
	}

	if userName, ok := d.GetOk("user_name"); ok {
		input.UserName = helpers.String(userName.(string))
	}
	_, resp, err := appIDClient.UpdateCloudDirectoryUserWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID Cloud Directory user: %s\n%s", err, resp)
	}

	if d.HasChanges("password") {
		password = d.Get("password").(string)

		_, resp, err = appIDClient.ChangePasswordWithContext(ctx, &appid.ChangePasswordOptions{
			TenantID:    &tenantID,
			UUID:        &userID,
			NewPassword: &password,
		})

		if err != nil {
			return diag.Errorf("Error updating AppID Cloud Directory user: %s\n%s", err, resp)
		}
	}

	return resourceIBMAppIDCloudDirectoryUserRead(ctx, d, meta)
}

func expandAppIDUserEmails(e []interface{}) []appid.CreateNewUserEmailsItem {
	if len(e) == 0 {
		return nil
	}

	result := make([]appid.CreateNewUserEmailsItem, len(e))

	for i, item := range e {
		eMap := item.(map[string]interface{})

		email := appid.CreateNewUserEmailsItem{
			Value: helpers.String(eMap["value"].(string)),
		}

		if primary, ok := eMap["primary"]; ok {
			email.Primary = helpers.Bool(primary.(bool))
		}

		result[i] = email
	}

	return result
}
