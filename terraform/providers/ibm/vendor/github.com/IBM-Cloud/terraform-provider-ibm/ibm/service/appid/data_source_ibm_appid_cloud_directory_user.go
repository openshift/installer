package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDCloudDirectoryUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDCloudDirectoryUserRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Determines if the user account is active or not",
				Computed:    true,
			},
			"locked_until": {
				Description: "Integer (epoch time in milliseconds), determines till when the user account will be locked",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud Directory user ID",
			},
			"subject": {
				Description: "The user's identifier ('subject' in identity token)",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud Directory user display name",
			},
			"user_name": {
				Description: "Optional username",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "Current user status: `PENDING` or `CONFIRMED`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "A set of user emails",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value": {
							Description: "User email",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"primary": {
							Description: "`true` if this is primary email",
							Type:        schema.TypeBool,
							Computed:    true,
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

func dataSourceIBMAppIDCloudDirectoryUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	userID := d.Get("user_id").(string)

	user, resp, err := appIDClient.GetCloudDirectoryUserWithContext(ctx, &appid.GetCloudDirectoryUserOptions{
		TenantID: &tenantID,
		UserID:   &userID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID Cloud Directory user: %s\n%s", err, resp)
	}

	d.Set("tenant_id", tenantID)

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
		log.Printf("[DEBUG] Error getting AppID user attributes: %s\n%s", err, resp)
		return diag.Errorf("Error getting AppID user attributes: %s", err)
	}

	if attr.Sub != nil {
		d.Set("subject", *attr.Sub)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *user.ID))
	return nil
}

func flattenAppIDUserMetadata(m *appid.GetUserMeta) []interface{} {
	var result []interface{}

	meta := map[string]interface{}{}

	if m.Created != nil {
		meta["created"] = m.Created.String()
	}

	if m.LastModified != nil {
		meta["last_modified"] = m.LastModified.String()
	}

	result = append(result, meta)

	return result
}

func flattenAppIDUserEmails(e []appid.GetUserEmailsItem) []interface{} {
	var result []interface{}

	for _, v := range e {
		email := map[string]interface{}{
			"value": *v.Value,
		}

		if v.Primary != nil {
			email["primary"] = *v.Primary
		}

		result = append(result, email)
	}

	return result
}
